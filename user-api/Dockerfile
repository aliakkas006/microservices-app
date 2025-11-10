# Stage 1: Build
FROM node:18-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies 
COPY package*.json ./

# Install all dependencies 
RUN npm install

# Copy source code
COPY . .

# Stage 2: Production image
FROM node:18-alpine

WORKDIR /app

# Copy only production dependencies from builder stage
COPY --from=builder /app/node_modules ./node_modules

# Copy source code or built code
COPY --from=builder /app ./

# Use a non-root user for security
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Expose the port your app listens on
EXPOSE 5000

CMD ["node", "src/server.js"]
