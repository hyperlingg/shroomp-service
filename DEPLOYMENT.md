# GitHub Actions Deployment Setup

This guide explains how to set up automated deployment of the Shroomp backend to Google Cloud Run.

## Prerequisites

- Google Cloud Project with billing enabled
- Artifact Registry repository created
- GitHub repository for this service

## Setup Steps

### 1. Create a Service Account in GCP

```bash
# Set your project ID
export PROJECT_ID="your-project-id"

# Create service account
gcloud iam service-accounts create github-actions \
  --display-name="GitHub Actions Deployment"

# Grant necessary permissions
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/artifactregistry.writer"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

# Create and download key
gcloud iam service-accounts keys create github-actions-key.json \
  --iam-account=github-actions@${PROJECT_ID}.iam.gserviceaccount.com
```

### 2. Add GitHub Secrets

Go to your GitHub repository → Settings → Secrets and variables → Actions → New repository secret

Add the following secrets:

| Secret Name | Description | Example Value |
|-------------|-------------|---------------|
| `GCP_PROJECT_ID` | Your GCP project ID | `my-project-123456` |
| `GCP_REGION` | GCP region for deployment | `us-central1` |
| `GCP_ARTIFACT_REGISTRY` | Artifact Registry repository name | `shroomp` |
| `GCP_SA_KEY` | Content of `github-actions-key.json` | `{ "type": "service_account", ... }` |

**Important:** For `GCP_SA_KEY`, copy the entire contents of the `github-actions-key.json` file.

### 3. Create Artifact Registry Repository (if not exists)

```bash
gcloud artifacts repositories create shroomp \
  --repository-format=docker \
  --location=us-central1 \
  --description="Shroomp container images"
```

### 4. Verify GitHub Action Configuration

The workflow file `.github/workflows/deploy.yml` will:

1. Trigger on:
   - Push to `main` branch
   - Manual workflow dispatch

2. Build the Docker image from `./Dockerfile`

3. Push to Artifact Registry with tags:
   - `latest` (always points to most recent)
   - `<commit-sha>` (specific version)

4. Deploy to Cloud Run with:
   - 512Mi memory
   - 1 CPU
   - Auto-scaling: 0-10 instances
   - Port 8080
   - Public access (unauthenticated)

### 5. Deploy

Push your code to the `main` branch:

```bash
git add .
git commit -m "Add GitHub Actions deployment"
git push origin main
```

Or trigger manually:
- Go to GitHub → Actions → "Deploy to Cloud Run" → Run workflow

### 6. Monitor Deployment

1. Go to GitHub Actions tab to see the workflow progress
2. Once complete, the deployment URL will be shown in the logs
3. Access your API at: `https://shroomp-backend-<hash>-uc.a.run.app`

## Configuration Options

### Adjust Cloud Run Resources

Edit `.github/workflows/deploy.yml` to change:

```yaml
--memory 512Mi      # Change to 1Gi, 2Gi, etc.
--cpu 1             # Change to 2, 4, etc.
--min-instances 0   # Change for faster cold starts
--max-instances 10  # Change based on traffic needs
```

### Change Deployment Region

Update the `GCP_REGION` secret to deploy to a different region:
- `us-central1` (Iowa)
- `us-east1` (South Carolina)
- `europe-west1` (Belgium)
- `asia-northeast1` (Tokyo)

### Add Environment Variables

Add to the `gcloud run deploy` command:

```yaml
--set-env-vars "PORT=8080,ENV=production,LOG_LEVEL=info"
```

Or use secrets:

```yaml
--set-secrets "API_KEY=api-key:latest"
```

## Data Persistence

⚠️ **Important:** Cloud Run containers are stateless. Your current `data.json` file storage won't persist.

### Options:

1. **Cloud Storage** (Recommended):
   - Store data.json in a Cloud Storage bucket
   - Mount using Cloud Run volume mounts

2. **Cloud SQL**:
   - Use PostgreSQL or MySQL instead of JSON file
   - Connect via Cloud SQL Proxy

3. **Firestore**:
   - Use NoSQL database for mushroom sightings
   - Native GCP integration

### Example: Cloud Storage Setup

```bash
# Create bucket
gsutil mb gs://${PROJECT_ID}-shroomp-data

# Grant service account access
gsutil iam ch serviceAccount:github-actions@${PROJECT_ID}.iam.gserviceaccount.com:objectAdmin \
  gs://${PROJECT_ID}-shroomp-data
```

Update deployment in `.github/workflows/deploy.yml`:

```yaml
--update-env-vars "DATA_FILE=gs://${PROJECT_ID}-shroomp-data/data.json"
```

Then update your Go code to use Cloud Storage SDK:

```go
import "cloud.google.com/go/storage"
```

## Troubleshooting

### Deployment fails with "Permission denied"

- Verify service account has all three roles listed above
- Check that Artifact Registry repository exists
- Ensure Cloud Run API is enabled

### Image push fails

```bash
# Enable Artifact Registry API
gcloud services enable artifactregistry.googleapis.com

# Verify authentication
gcloud auth configure-docker us-central1-docker.pkg.dev
```

### Service not accessible

- Check that `--allow-unauthenticated` is set
- Verify firewall rules in GCP
- Check Cloud Run logs: `gcloud run services logs read shroomp-backend --region=us-central1`

## Cost Estimation

Cloud Run pricing (as of 2025):

- **CPU**: $0.00002400/vCPU-second
- **Memory**: $0.00000250/GiB-second
- **Requests**: $0.40/million requests
- **Free tier**: 2 million requests/month

Estimated monthly cost for low traffic (~10,000 requests/month):
- **~$0.50 - $2.00/month** (mostly free tier)

## Security Recommendations

For production:

1. **Remove `--allow-unauthenticated`** and add authentication
2. **Enable Cloud Armor** for DDoS protection
3. **Use Secret Manager** for sensitive data
4. **Add rate limiting**
5. **Enable Cloud Logging and Monitoring**
6. **Use custom domain with SSL**

## Testing Locally

Before pushing, test the Docker build:

```bash
# Build
docker build -t shroomp-backend:test .

# Run
docker run -p 8080:8080 shroomp-backend:test

# Test
curl http://localhost:8080/items
```

## Next Steps

- [ ] Set up database for persistent storage
- [ ] Add authentication (Firebase Auth, OAuth)
- [ ] Configure custom domain
- [ ] Set up monitoring and alerts
- [ ] Add staging environment
- [ ] Implement CI/CD testing before deployment
