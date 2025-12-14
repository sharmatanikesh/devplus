# Use the official Kestra image
FROM kestra/kestra:latest

# Cloud Run expects the service to listen on port 8080
ENV MICRONAUT_SERVER_PORT=8080

# Expose the port (for local testing only)
EXPOSE 8080

# Start Kestra in local mode (simple mode without Kubernetes)
CMD ["server", "local"]
