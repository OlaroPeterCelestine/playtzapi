# CORS Configuration Fix

## Problem

Your React app running on `http://localhost:3001` was getting CORS errors when trying to access the API at `https://playtzapi-production.up.railway.app`.

**Error:**
```
Access to fetch at 'https://playtzapi-production.up.railway.app/api/v1/auth/me' 
from origin 'http://localhost:3001' has been blocked by CORS policy: 
No 'Access-Control-Allow-Origin' header is present on the requested resource.
```

## Solution

The CORS configuration has been updated to:
1. **Allow credentials** (cookies) for cross-origin requests
2. **Explicitly allow** common development origins
3. **Support preflight** OPTIONS requests
4. **Configurable** via environment variable

## Configuration

### Default Allowed Origins
- `http://localhost:3000`
- `http://localhost:3001`
- `http://localhost:5173` (Vite default)
- `http://localhost:8080`

### Environment Variable

You can configure allowed origins via the `CORS_ORIGINS` environment variable:

```bash
# In .env file or Railway environment variables
CORS_ORIGINS=http://localhost:3001,https://your-frontend-domain.com
```

**Note:** Separate multiple origins with commas.

## For Railway Production

Add the `CORS_ORIGINS` environment variable in Railway:

1. Go to your Railway project
2. Select your service
3. Go to **Variables** tab
4. Add new variable:
   - **Name:** `CORS_ORIGINS`
   - **Value:** `https://your-frontend-domain.com,http://localhost:3001`

## Testing

After deploying, test CORS with:

```bash
# Test from your React app
curl -X GET https://playtzapi-production.up.railway.app/api/v1/auth/me \
  -H "Origin: http://localhost:3001" \
  -H "Cookie: session_id=your-session-id" \
  -v
```

You should see:
```
< HTTP/1.1 200 OK
< Access-Control-Allow-Origin: http://localhost:3001
< Access-Control-Allow-Credentials: true
```

## Important Notes

1. **Credentials:** The API now allows credentials (cookies) to be sent cross-origin. Make sure your frontend includes `credentials: 'include'` in fetch requests.

2. **502 Bad Gateway:** If you're still getting 502 errors, the Railway server might be down. Check:
   - Railway dashboard for service status
   - Server logs in Railway
   - Database connection

3. **Frontend Configuration:** Your React app should use:
   ```javascript
   fetch('https://playtzapi-production.up.railway.app/api/v1/auth/me', {
     credentials: 'include', // Important!
     headers: {
       'Content-Type': 'application/json'
     }
   })
   ```

## Troubleshooting

### Still getting CORS errors?

1. **Check allowed origins:** Make sure your frontend URL is in the `CORS_ORIGINS` list
2. **Check credentials:** Ensure `credentials: 'include'` is set in fetch requests
3. **Check server status:** Verify the Railway server is running (check `/health` endpoint)
4. **Check headers:** Make sure you're not sending disallowed headers

### 502 Bad Gateway?

1. **Check Railway logs:** The server might have crashed
2. **Check database:** Database connection might be failing
3. **Restart service:** Try redeploying on Railway
4. **Check environment variables:** Make sure all required env vars are set

## Example Frontend Code

```javascript
// Login
const loginResponse = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/login', {
  method: 'POST',
  credentials: 'include', // CRITICAL: Allows cookies
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    username: 'admin',
    password: 'admin123'
  })
});

// Get current user
const userResponse = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/me', {
  credentials: 'include', // CRITICAL: Sends cookies
  headers: {
    'Content-Type': 'application/json',
  }
});
```

