# Login Endpoint Guide

## Endpoint

**URL:** `https://playtzapi-production.up.railway.app/api/v1/auth/login`  
**Method:** `POST`  
**Content-Type:** `application/json`

## Request

### Headers
```
Content-Type: application/json
Origin: http://localhost:3001 (for CORS)
```

### Body
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Note:** You can use either `username` or `email` in the `username` field.

## Response

### Success (200 OK)
```json
{
  "success": true,
  "message": "Login successful",
  "session_id": "abc123...",
  "user": {
    "id": "822181e3-f800-4139-ba8a-cfb2f61d9ee6",
    "email": "admin@playtz.com",
    "username": "admin",
    "first_name": "Admin",
    "last_name": "User",
    "role_id": "ab66ab81-4366-40c9-b9ad-510b3f87b62d",
    "role_name": "Admin",
    "active": true,
    "created_at": "2026-01-01T04:15:56Z",
    "updated_at": "2026-01-01T04:15:56Z"
  }
}
```

**Cookie Set:** `session_id` (HTTP-only, 10 minutes)

### Error Responses

#### 400 Bad Request
```json
{
  "error": "Username and password are required"
}
```

#### 401 Unauthorized
```json
{
  "error": "Invalid username or password"
}
```

#### 403 Forbidden
```json
{
  "error": "Account is inactive"
}
```

## JavaScript/Fetch Example

```javascript
const login = async (username, password) => {
  try {
    const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/login', {
      method: 'POST',
      credentials: 'include', // CRITICAL: Allows cookies
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        username: username,
        password: password
      })
    });

    const data = await response.json();

    if (response.ok && data.success) {
      console.log('Login successful!', data.user);
      // Session cookie is automatically set
      return data;
    } else {
      console.error('Login failed:', data.error);
      throw new Error(data.error || 'Login failed');
    }
  } catch (error) {
    console.error('Network error:', error);
    throw error;
  }
};

// Usage
login('admin', 'admin123')
  .then(user => {
    console.log('Logged in as:', user.user.username);
  })
  .catch(error => {
    console.error('Login error:', error);
  });
```

## cURL Example

```bash
curl -X POST https://playtzapi-production.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3001" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt \
  -v
```

## React Example

```jsx
import { useState } from 'react';

function LoginForm() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/login', {
        method: 'POST',
        credentials: 'include', // IMPORTANT!
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        // Login successful - cookie is set automatically
        console.log('Logged in:', data.user);
        // Redirect or update state
        window.location.href = '/dashboard';
      } else {
        setError(data.error || 'Login failed');
      }
    } catch (err) {
      setError('Network error. Please try again.');
      console.error('Login error:', err);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      {error && <div className="error">{error}</div>}
      <button type="submit">Login</button>
    </form>
  );
}
```

## Important Notes

1. **Credentials Required:** Always include `credentials: 'include'` in fetch requests to send/receive cookies.

2. **Session Cookie:** After successful login, a `session_id` cookie is automatically set. This cookie is:
   - HTTP-only (not accessible via JavaScript)
   - Valid for 10 minutes of inactivity
   - Required for authenticated requests

3. **CORS:** The endpoint supports CORS from:
   - `http://localhost:3001`
   - `http://localhost:3000`
   - `http://localhost:5173`
   - `http://localhost:8080`
   - Or any origin set in `CORS_ORIGINS` environment variable

4. **Subsequent Requests:** After login, include the session cookie in all authenticated requests:
   ```javascript
   fetch('https://playtzapi-production.up.railway.app/api/v1/auth/me', {
     credentials: 'include', // Sends the session cookie
     headers: {
       'Content-Type': 'application/json',
     }
   })
   ```

## Testing

### Test Login
```bash
curl -X POST https://playtzapi-production.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt
```

### Test Session (After Login)
```bash
curl -X GET https://playtzapi-production.up.railway.app/api/v1/auth/me \
  -b cookies.txt
```

### Test Logout
```bash
curl -X POST https://playtzapi-production.up.railway.app/api/v1/auth/logout \
  -b cookies.txt
```

## Troubleshooting

### CORS Error
- Make sure `credentials: 'include'` is set
- Check that your origin is allowed
- Verify the server has the latest CORS configuration deployed

### 502 Bad Gateway
- Server might be down - check Railway dashboard
- Check Railway logs for errors
- Try restarting the service

### Invalid Credentials
- Verify username/password are correct
- Check if user account is active
- Ensure user exists in database

