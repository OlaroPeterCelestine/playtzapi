# Deploy CORS Fix to Railway

## Quick Deploy Steps

### 1. Commit and Push Changes

```bash
cd /Users/olaro/Documents/GitHub/playtzapi

# Check what changed
git status

# Add all changes
git add .

# Commit
git commit -m "Fix CORS configuration for cross-origin requests with credentials"

# Push to Railway (if connected to GitHub)
git push origin main

# OR if Railway is connected directly, just push
git push
```

### 2. Railway Will Auto-Deploy

Railway should automatically detect the push and redeploy. Check:
- Railway Dashboard → Your Service → Deployments
- Wait for deployment to complete (usually 2-5 minutes)

### 3. Add Environment Variable (Optional)

If you need to add custom origins:

1. Go to Railway Dashboard
2. Select your service
3. Go to **Variables** tab
4. Click **+ New Variable**
5. Add:
   - **Name:** `CORS_ORIGINS`
   - **Value:** `http://localhost:3001,https://your-frontend-domain.com`
6. Save (this will trigger a redeploy)

### 4. Verify Deployment

After deployment completes, test:

```bash
# Test health endpoint
curl https://playtzapi-production.up.railway.app/health

# Test CORS headers
curl -X OPTIONS https://playtzapi-production.up.railway.app/api/v1/auth/me \
  -H "Origin: http://localhost:3001" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

You should see:
```
< HTTP/1.1 204 No Content
< Access-Control-Allow-Origin: http://localhost:3001
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS, PATCH
```

## Troubleshooting 502 Bad Gateway

If you're still getting 502 errors after deployment:

### Check Railway Logs

1. Go to Railway Dashboard
2. Select your service
3. Click **Logs** tab
4. Look for errors like:
   - Database connection failures
   - Port binding issues
   - Missing environment variables

### Common Issues

1. **Database Connection Failed**
   - Check `DATABASE_URL` is set in Railway
   - Verify database is running

2. **Port Not Set**
   - Railway sets `PORT` automatically
   - Make sure your code uses `os.Getenv("PORT")`

3. **Build Failed**
   - Check build logs in Railway
   - Ensure `go.mod` is up to date

### Restart Service

If needed, restart the service:
1. Railway Dashboard → Your Service
2. Click **Settings**
3. Click **Restart**

## Verify CORS is Working

After deployment, test from your React app:

```javascript
// This should now work
fetch('https://playtzapi-production.up.railway.app/api/v1/auth/me', {
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  }
})
.then(res => {
  console.log('CORS working!', res);
})
.catch(err => {
  console.error('CORS error:', err);
});
```

## Expected Behavior

✅ **After deployment:**
- CORS errors should be gone
- 502 errors should be resolved (if server was down)
- Cookies should be sent/received properly
- Preflight OPTIONS requests should work

❌ **If still failing:**
- Check Railway logs for server errors
- Verify environment variables are set
- Check if database is accessible
- Try restarting the Railway service

