# Playtz API Documentation

## Base URL

**Production:** `https://playtzapi-production.up.railway.app/api/v1`  
**Local Development:** `http://localhost:8080/api/v1`

---

## Authentication

All API endpoints (except auth endpoints) require authentication via session cookies.

### Login

**Endpoint:** `POST /auth/login`

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (Success - 200):**
```json
{
  "success": true,
  "message": "Login successful",
  "session_id": "abc123...",
  "user": {
    "id": "uuid",
    "email": "admin@playtz.com",
    "username": "admin",
    "first_name": "Admin",
    "last_name": "User",
    "role_id": "uuid",
    "role_name": "admin",
    "active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response (Error - 401):**
```json
{
  "error": "Invalid username or password"
}
```

**Response (Error - 403):**
```json
{
  "error": "Account is inactive"
}
```

**Notes:**
- Session cookie is automatically set with `HttpOnly`, `Secure`, and `SameSite=None`
- Session expires after 10 minutes of inactivity
- Use `credentials: 'include'` in fetch requests to send cookies
- **401 errors on login are expected for failed attempts** - handle them gracefully in the UI without logging to console

---

### Logout

**Endpoint:** `POST /auth/logout`

**Request:** No body required (uses session cookie)

**Response (Success - 200):**
```json
{
  "message": "Logged out successfully",
  "success": true
}
```

**Notes:**
- Clears session from server
- Clears session cookie
- Works even without valid session

---

### Check Current User

**Endpoint:** `GET /auth/me`

**Request:** No body required (uses session cookie)

**Response (Authenticated - 200):**
```json
{
  "authenticated": true,
  "user": {
    "id": "uuid",
    "email": "admin@playtz.com",
    "username": "admin",
    "first_name": "Admin",
    "last_name": "User",
    "role_id": "uuid",
    "role_name": "admin",
    "active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Response (Not Authenticated - 200):**
```json
{
  "authenticated": false,
  "user": null
}
```

**Notes:**
- Always returns 200 (never 401)
- Use `authenticated` field to check status
- Safe to call on login page without errors

---

### Change Password

**Endpoint:** `POST /auth/change-password`  
**Authentication:** Required

**Request:**
```json
{
  "current_password": "oldpassword",
  "new_password": "newpassword123"
}
```

**Response (Success - 200):**
```json
{
  "message": "Password changed successfully"
}
```

**Response (Error - 401):**
```json
{
  "error": "Current password is incorrect"
}
```

---

## Protected Endpoints

All endpoints below require authentication. Include `credentials: 'include'` in fetch requests.

### Admin Dashboard

**Endpoint:** `GET /admin/dashboard`  
**Authentication:** Required

**Response:**
```json
{
  "user": { ... },
  "stats": {
    "total_users": 10,
    "total_news": 5,
    "total_events": 3,
    "total_merchandise": 8,
    "total_orders": 15,
    "pending_orders": 2
  },
  "recent_news": [ ... ],
  "recent_events": [ ... ],
  "recent_orders": [ ... ]
}
```

---

## News Endpoints

### List News

**Endpoint:** `GET /news`  
**Authentication:** Required

**Response:**
```json
[
  {
    "id": "uuid",
    "title": "News Title",
    "content": "News content...",
    "author": "Author Name",
    "published": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### Get News by ID

**Endpoint:** `GET /news/:id`  
**Authentication:** Required

**Response:**
```json
{
  "id": "uuid",
  "title": "News Title",
  "content": "News content...",
  "author": "Author Name",
  "published": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Create News

**Endpoint:** `POST /news`  
**Authentication:** Required

**Request:**
```json
{
  "title": "News Title",
  "content": "News content...",
  "author": "Author Name",
  "published": false
}
```

**Response (201):**
```json
{
  "id": "uuid",
  "title": "News Title",
  "content": "News content...",
  "author": "Author Name",
  "published": false,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Update News

**Endpoint:** `PUT /news/:id`  
**Authentication:** Required

**Request:**
```json
{
  "title": "Updated Title",
  "content": "Updated content...",
  "author": "Updated Author",
  "published": true
}
```

**Response (200):**
```json
{
  "id": "uuid",
  "title": "Updated Title",
  "content": "Updated content...",
  "author": "Updated Author",
  "published": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Delete News

**Endpoint:** `DELETE /news/:id`  
**Authentication:** Required

**Response (200):**
```json
{
  "message": "News deleted successfully"
}
```

---

## Events Endpoints

### List Events

**Endpoint:** `GET /events`  
**Authentication:** Required

### Get Event by ID

**Endpoint:** `GET /events/:id`  
**Authentication:** Required

### Create Event

**Endpoint:** `POST /events`  
**Authentication:** Required

**Request:**
```json
{
  "title": "Event Title",
  "description": "Event description...",
  "date": "2024-12-31T00:00:00Z",
  "location": "Venue Name",
  "active": true
}
```

### Update Event

**Endpoint:** `PUT /events/:id`  
**Authentication:** Required

### Delete Event

**Endpoint:** `DELETE /events/:id`  
**Authentication:** Required

---

## Merchandise Endpoints

### List Merchandise

**Endpoint:** `GET /merch`  
**Authentication:** Required

### Get Merchandise by ID

**Endpoint:** `GET /merch/:id`  
**Authentication:** Required

### Create Merchandise

**Endpoint:** `POST /merch`  
**Authentication:** Required

**Request:**
```json
{
  "name": "Product Name",
  "description": "Product description...",
  "price": 29.99,
  "stock": 100,
  "active": true
}
```

### Update Merchandise

**Endpoint:** `PUT /merch/:id`  
**Authentication:** Required

### Delete Merchandise

**Endpoint:** `DELETE /merch/:id`  
**Authentication:** Required

---

## Careers Endpoints

### List Careers

**Endpoint:** `GET /careers`  
**Authentication:** Required

### Get Career by ID

**Endpoint:** `GET /careers/:id`  
**Authentication:** Required

### Create Career

**Endpoint:** `POST /careers`  
**Authentication:** Required

**Request:**
```json
{
  "title": "Job Title",
  "description": "Job description...",
  "department": "Department Name",
  "location": "Location",
  "active": true
}
```

### Update Career

**Endpoint:** `PUT /careers/:id`  
**Authentication:** Required

### Delete Career

**Endpoint:** `DELETE /careers/:id`  
**Authentication:** Required

---

## Rooms Endpoints

### List Rooms

**Endpoint:** `GET /rooms`  
**Authentication:** Required

### Get Room by ID

**Endpoint:** `GET /rooms/:id`  
**Authentication:** Required

### Create Room

**Endpoint:** `POST /rooms`  
**Authentication:** Required

### Update Room

**Endpoint:** `PUT /rooms/:id`  
**Authentication:** Required

### Delete Room

**Endpoint:** `DELETE /rooms/:id`  
**Authentication:** Required

---

## Mixes Endpoints

### List Mixes

**Endpoint:** `GET /mixes`  
**Authentication:** Required

### Get Mix by ID

**Endpoint:** `GET /mixes/:id`  
**Authentication:** Required

### Create Mix

**Endpoint:** `POST /mixes`  
**Authentication:** Required

### Update Mix

**Endpoint:** `PUT /mixes/:id`  
**Authentication:** Required

### Delete Mix

**Endpoint:** `DELETE /mixes/:id`  
**Authentication:** Required

### Add Track to Mix

**Endpoint:** `POST /mixes/:id/tracks`  
**Authentication:** Required

### Add Multiple Tracks to Mix

**Endpoint:** `POST /mixes/:id/tracks/bulk`  
**Authentication:** Required

### Remove Track from Mix

**Endpoint:** `DELETE /mixes/:id/tracks`  
**Authentication:** Required

---

## Users Endpoints

### List Users

**Endpoint:** `GET /users`  
**Authentication:** Required

### Get User by ID

**Endpoint:** `GET /users/:id`  
**Authentication:** Required

### Create User

**Endpoint:** `POST /users`  
**Authentication:** Required

**Request:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "optional-password",
  "first_name": "First",
  "last_name": "Last",
  "role_id": "role-uuid"
}
```

**Response (with default password):**
```json
{
  "user": { ... },
  "default_password": "generated-password",
  "message": "User created with default password. Please change it on first login."
}
```

**Response (with provided password):**
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "username": "username",
  ...
}
```

### Update User

**Endpoint:** `PUT /users/:id`  
**Authentication:** Required

### Delete User

**Endpoint:** `DELETE /users/:id`  
**Authentication:** Required

### Update User Role

**Endpoint:** `PUT /users/:id/role`  
**Authentication:** Required

**Request:**
```json
{
  "role_id": "new-role-uuid"
}
```

---

## Roles Endpoints

### List Roles

**Endpoint:** `GET /roles`  
**Authentication:** Required

### Get Role by ID

**Endpoint:** `GET /roles/:id`  
**Authentication:** Required

### Create Role

**Endpoint:** `POST /roles`  
**Authentication:** Required

**Request:**
```json
{
  "name": "role_name",
  "description": "Role description",
  "permissions": ["permission1", "permission2"],
  "active": true
}
```

### Update Role

**Endpoint:** `PUT /roles/:id`  
**Authentication:** Required

### Delete Role

**Endpoint:** `DELETE /roles/:id`  
**Authentication:** Required

---

## Upload Endpoints

### Upload Single Image

**Endpoint:** `POST /upload`  
**Authentication:** Required  
**Content-Type:** `multipart/form-data`

**Request:** Form data with `file` field

**Response:**
```json
{
  "url": "https://res.cloudinary.com/...",
  "secure_url": "https://res.cloudinary.com/...",
  "public_id": "image-id",
  "width": 1920,
  "height": 1080,
  "format": "webp"
}
```

### Upload Multiple Images

**Endpoint:** `POST /upload/multiple`  
**Authentication:** Required  
**Content-Type:** `multipart/form-data`

**Request:** Form data with `files[]` field (array)

**Response:**
```json
{
  "urls": [
    {
      "url": "https://res.cloudinary.com/...",
      "secure_url": "https://res.cloudinary.com/...",
      "public_id": "image-id-1"
    },
    {
      "url": "https://res.cloudinary.com/...",
      "secure_url": "https://res.cloudinary.com/...",
      "public_id": "image-id-2"
    }
  ]
}
```

---

## Cart Endpoints

### Get Cart

**Endpoint:** `GET /cart`  
**Authentication:** Required

### Add to Cart

**Endpoint:** `POST /cart/add`  
**Authentication:** Required

**Request:**
```json
{
  "merch_id": "uuid",
  "quantity": 2
}
```

### Update Cart Item

**Endpoint:** `PUT /cart/update`  
**Authentication:** Required

**Request:**
```json
{
  "merch_id": "uuid",
  "quantity": 3
}
```

### Remove from Cart

**Endpoint:** `DELETE /cart/remove`  
**Authentication:** Required

**Request:**
```json
{
  "merch_id": "uuid"
}
```

### Clear Cart

**Endpoint:** `DELETE /cart/clear`  
**Authentication:** Required

---

## Orders Endpoints

### List Orders

**Endpoint:** `GET /orders`  
**Authentication:** Required

### Get Order by ID

**Endpoint:** `GET /orders/:id`  
**Authentication:** Required

### Create Order (Checkout)

**Endpoint:** `POST /checkout`  
**Authentication:** Required

**Request:**
```json
{
  "items": [
    {
      "merch_id": "uuid",
      "quantity": 2,
      "price": 29.99
    }
  ],
  "total": 59.98,
  "shipping_address": "123 Main St",
  "billing_address": "123 Main St"
}
```

### Update Order Status

**Endpoint:** `PUT /orders/:id/status`  
**Authentication:** Required

**Request:**
```json
{
  "status": "completed"
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "Invalid request body"
}
```

### 401 Unauthorized
```json
{
  "error": "Authentication required"
}
```
or
```json
{
  "error": "Invalid or expired session"
}
```

### 403 Forbidden
```json
{
  "error": "Access denied"
}
```
or
```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found
```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error"
}
```

---

## JavaScript Examples

### Login
```javascript
const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/login', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include',
  body: JSON.stringify({
    username: 'admin',
    password: 'admin123'
  })
});

const data = await response.json();
if (response.ok && data.success) {
  // Redirect to dashboard
  window.location.href = '/dashboard';
} else {
  // Show error to user (don't log to console - 401 is expected for failed attempts)
  // Display error message in UI instead of console.error()
  showErrorMessage(data.error || 'Login failed');
}
```

### Authenticated Request
```javascript
const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/news', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include' // Required for cookies
});

const data = await response.json();
```

### Create News
```javascript
const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/news', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include',
  body: JSON.stringify({
    title: 'News Title',
    content: 'News content...',
    author: 'Author Name',
    published: true
  })
});

const data = await response.json();
```

### Logout
```javascript
const response = await fetch('https://playtzapi-production.up.railway.app/api/v1/auth/logout', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json',
  }
});

// Clear local storage
localStorage.clear();
sessionStorage.clear();

// Redirect to login
window.location.href = '/login';
```

---

## Important Notes

1. **All endpoints require authentication** except:
   - `POST /auth/login`
   - `POST /auth/logout`
   - `GET /auth/me`
   - `GET /health`

2. **Always include `credentials: 'include'`** in fetch requests to send cookies

3. **Session expires after 10 minutes** of inactivity

4. **CORS is configured** for:
   - `http://localhost:3000`
   - `http://localhost:3001`
   - `http://localhost:5173`
   - `http://localhost:8080`
   - `https://playtzadmin.vercel.app`

5. **Default Admin Credentials:**
   - Username: `admin`
   - Password: `admin123`

6. **Session Management:**
   - Sessions are stored server-side
   - Cookies are `HttpOnly`, `Secure`, and `SameSite=None`
   - Session is automatically cleared on logout

7. **Error Handling Best Practices:**
   - **401 errors on `/auth/login` are expected** for failed login attempts
   - Handle these gracefully in the UI without logging to console
   - Only log unexpected errors (500, network failures, etc.)
   - Example: Show error message to user instead of `console.error()`

---

## Health Check

**Endpoint:** `GET /health`  
**Authentication:** Not required

**Response:**
```json
{
  "status": "ok"
}
```

---

## Frontend Error Handling Guide

### Suppressing Console Errors for Login Failures

**401 errors on login are expected** when credentials are incorrect. Handle them gracefully without logging to console:

#### React/Next.js Example (AuthContext.tsx)

```typescript
// ❌ DON'T: Log 401 errors to console
const response = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ username, password })
});

if (!response.ok) {
  console.error('Login failed:', response.status); // ❌ Don't do this
}

// ✅ DO: Handle errors silently and show to user
const response = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ username, password })
});

const data = await response.json();

if (response.ok && data.success) {
  // Success - redirect or update state
  setUser(data.user);
} else {
  // Handle error gracefully - show to user, don't log to console
  if (response.status === 401) {
    // 401 is expected for wrong credentials - handle silently
    setError(data.error || 'Invalid credentials');
    // Don't use console.error() here
  } else if (response.status >= 500) {
    // Only log unexpected server errors
    console.error('Server error:', data.error);
    setError('Server error. Please try again later.');
  } else {
    // Other errors (400, 403, etc.)
    setError(data.error || 'Login failed');
  }
}
```

#### Suppress Specific Console Errors

If you need to suppress console errors globally for 401 on login:

```typescript
// Override console.error temporarily for login endpoint
const originalError = console.error;
console.error = (...args) => {
  // Suppress 401 errors from login endpoint
  if (args[0]?.includes?.('auth/login') && args[0]?.includes?.('401')) {
    return; // Don't log this error
  }
  originalError.apply(console, args);
};
```

**Note:** It's better to handle errors properly in your code rather than suppressing console errors globally.

---

*Last Updated: 2024*

