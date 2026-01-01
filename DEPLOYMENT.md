# Railway Deployment Guide

## Prerequisites
1. Railway account (sign up at https://railway.app)
2. GitHub repository with your code
3. Railway CLI (optional, for local testing)

## Deployment Steps

### Option 1: Deploy via Railway Dashboard (Recommended)

1. **Connect Repository**
   - Go to https://railway.app
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose your repository

2. **Configure Environment Variables**
   - In Railway dashboard, go to your project → Variables
   - Add the following environment variables:
   
   ```
   DATABASE_URL=postgresql://postgres:yeTIxpqsTxzwURaMfBbnAWFjdGfUVAvT@shinkansen.proxy.rlwy.net:20263/railway
   CLOUDINARY_CLOUD_NAME=dgj7pnuwh
   CLOUDINARY_API_KEY=736681341159732
   CLOUDINARY_API_SECRET=DksTG6mc28e8ZJUytm-ZWjOnGdY
   ```
   
   Note: `PORT` is automatically set by Railway, you don't need to set it.

3. **Deploy**
   - Railway will automatically detect the `railway.toml` file
   - It will build and deploy your Go application
   - The build command: `go mod download && go mod tidy && go build -o main .`
   - The start command: `./main`

4. **Monitor Deployment**
   - Check the "Deployments" tab for build logs
   - Check the "Metrics" tab for application health
   - Health check endpoint: `/health`

### Option 2: Deploy via Railway CLI

1. **Install Railway CLI**
   ```bash
   npm i -g @railway/cli
   ```

2. **Login**
   ```bash
   railway login
   ```

3. **Initialize Project**
   ```bash
   railway init
   ```

4. **Set Environment Variables**
   ```bash
   railway variables set DATABASE_URL="postgresql://postgres:yeTIxpqsTxzwURaMfBbnAWFjdGfUVAvT@shinkansen.proxy.rlwy.net:20263/railway"
   railway variables set CLOUDINARY_CLOUD_NAME="dgj7pnuwh"
   railway variables set CLOUDINARY_API_KEY="736681341159732"
   railway variables set CLOUDINARY_API_SECRET="DksTG6mc28e8ZJUytm-ZWjOnGdY"
   ```

5. **Deploy**
   ```bash
   railway up
   ```

## Configuration Files

- `railway.toml` - Railway configuration
- `.railwayignore` - Files to exclude from deployment

## Environment Variables Required

| Variable | Description | Required |
|----------|-------------|----------|
| `DATABASE_URL` | PostgreSQL connection string | Yes |
| `CLOUDINARY_CLOUD_NAME` | Cloudinary cloud name | Yes |
| `CLOUDINARY_API_KEY` | Cloudinary API key | Yes |
| `CLOUDINARY_API_SECRET` | Cloudinary API secret | Yes |
| `PORT` | Server port (auto-set by Railway) | No |

## Health Check

Once deployed, your API will be available at:
- Health check: `https://your-app.railway.app/health`
- API base: `https://your-app.railway.app/api/v1`

## Troubleshooting

1. **Build fails**: Check Go version compatibility (requires Go 1.24.0+)
2. **Database connection fails**: Verify `DATABASE_URL` is correct
3. **Image upload fails**: Verify Cloudinary credentials
4. **Port binding errors**: Railway sets PORT automatically, don't override it

## Post-Deployment

1. Test the health endpoint: `curl https://your-app.railway.app/health`
2. Test an API endpoint: `curl https://your-app.railway.app/api/v1/news`
3. Monitor logs in Railway dashboard
4. Set up custom domain (optional) in Railway settings

