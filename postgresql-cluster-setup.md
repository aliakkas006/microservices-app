# PostgreSQL Cluster Setup (Primary + Replica) in Separate VMs

This project demonstrates setting up a **PostgreSQL cluster** with a **Primary (Master) DB** and a **Standby (Replica) DB**, running inside **two separate VMs**.

---

## High-Level Plan

### VMs

| VM  | Role             | Description       |
|-----|-----------------|-----------------|
| VM1 | `master_pg`     | Primary Database |
| VM2 | `slave_pg`      | Replica Database |

---

## Replication Type

We use **Streaming Replication**, a built-in PostgreSQL feature:

- The primary database continuously ships **WAL (Write-Ahead Log)** entries to the replica.
- The replica applies these WAL logs to stay in sync with the primary.

---

## Networking

Since the VMs are separate:

- We **do not use `network_mode: host`**.
- Options for connectivity:
  - **Custom bridge network** in Docker Compose so containers can communicate across hosts.

---

## Docker Install (Master + Slave)

- `docker_install.sh`

```bash
#!/bin/bash

# Define the target OS

#TARGET_OS="ubuntu"  # Change to "centos" if needed

TARGET_OS=$(cat /etc/os-release | awk -F= '/^ID=/ {gsub(/"/, "", $2); print $2}')
VERSION=$(cat /etc/os-release | awk -F= '/^VERSION=/ {gsub(/"/, "", $2); print $2}')

# Docker installation function for Ubuntu
install_docker_ubuntu() {
    echo "Installing Docker on Ubuntu $VERSION..."
    for pkg in docker.io docker-doc docker-compose docker-compose-v2 podman-docker containerd runc; do sudo apt-get remove $pkg; done
    sudo apt-get update && sudo apt-get upgrade -y
    sudo apt-get install -y ca-certificates curl
    sudo install -m 0755 -d /etc/apt/keyrings
    sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
    sudo chmod a+r /etc/apt/keyrings/docker.asc
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
    sudo apt-get update
    sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    sudo systemctl start docker
    sudo systemctl enable docker
    sudo apt-get install net-tools vim -y
    sudo apt-get install chrony -y
    sudo systemctl restart chrony
    sudo sed -i '/pool/s/^/#/g' /etc/chrony/chrony.conf
    sudo sed -i '/pool 2.ubuntu.pool.ntp.org/a server 0.asia.pool.ntp.org iburst' /etc/chrony/chrony.conf
    sudo systemctl restart chrony
    chronyc sources
    sudo usermod -aG docker $USER
    echo "Docker installed successfully on Ubuntu."
}

# Docker installation function for CentOS
install_docker_centos() {
    echo "Installing Docker on CentOS $VERSION..."
    sudo sed -i 's/mirrorlist/#mirrorlist/g' /etc/yum.repos.d/CentOS-*
    sudo sed -i 's|#baseurl=http://mirror.centos.org|baseurl=http://vault.centos.org|g' /etc/yum.repos.d/CentOS-*
    sudo yum update -y
    sudo yum remove -y docker docker-client docker-client-latest docker-common docker-latest docker-latest-logrotate docker-logrotate docker-engine
    sudo yum install -y yum-utils net-tools
    sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
    sudo yum install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
    sudo systemctl start docker
    sudo systemctl enable docker
    echo "Docker installed successfully on CentOS."
}

# Main script logic
if [[ "$TARGET_OS" == "ubuntu" ]]; then
    install_docker_ubuntu
elif [[ "$TARGET_OS" == "centos" ]]; then
    install_docker_centos
else
    echo "Unsupported OS. Please set TARGET_OS to 'ubuntu' or 'centos'."
    exit 1
fi

```

- `Dockerfile`
```bash
FROM postgres:17.2

RUN apt-get update && \
    apt-get install -y rsync && \
    apt-get install -y postgresql-17-postgis-3 postgis postgresql-17-postgis-3-scripts && \
    rm -rf /var/lib/apt/lists/*

```

## Prepare Master DB

### 1. PostgreSQL Master Config
- `config/postgresql.conf`
```bash
# Connection Settings
listen_addresses = '*'
port = 5432

# Replication Settings
wal_level = replica
max_wal_senders = 10
max_replication_slots = 10
hot_standby = on

# Archive Settings
#archive_mode = on
#archive_command = 'test ! -f /var/lib/postgresql/archive/%f && cp %p /var/lib/postgresql/archive/%f'

# Memory and Performance
shared_buffers = 128MB
work_mem = 4MB

```

- `config/pg_hba.conf`

```bash
# TYPE  DATABASE        USER            ADDRESS                 METHOD
local   all             all                                     trust
host    all             all             127.0.0.1/32            md5
host    all             all             ::1/128                 md5
host    replication     all             192.168.56.0/24         md5
host    all             all             0.0.0.0/0               md5
host    replication     replicator      192.168.56.0/24         md5

```

### 2. `docker-compose.yaml`

```yaml
#version: '3'
services:
  master_pg:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: masterdb
    hostname: masterdb
    restart: unless-stopped
    volumes:
      - ./data:/var/lib/postgresql/data
      - ./config/postgresql.conf:/etc/postgresql/postgresql.conf
      - ./config/pg_hba.conf:/etc/postgresql/pg_hba.conf
    ports:
      - "5432:5432"
    env_file:
      - path: .env
        required: true
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U $POSTGRES_USER -d $POSTGRES_DB -h 0.0.0.0 -p 5432"]
      interval: 10s
      timeout: 5s
      retries: 5
    command: >
      bash -c "
      docker-entrypoint.sh postgres -c config_file=/etc/postgresql/postgresql.conf -c hba_file=/etc/postgresql/pg_hba.conf
      "

```

- `.env`
```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=postgres

```

### 3. Check the Setup

- inside the master_pg container:
`sudo docker exec -it masterdb psql -U postgres -d postgres `

- Creat a replication user on master:
`CREATE USER replicator WITH REPLICATION ENCRYPTED PASSWORD 'replicator_password'`

- Check the Role:
`\du`

- Check master's WAL settings:
`SELECT name, setting FROM pg_settings 
WHERE name IN ('wal_level', 'max_wal_senders', 'max_replication_slots');`

## Prepare Slave DB

### 1. Initialize Slave DB

- `init_slave.sh`

```bash
#!/bin/bash
set -e

echo "Checking if slave initialization is needed..."

# Only initialize if data directory is empty
if [ -z "$(ls -A /var/lib/postgresql/data)" ]; then
    echo "Initializing slave database from master ${POSTGRES_MASTER_HOST}:${POSTGRES_MASTER_PORT}..."
    
    # Wait for master to be available using the replication user
    echo "Waiting for master database to be ready..."
    until PGPASSWORD=${REPLICATION_PASSWORD} pg_isready -h ${POSTGRES_MASTER_HOST} -p ${POSTGRES_MASTER_PORT} -U ${REPLICATION_USER}; do
        echo "Master not ready yet, retrying in 5 seconds..."
        sleep 5
    done
    
    # Perform base backup using replication user
    echo "Starting base backup from master..."
    export PGPASSWORD=${REPLICATION_PASSWORD}
    pg_basebackup -h ${POSTGRES_MASTER_HOST} -p ${POSTGRES_MASTER_PORT} -D /var/lib/postgresql/data \
                 -U ${REPLICATION_USER} -v -P --wal-method=stream --progress -R
    
    # Additional configuration - use replication user credentials
    echo "primary_conninfo = 'host=${POSTGRES_MASTER_HOST} port=${POSTGRES_MASTER_PORT} user=${REPLICATION_USER} password=${REPLICATION_PASSWORD} application_name=slavedb'" >> /var/lib/postgresql/data/postgresql.auto.conf
    echo "hot_standby = on" >> /var/lib/postgresql/data/postgresql.auto.conf
    
    # Create standby signal file
    touch /var/lib/postgresql/data/standby.signal
    
    echo "Slave initialization complete."
else
    echo "Data directory not empty, assuming already initialized."
fi

```

- `chmod +x init_slave.sh`

### 2. `docker-compose.yaml`

```yaml
# version: '3.8'
services:
  slave_pg:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: slavedb
    hostname: slavedb
    restart: unless-stopped
    volumes:
      - ./data:/var/lib/postgresql/data
      - ./init_slave.sh:/docker-entrypoint-initdb.d/init-slave.sh
    ports:
      - "5432:5432"
    env_file:
      - .env
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_MASTER_HOST=${POSTGRES_MASTER_HOST}
      - POSTGRES_MASTER_PORT=${POSTGRES_MASTER_PORT}
      - REPLICATION_USER=${REPLICATION_USER}
      - REPLICATION_PASSWORD=${REPLICATION_PASSWORD}

```

- `.env`
```bash
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=postgres
POSTGRES_MASTER_HOST=192.168.56.100
POSTGRES_MASTER_PORT=5432

REPLICATION_USER=replicator
REPLICATION_PASSWORD=replicator_password

```

### 2. Verify the slave_pg
 - inside the slave_pg container:
`sudo docker exec -it slavedb psql -U postgres -d postgres `

- Check if slave is in recovery mode:
`SELECT pg_is_in_recovery();`

- Check replication status:
`SELECT * FROM pg_stat_wal_receiver;`



