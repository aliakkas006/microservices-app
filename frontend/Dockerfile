# Stage 1: Install dependencies
FROM node:20-alpine AS deps

WORKDIR /app

# Install dependencies based on lock file for reproducibility
COPY package.json package-lock.json ./
RUN npm ci --legacy-peer-deps


# Stage 2: Dev
FROM deps AS development

WORKDIR /app

# Copy all source files
COPY . .

# Expose React's default port
EXPOSE 3000

# Run React in development mode
CMD ["npm", "start"]

