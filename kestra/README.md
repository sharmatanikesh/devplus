# Kestra Deployment Configuration

This folder contains the infrastructure setup for deploying Kestra workflow engine to Google Cloud Run.

## Quick Deploy to Cloud Run

### Option 1: Build and Deploy with Docker

**Local build and test:**
```bash
cd /Users/tanikesh/Documents/devplus/kestra

# Build the image
docker build -t devplus-kestra .

# Test locally
docker run -p 8080:8080 --env-file .env devplus-kestra

# Access Kestra UI at http://localhost:8080
```

**Deploy to Cloud Run:**
```bash
# Build and push to GCR
gcloud builds submit --config cloudbuild.yaml

# Or manually:
docker build -t gcr.io/YOUR_PROJECT_ID/devplus-kestra .
docker push gcr.io/YOUR_PROJECT_ID/devplus-kestra

gcloud run deploy devplus-kestra \
  --image gcr.io/YOUR_PROJECT_ID/devplus-kestra \
  --region asia-south1 \
  --port 8080 \
  --allow-unauthenticated
```

### Option 2: Using Cloud Run Console

1. Go to https://console.cloud.google.com/run
2. Click **"Create Service"**
3. Select **"Continuously deploy from a repository"**
4. Connect GitHub repo
5. Set build configuration:
   - Dockerfile: `kestra/Dockerfile`
   - Build context: `kestra`
6. Service name: `devplus-kestra`
7. Region: `asia-south1`
8. Port: `8080`
9. Click **"Create"**

### Option 3: Using Pre-built Image (Quickest)

1. Go to https://console.cloud.google.com/run
2. Click **"Create Service"**
3. Select **"Deploy one revision from an existing container image"**
4. Container image URL: `kestra/kestra:latest`
5. Service name: `devplus-kestra`
6. Region: `asia-south1`
7. Port: `8080`
8. CPU allocation: **CPU is always allocated**
9. Min instances: `1` (keeps Kestra running)
10. Max instances: `1`
11. Authentication: Allow unauthenticated (or require auth for security)
12. Click **"Create"**

### Option 2: Using gcloud CLI

```bash
gcloud run deploy devplus-kestra \
  --image kestra/kestra:latest \
  --platform managed \
  --region asia-south1 \
  --port 8080 \
  --allow-unauthenticated \
  --min-instances 1 \
  --max-instances 1 \
  --cpu-always-allocated \
  --command "server" \
  --args "standalone"
```

## Environment Variables

Add these in Cloud Run Console → Variables & Secrets:

```bash
SECRET_GEMINI_API_KEY=QUl6YVN5QTZMWmF1T1AyYWdJYXpzZ1RGZ
```

## After Deployment

1. Cloud Run will give you a URL like:
   ```
   https://devplus-kestra-XXXXX.asia-south1.run.app
   ```

2. Update your backend service environment variable:
   ```bash
   KESTRA_URL=https://devplus-kestra-XXXXX.asia-south1.run.app
   ```

3. Access Kestra UI at the Cloud Run URL

## Workflows

Your workflow files are in `../backend/workflows/`:
- `ai-pull-request-analysis.yaml`
- `ai-repo-analysis.yaml`

These will be accessible to Kestra once deployed.

## Cost

- **Always-on instance**: ~$10-15/month
- **Auto-scaling**: Pay per request (cheaper for low usage)

Choose based on your workflow frequency.
