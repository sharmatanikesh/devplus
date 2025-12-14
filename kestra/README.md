# Kestra - AI Workflow Engine

This directory contains the setup for [Kestra](https://kestra.io/), the AI workflow orchestration engine that powers DevPlus's intelligent PR analysis, release notes generation, and risk assessment.

## What is Kestra?

Kestra is an open-source orchestration platform that manages AI workflows in DevPlus. The workflows defined in `../backend/workflows/` are executed by Kestra to provide:

- AI-powered pull request analysis
- Automated release notes generation
- Release risk assessment
- Repository analysis

## Prerequisites

- Docker installed locally
- (Optional) Google Cloud SDK (`gcloud`) for cloud deployment

## Local Development

### Quick Start

Run Kestra using the official Docker image:

```bash
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  kestra/kestra:latest
```

Or using PowerShell on Windows:

```powershell
docker run -d `
  --name kestra `
  -p 8080:8080 `
  -v ${PWD}/storage:/app/storage `
  kestra/kestra:latest
```

### Using Custom Dockerfile

If you need custom configurations:

1. **Build the Docker Image**

```bash
docker build -t kestra-local .
```

2. **Run the Container**

```bash
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -v $(pwd)/storage:/app/storage \
  kestra-local
```

### Access Kestra UI

Open your browser and navigate to:

```
http://localhost:8080
```

Default credentials (if authentication is enabled):

- Username: `admin`
- Password: `kestra`

### Uploading Workflows

The DevPlus Kestra workflows are located in `../backend/workflows/`:

- `ai-pull-request-analysis.yaml` - AI PR review workflow
- `ai-release-risk.yaml` - Release risk assessment
- `ai-repo-analysis.yaml` - Repository analysis

Upload these workflows through the Kestra UI or use the Kestra CLI.

### Managing the Container

**Stop the container:**

```bash
docker stop kestra
```

**Start the container:**

```bash
docker start kestra
```

**Remove the container:**

```bash
docker stop kestra
docker rm kestra
```

**View logs:**

```bash
docker logs kestra -f
```

## Environment Variables

Configure Kestra using environment variables:

```bash
docker run -d \
  --name kestra \
  -p 8080:8080 \
  -e KESTRA_CONFIGURATION=/app/application.yaml \
  -e MICRONAUT_ENVIRONMENTS=dev \
  kestra/kestra:latest
```

Common environment variables:

- `KESTRA_CONFIGURATION` - Path to custom configuration file
- `MICRONAUT_ENVIRONMENTS` - Environment (dev, prod)
- `KESTRA_ENCRYPTION_SECRET_KEY` - Secret key for encryption

See [.env.example](./.env.example) for more configuration options.

## Integration with DevPlus

Kestra integrates with the DevPlus backend through webhooks. The backend triggers workflows and receives callbacks:

1. Backend triggers workflow: `POST http://localhost:8080/api/v1/executions/{namespace}/{flowId}`
2. Kestra executes AI workflow
3. Kestra calls back to backend: `POST {BACKEND_URL}/api/v1/webhook/ai`

Configure the backend URL in your workflow YAML files.

## Production Deployment

### Deploy to Google Cloud Run

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

## Troubleshooting

### Common Issues

**Container won't start:**

```bash
# Check logs
docker logs kestra

# Verify port is not in use
netstat -an | grep 8080  # Linux/Mac
netstat -an | findstr 8080  # Windows
```

**Cannot access UI:**

- Ensure Docker container is running: `docker ps | grep kestra`
- Verify port mapping: `docker port kestra`
- Check firewall settings

**Workflow execution fails:**

- Check workflow YAML syntax in Kestra UI
- Verify backend API is accessible from Kestra
- Review execution logs in Kestra UI

### Cloud Run Issues

**Health check failures:**
The Dockerfile includes a health check at `/health`. Ensure Kestra is responding on port 8080.

**Permission issues:**
Cloud Run runs containers as non-root by default, but this Dockerfile uses `USER root` for Kestra compatibility.

## Useful Commands

### Local Docker Commands

```bash
# View container stats
docker stats kestra

# Execute command in container
docker exec -it kestra sh

# Restart container
docker restart kestra

# Pull latest Kestra image
docker pull kestra/kestra:latest
```

### Cloud Run Commands

**Update service:**

```bash
gcloud run services update $SERVICE_NAME \
  --region $REGION \
  --memory 4Gi
```

**Delete service:**

```bash
gcloud run services delete $SERVICE_NAME --region $REGION
```

**List all services:**

```bash
gcloud run services list
```

## Additional Resources

- [Kestra Documentation](https://kestra.io/docs)
- [Kestra Docker Hub](https://hub.docker.com/r/kestra/kestra)
- [DevPlus Backend README](../backend/README.md) - Integration details
- [DevPlus Workflows](../backend/workflows/) - Workflow definitions
- [Google Cloud Run Documentation](https://cloud.google.com/run/docs) - For cloud deployment
