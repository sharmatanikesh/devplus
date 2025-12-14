# Use official Kestra image as base
FROM kestra/kestra:latest

# Switch to root for any customizations
USER root

# Set working directory
WORKDIR /app

# Expose Kestra port
EXPOSE 8080

# Set environment variables (can be overridden at runtime)
ENV KESTRA_CONFIGURATION=""

# Run Kestra server in standalone mode
CMD ["server", "standalone"]
