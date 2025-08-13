#!/usr/bin/env bash
set -euo pipefail

# Configuration with defaults
BOOTSTRAP="${BOOTSTRAP:-kafka-1:19092}"
TOPIC="${TOPIC_NAME:-todo_created}"
PARTS="${TOPIC_PARTITIONS:-3}"
RF="${TOPIC_RF:-2}"
MIN_ISR=2
MAX_RETRIES=30
RETRY_DELAY=2

echo "Starting Kafka topic initialization for topic '$TOPIC'"

# Verify required commands exist
command -v /opt/bitnami/kafka/bin/kafka-metadata-quorum.sh >/dev/null 2>&1 || {
    echo "Error: kafka-metadata-quorum.sh not found"
    exit 1
}

command -v /opt/bitnami/kafka/bin/kafka-topics.sh >/dev/null 2>&1 || {
    echo "Error: kafka-topics.sh not found"
    exit 1
}

# Wait for Kafka cluster to be ready
echo "Waiting for Kafka to be reachable at ${BOOTSTRAP%%,*} (max $MAX_RETRIES attempts)..."
for ((i=1; i<=MAX_RETRIES; i++)); do
    if /opt/bitnami/kafka/bin/kafka-metadata-quorum.sh \
        --bootstrap-server "${BOOTSTRAP%%,*}" \
        describe > /dev/null 2>&1; then
        echo "Kafka cluster is ready"
        break
    fi
    
    if [[ $i -eq $MAX_RETRIES ]]; then
        echo "Error: Kafka cluster not ready after $MAX_RETRIES attempts"
        exit 1
    fi
    
    echo "Retrying ($i/$MAX_RETRIES)..."
    sleep $RETRY_DELAY
done

# Create topic with configuration
echo "Creating topic '$TOPIC' with:"
echo "  Partitions: $PARTS"
echo "  Replication Factor: $RF"
echo "  min.insync.replicas: $MIN_ISR"

/opt/bitnami/kafka/bin/kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP" \
    --create \
    --if-not-exists \
    --topic "$TOPIC" \
    --partitions "$PARTS" \
    --replication-factor "$RF" \
    --config "min.insync.replicas=$MIN_ISR" || {
    echo "Error: Failed to create topic '$TOPIC'"
    exit 1
}

# Verify topic creation and configuration
echo "Verifying topic configuration..."
/opt/bitnami/kafka/bin/kafka-topics.sh \
    --bootstrap-server "$BOOTSTRAP" \
    --describe \
    --topic "$TOPIC" || {
    echo "Error: Topic verification failed"
    exit 1
}

echo "Topic '$TOPIC' initialized successfully"
exit 0