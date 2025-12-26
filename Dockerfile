# Stage 1: Build Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

# Install dependencies
COPY frontend/package*.json ./
RUN npm install

# Build frontend
COPY frontend .
RUN npm run build

# Stage 2: Build Backend
FROM golang:1.24-alpine AS backend-builder
WORKDIR /app

# Install gcc/musl-dev for SQLite CGO support
RUN apk add --no-cache gcc musl-dev

# Download Go dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go ./

# Copy built frontend assets from Stage 1
# Note: npm run build outputs to ../dist relative to frontend, so it's at /app/dist
COPY --from=frontend-builder /app/dist ./dist

# Build the application
# CGO_ENABLED=1 is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o toolbox .

# Stage 3: Runtime
FROM alpine:latest
WORKDIR /app

# Install timezone data and certificates
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=backend-builder /app/toolbox .

# Create data directory for SQLite persistence
RUN mkdir /data
VOLUME /data

# Expose port
EXPOSE 8080

CMD ["./toolbox"]
