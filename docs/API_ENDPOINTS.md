# Complete API Endpoints List

**Base URL:** `http://localhost:8080`  
**API Base:** `http://localhost:8080/api/v1`

---

## Health Check
- `GET /health` - Server health check

---

## Web Pages (HTML)
- `GET /admin/login` - Admin login page
- `GET /admin/dashboard` - Admin dashboard (requires authentication)

---

## Authentication Endpoints

### Public Endpoints
- `POST /api/v1/auth/login` - Login user
  - **Request Body:**
    ```json
    {
      "username": "admin",
      "password": "admin123"
    }
    ```
  - **Response:** Returns session ID and user information
  - **Sets Cookie:** `session_id` (HTTP-only, 10 minutes)

- `POST /api/v1/auth/logout` - Logout user (requires authentication)
  - **Clears:** Session cookie

### Protected Endpoints (Require Authentication)
- `GET /api/v1/auth/me` - Get current authenticated user
  - **Requires:** Valid session cookie

- `POST /api/v1/auth/change-password` - Change user password
  - **Requires:** Valid session cookie
  - **Request Body:**
    ```json
    {
      "current_password": "oldpassword",
      "new_password": "newpassword123"
    }
    ```
  - **Note:** New password must be at least 6 characters

---

## Admin Endpoints (Protected)

All admin endpoints require authentication via session cookie.

- `GET /api/v1/admin/dashboard` - Get admin dashboard data
  - **Returns:** User info, statistics, recent news/events/orders
  - **Requires:** Valid session cookie

---

## News Endpoints
- `GET /api/v1/news` - Get all news articles
- `POST /api/v1/news` - Create news article
- `GET /api/v1/news/:id` - Get news article by ID
- `PUT /api/v1/news/:id` - Update news article
- `DELETE /api/v1/news/:id` - Delete news article

---

## Events Endpoints
- `GET /api/v1/events` - Get all events
- `POST /api/v1/events` - Create event
- `GET /api/v1/events/:id` - Get event by ID
- `PUT /api/v1/events/:id` - Update event
- `DELETE /api/v1/events/:id` - Delete event

---

## Merchandise Endpoints
- `GET /api/v1/merch` - Get all merchandise
- `POST /api/v1/merch` - Create merchandise item
- `GET /api/v1/merch/:id` - Get merchandise by ID
- `PUT /api/v1/merch/:id` - Update merchandise
- `DELETE /api/v1/merch/:id` - Delete merchandise

---

## Careers Endpoints
- `GET /api/v1/careers` - Get all career listings
- `POST /api/v1/careers` - Create career listing
- `GET /api/v1/careers/:id` - Get career by ID
- `PUT /api/v1/careers/:id` - Update career listing
- `DELETE /api/v1/careers/:id` - Delete career listing

---

## Shopping Cart Endpoints
- `GET /api/v1/cart` - Get cart
- `POST /api/v1/cart/add` - Add item to cart
- `PUT /api/v1/cart/update` - Update cart item
- `DELETE /api/v1/cart/remove` - Remove item from cart
- `DELETE /api/v1/cart/clear` - Clear entire cart

---

## Checkout & Orders Endpoints
- `POST /api/v1/checkout` - Create order from cart
- `GET /api/v1/orders` - Get all orders
- `GET /api/v1/orders/:id` - Get order by ID
- `PUT /api/v1/orders/:id/status` - Update order status

---

## Rooms Endpoints
- `GET /api/v1/rooms` - Get all rooms
- `POST /api/v1/rooms` - Create room
- `GET /api/v1/rooms/:id` - Get room by ID
- `PUT /api/v1/rooms/:id` - Update room
- `DELETE /api/v1/rooms/:id` - Delete room

---

## Mixes Endpoints
- `GET /api/v1/mixes` - Get all mixes
- `POST /api/v1/mixes` - Create mix
- `GET /api/v1/mixes/:id` - Get mix by ID
- `PUT /api/v1/mixes/:id` - Update mix
- `DELETE /api/v1/mixes/:id` - Delete mix
- `POST /api/v1/mixes/:id/tracks` - Add track to mix
- `POST /api/v1/mixes/:id/tracks/bulk` - Add multiple tracks to mix
- `DELETE /api/v1/mixes/:id/tracks` - Remove track from mix

---

## Upload Endpoints
- `POST /api/v1/upload` - Upload single image to Cloudinary
  - **Content-Type:** `multipart/form-data`
  - **Fields:** `file`, `folder` (optional)
  - **Returns:** Image URL

- `POST /api/v1/upload/multiple` - Upload multiple images to Cloudinary
  - **Content-Type:** `multipart/form-data`
  - **Fields:** `files[]`, `folder` (optional)
  - **Returns:** Array of image URLs

---

## Users Endpoints
- `GET /api/v1/users` - Get all users
- `POST /api/v1/users` - Create user
  - **Request Body:**
    ```json
    {
      "email": "user@example.com",
      "username": "username",
      "password": "optional", // If omitted, generates default password
      "first_name": "First",
      "last_name": "Last",
      "role_id": "role-uuid"
    }
    ```
  - **Response:** If password omitted, includes `default_password` in response
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `PUT /api/v1/users/:id/role` - Update user role
  - **Request Body:**
    ```json
    {
      "role_id": "new-role-uuid"
    }
    ```

---

## Roles Endpoints
- `GET /api/v1/roles` - Get all roles
- `POST /api/v1/roles` - Create role
  - **Request Body:**
    ```json
    {
      "name": "Role Name",
      "description": "Role description",
      "permissions": ["read", "write", "delete"]
    }
    ```
- `GET /api/v1/roles/:id` - Get role by ID
- `PUT /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role

---

## Debug Endpoints (Development Only)
- `GET /debug/cloudinary` - Check Cloudinary configuration
  - **Returns:** Cloudinary environment variables (masked)

---

## Authentication & Session Management

### How Authentication Works

1. **Login:** Send POST request to `/api/v1/auth/login` with username and password
2. **Session Cookie:** Server sets HTTP-only cookie `session_id` (valid for 10 minutes)
3. **Authenticated Requests:** Include the session cookie in subsequent requests
4. **Session Expiry:** Sessions expire after 10 minutes of inactivity
5. **Logout:** Send POST to `/api/v1/auth/logout` to clear session

### Example Login Flow

```bash
# 1. Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt

# 2. Use authenticated endpoint
curl -X GET http://localhost:8080/api/v1/auth/me \
  -b cookies.txt

# 3. Access admin dashboard
curl -X GET http://localhost:8080/api/v1/admin/dashboard \
  -b cookies.txt

# 4. Logout
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

### Protected Endpoints

The following endpoints require authentication:
- `GET /api/v1/auth/me`
- `POST /api/v1/auth/change-password`
- `POST /api/v1/auth/logout`
- `GET /api/v1/admin/dashboard`

---

## Summary

**Total API Endpoints: 55+**

### Endpoint Categories:
- **Health Check:** 1
- **Web Pages:** 2
- **Authentication:** 4
- **Admin:** 1
- **News:** 5
- **Events:** 5
- **Merchandise:** 5
- **Careers:** 5
- **Shopping Cart:** 5
- **Orders:** 4
- **Rooms:** 5
- **Mixes:** 7
- **Upload:** 2
- **Users:** 6
- **Roles:** 5
- **Debug:** 1

**Base URL for API:** `http://localhost:8080/api/v1`  
**Base URL for Web:** `http://localhost:8080`

---

## Testing

See [TESTING_LOGIN.md](./TESTING_LOGIN.md) for detailed testing instructions.

**Quick Test:**
```bash
./scripts/test_login.sh
```

**Test Credentials:**
- Username: `admin`
- Password: `admin123`
