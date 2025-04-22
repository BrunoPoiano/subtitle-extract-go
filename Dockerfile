# Use Go 1.24.2 with Alpine 3.21 as the base image for a lightweight container
FROM golang:1.24.2-alpine3.21

# Set the working directory for the application
WORKDIR /app/subextract

# Copy necessary files for the Go application
COPY start.sh main.go go.mod ./

# Build the Go application
RUN go build .

# Install cron daemon for scheduled tasks
RUN apk add --no-cache dcron

# Install timezone data for proper time handling
RUN apk add --no-cache tzdata

# Install ffmpeg for video processing capabilities
RUN apk add --no-cache ffmpeg

# Ensure we're in the correct working directory
WORKDIR /app/subextract

# Create directory for storing videos to be processed
RUN mkdir -p /app/subextract/videos

# Set up cron job to run the application every 6 hours
# Output is redirected to log file for debugging
RUN echo " * */6 * * * /app/subextract/main >> /var/log/cron.log 2>&1" > /etc/crontabs/root

# Create empty log file for cron output
RUN touch /var/log/cron.log

# Copy startup script and make it executable
COPY start.sh /app/subextract/start.sh
RUN chmod +x /app/subextract/start.sh

# Execute the startup script when container launches
CMD ["/bin/sh", "/app/subextract/start.sh"]
