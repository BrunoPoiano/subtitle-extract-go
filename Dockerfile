# Use a multi-stage build approach to reduce image size
FROM golang:1.24.2-alpine3.21 AS builder

# Set the working directory for the build stage
WORKDIR /build

# Copy only files needed for building
COPY main.go go.mod ./

# Build the Go application with static linking to reduce dependencies
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o subextract .

# Use a smaller Alpine base for the final image
FROM alpine:3.21

# Set the working directory for the application
WORKDIR /app/subextract

# Install only the necessary packages in a single RUN to reduce layers
RUN apk add --no-cache dcron tzdata ffmpeg && \
  mkdir -p /app/subextract/videos && \
  echo "0 */6 * * * /app/subextract/subextract >> /var/log/cron.log 2>&1" > /etc/crontabs/root && \
  touch /var/log/cron.log

# Copy only the compiled binary from the builder stage
COPY --from=builder /build/subextract ./

# Copy startup script and make it executable
COPY start.sh ./
RUN chmod +x ./start.sh

# Execute the startup script when container launches
CMD ["/bin/sh", "./start.sh"]
