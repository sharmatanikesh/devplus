# Kestra Setup for Google Cloud Run

This repository contains a Dockerfile to run [Kestra](https://kestra.io/) on Google Cloud Run.

## Prerequisites

- Docker installed locally
- Google Cloud SDK (`gcloud`) installed (for Cloud Run deployment)
- A Google Cloud project with billing enabled

## Local Development

### 1. Build the Docker Image

```bash
docker build -t kestra-local .
```

### 2. Run Locally

```bash
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  kestra-local
```

### 3. Access Kestra

Open your browser and navigate to:
```
http://localhost:8080
```

### 4. Stop the Container

```bash
docker stop kestra
docker rm kestra
```

## Deploy to Google Cloud Run

### Option 1: Using Cloud Build (Recommended)

```bash
# Set your project ID
export PROJECT_ID=your-project-id
export REGION=us-central1
export SERVICE_NAME=kestra

# Build and push to Google Container Registry
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME

# Deploy to Cloud Run
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 2Gi \
  --cpu 2 \
  --timeout 3600 \
  --max-instances 10
```

### Option 2: Using Artifact Registry

```bash
# Set your project ID
export PROJECT_ID=your-project-id
export REGION=us-central1
export SERVICE_NAME=kestra
export REPOSITORY=kestra-repo

# Create Artifact Registry repository (one-time setup)
gcloud artifacts repositories create $REPOSITORY \
  --repository-format=docker \
  --location=$REGION \
  --description="Kestra Docker repository"

# Configure Docker to use gcloud as credential helper
gcloud auth configure-docker $REGION-docker.pkg.dev

# Build and tag the image
docker build -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$SERVICE_NAME:latest .

# Push to Artifact Registry
docker push $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$SERVICE_NAME:latest

# Deploy to Cloud Run
gcloud run deploy $SERVICE_NAME \
  --image $REGION-docker.pkg.dev/$PROJECT_ID/$REPOSITORY/$SERVICE_NAME:latest \
  --platform managed \
  --region $REGION \
  --allow-unauthenticated \
  --port 8080 \
  --memory 2Gi \
  --cpu 2 \
  --timeout 3600 \
  --max-instances 10
```

## Configuration

### Environment Variables

You can add custom environment variables during deployment:

```bash
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --set-env-vars "KESTRA_CONFIGURATION=/app/application.yaml" \
  --set-env-vars "CUSTOM_VAR=value"
```

### Custom Configuration File

If you need a custom `application.yaml`:

1. Create your `application.yaml` file in this directory
2. Update the Dockerfile to copy it:
   ```dockerfile
   COPY application.yaml /app/application.yaml
   ENV KESTRA_CONFIGURATION=/app/application.yaml
   ```
3. Rebuild and redeploy

### Resource Limits

Adjust memory and CPU based on your workload:

```bash
gcloud run deploy $SERVICE_NAME \
  --memory 4Gi \
  --cpu 4
```

## Monitoring

### View Logs

```bash
gcloud run services logs read $SERVICE_NAME --region $REGION
```

### Check Service Status

```bash
gcloud run services describe $SERVICE_NAME --region $REGION
```

## Troubleshooting

### Container Won't Start

Check the logs:
```bash
gcloud run services logs tail $SERVICE_NAME --region $REGION
```

### Health Check Failures

The Dockerfile includes a health check at `/health`. Ensure Kestra is responding on port 8080.

### Permission Issues

Cloud Run runs containers as non-root by default, but this Dockerfile uses `USER root` for Kestra compatibility. If you encounter permission issues, verify the USER directive in the Dockerfile.

## Useful Commands

### Update Service

```bash
gcloud run services update $SERVICE_NAME \
  --region $REGION \
  --memory 4Gi
```

### Delete Service

```bash
gcloud run services delete $SERVICE_NAME --region $REGION
```

### List All Services

```bash
gcloud run services list
```

## Additional Resources

- [Kestra Documentation](https://kestra.io/docs)
- [Google Cloud Run Documentation](https://cloud.google.com/run/docs)
- [Kestra Docker Hub](https://hub.docker.com/r/kestra/kestra)
