# Testing Login and Dashboard

## Quick Test Guide

### Method 1: Browser (Easiest)

1. **Open Login Page:**
   ```
   http://localhost:8080/admin/login
   ```

2. **Login Credentials:**
   - Username: `admin`
   - Password: `admin123`

3. **After Login:**
   - You'll be redirected to: `http://localhost:8080/admin/dashboard`

### Method 2: API Testing (cURL)

#### Step 1: Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "session_id": "...",
  "user": {
    "id": "...",
    "username": "admin",
    "email": "admin@playtz.com",
    "role_name": "Admin"
  }
}
```

#### Step 2: Verify Session
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -b cookies.txt
```

#### Step 3: Access Dashboard API
```bash
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -b cookies.txt
```

#### Step 4: Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

### Method 3: Automated Test Script

Run the test script:
```bash
./scripts/test_login.sh
```

This will test:
- ✅ Server health
- ✅ Login API
- ✅ Session validation
- ✅ Dashboard access
- ✅ Logout

### Method 4: JavaScript/Fetch

```javascript
// Login
const loginResponse = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include', // Important for cookies
  body: JSON.stringify({
    username: 'admin',
    password: 'admin123'
  })
});

const loginData = await loginResponse.json();
console.log('Login:', loginData);

// Get Current User
const userResponse = await fetch('http://localhost:8080/api/v1/auth/me', {
  credentials: 'include'
});
const userData = await userResponse.json();
console.log('User:', userData);

// Access Dashboard
const dashboardResponse = await fetch('http://localhost:8080/api/v1/admin/dashboard', {
  credentials: 'include'
});
const dashboardData = await dashboardResponse.json();
console.log('Dashboard:', dashboardData);
```

### Method 5: Postman/Insomnia

1. **Login Request:**
   - Method: `POST`
   - URL: `http://localhost:8080/api/v1/auth/login`
   - Headers: `Content-Type: application/json`
   - Body (JSON):
     ```json
     {
       "username": "admin",
       "password": "admin123"
     }
     ```
   - **Important:** Enable "Save cookies automatically

2. **Get Current User:**
   - Method: `GET`
   - URL: `http://localhost:8080/api/v1/auth/me`
   - Cookies will be sent automatically

3. **Dashboard:**
   - Method: `GET`
   - URL: `http://localhost:8080/api/v1/admin/dashboard`
   - Cookies will be sent automatically

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/login` | Login user | No |
| POST | `/api/v1/auth/logout` | Logout user | Yes |
| GET | `/api/v1/auth/me` | Get current user | Yes |
| POST | `/api/v1/auth/change-password` | Change password | Yes |

### Admin Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/v1/admin/dashboard` | Get dashboard data | Yes (Admin) |

## Test Credentials

- **Username:** `admin`
- **Password:** `admin123`
- **Email:** `admin@playtz.com`
- **Role:** Admin

## Troubleshooting

### Issue: "Invalid username or password"
- Check if the user exists: `GET /api/v1/users`
- Verify password is correct
- Check if user is active

### Issue: "Not authenticated"
- Make sure cookies are enabled
- Check if session expired (10 minutes inactivity)
- Try logging in again

### Issue: "Account is inactive"
- User account is disabled
- Contact admin to activate account

### Issue: Dashboard returns 401
- Session expired or invalid
- Login again to get new session

## Session Management

- Sessions expire after **10 minutes of inactivity**
- Session ID is stored in HTTP-only cookie: `session_id`
- Cookie path: `/`
- Cookie is httpOnly (not accessible via JavaScript)

## Example: Complete Login Flow

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt -v

# 2. Check session
curl -X GET http://localhost:8080/api/v1/auth/me \
  -b cookies.txt

# 3. Access dashboard
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -b cookies.txt

# 4. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

