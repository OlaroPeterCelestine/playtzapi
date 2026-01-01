# Quick Fix: CORS and 502 Errors

## The Problem

You're getting:
1. **CORS errors** - "No 'Access-Control-Allow-Origin' header"
2. **502 Bad Gateway** - Server not responding

## The Solution

The code is **already fixed** but needs to be **deployed to Railway**.

## Deploy Now (3 Steps)

### Step 1: Commit Changes
```bash
cd /Users/olaro/Documents/GitHub/playtzapi
git add .
git commit -m "Fix CORS configuration for cross-origin requests"
```

### Step 2: Push to Railway
```bash
git push
```

### Step 3: Wait for Deployment
- Go to Railway Dashboard
- Check your service → Deployments
- Wait 2-5 minutes for deployment to complete

## After Deployment

### Test the Fix

```bash
# Test health (should work)
curl https://playtzapi-production.up.railway.app/health

# Test CORS (should return proper headers)
curl -X OPTIONS https://playtzapi-production.up.railway.app/api/v1/auth/me \
  -H "Origin: http://localhost:3001" \
  -H "Access-Control-Request-Method: GET" \
  -v
```

### If Still Getting 502

1. **Check Railway Logs:**
   - Railway Dashboard → Your Service → Logs
   - Look for errors

2. **Common Issues:**
   - Database connection failed → Check `DATABASE_URL`
   - Build failed → Check build logs
   - Port issue → Railway sets `PORT` automatically

3. **Restart Service:**
   - Railway Dashboard → Your Service → Settings → Restart

## What Was Fixed

✅ CORS now allows `http://localhost:3001`  
✅ Credentials (cookies) enabled for cross-origin  
✅ Preflight OPTIONS requests handled  
✅ Configurable via `CORS_ORIGINS` env var  

## Frontend Code

Make sure your React app uses:

```javascript
fetch('https://playtzapi-production.up.railway.app/api/v1/auth/me', {
  credentials: 'include', // CRITICAL!
  headers: {
    'Content-Type': 'application/json'
  }
})
```

The `credentials: 'include'` is **required** for cookies to work cross-origin.

