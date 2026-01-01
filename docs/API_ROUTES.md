# Playtz API Routes Documentation

**Base URL:** `http://localhost:8080` (or your Railway deployment URL)  
**API Base:** `/api/v1`

---

## Authentication Routes

### Login
- **Endpoint:** `POST /api/v1/auth/login`
- **Auth Required:** No
- **Request Body:**
```json
{
  "username": "string",  // Username or email
  "password": "string"
}
```
- **Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "session_id": "string",
  "user": {
    "id": "string",
    "email": "string",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "role_id": "string",
    "role_name": "string",
    "active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```
- **Error Responses:**
  - `400` - Invalid request body
  - `401` - Invalid credentials
  - `403` - Account inactive

### Logout
- **Endpoint:** `POST /api/v1/auth/logout`
- **Auth Required:** No (but recommended)
- **Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

### Get Current User
- **Endpoint:** `GET /api/v1/auth/me`
- **Auth Required:** Yes
- **Response (200):**
```json
{
  "id": "string",
  "email": "string",
  "username": "string",
  "first_name": "string",
  "last_name": "string",
  "role_id": "string",
  "role_name": "string",
  "active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

## Admin Routes

### Get Admin Dashboard
- **Endpoint:** `GET /api/v1/admin/dashboard`
- **Auth Required:** Yes
- **Response (200):**
```json
{
  "user": {
    "id": "string",
    "email": "string",
    "username": "string",
    "first_name": "string",
    "last_name": "string",
    "role_id": "string",
    "role_name": "string",
    "active": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "stats": {
    "total_users": 0,
    "total_news": 0,
    "total_events": 0,
    "total_merchandise": 0,
    "total_orders": 0,
    "pending_orders": 0
  },
  "recent_news": [...],      // Only for Admin/Editor roles
  "recent_events": [...],    // Only for Admin/Editor roles
  "recent_orders": [...]     // Only for Admin/Manager roles
}
```

---

## News Routes

### Get All News
- **Endpoint:** `GET /api/v1/news`
- **Auth Required:** No
- **Query Parameters:** None
- **Response (200):** Array of news articles

### Create News
- **Endpoint:** `POST /api/v1/news`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "title": "string",
  "content": "string",
  "author": "string",
  "image": "string (URL)",
  "published": false
}
```
- **Response (201):** Created news article

### Get News by ID
- **Endpoint:** `GET /api/v1/news/:id`
- **Auth Required:** No
- **Response (200):** News article object
- **Error:** `404` - News article not found

### Update News
- **Endpoint:** `PUT /api/v1/news/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create News
- **Response (200):** Updated news article

### Delete News
- **Endpoint:** `DELETE /api/v1/news/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "News article deleted successfully"
}
```

---

## Events Routes

### Get All Events
- **Endpoint:** `GET /api/v1/events`
- **Auth Required:** No
- **Response (200):** Array of events

### Create Event
- **Endpoint:** `POST /api/v1/events`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "date": "2024-12-31",
  "time": "18:00:00",
  "location": "string",
  "image": "string (URL)",
  "active": true
}
```
- **Response (201):** Created event

### Get Event by ID
- **Endpoint:** `GET /api/v1/events/:id`
- **Auth Required:** No
- **Response (200):** Event object

### Update Event
- **Endpoint:** `PUT /api/v1/events/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Event
- **Response (200):** Updated event

### Delete Event
- **Endpoint:** `DELETE /api/v1/events/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Event deleted successfully"
}
```

---

## Merchandise Routes

### Get All Merchandise
- **Endpoint:** `GET /api/v1/merch`
- **Auth Required:** No
- **Response (200):** Array of merchandise items

### Create Merchandise
- **Endpoint:** `POST /api/v1/merch`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "name": "string",
  "description": "string",
  "price": 29.99,
  "image": "string (URL)",
  "stock": 100,
  "active": true
}
```
- **Response (201):** Created merchandise item

### Get Merchandise by ID
- **Endpoint:** `GET /api/v1/merch/:id`
- **Auth Required:** No
- **Response (200):** Merchandise item object

### Update Merchandise
- **Endpoint:** `PUT /api/v1/merch/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Merchandise
- **Response (200):** Updated merchandise item

### Delete Merchandise
- **Endpoint:** `DELETE /api/v1/merch/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Merchandise item deleted successfully"
}
```

---

## Careers Routes

### Get All Careers
- **Endpoint:** `GET /api/v1/careers`
- **Auth Required:** No
- **Response (200):** Array of career listings

### Create Career
- **Endpoint:** `POST /api/v1/careers`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "title": "string",
  "description": "string",
  "department": "string",
  "location": "string",
  "type": "full-time",  // "full-time", "part-time", "contract"
  "active": true
}
```
- **Response (201):** Created career listing

### Get Career by ID
- **Endpoint:** `GET /api/v1/careers/:id`
- **Auth Required:** No
- **Response (200):** Career listing object

### Update Career
- **Endpoint:** `PUT /api/v1/careers/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Career
- **Response (200):** Updated career listing

### Delete Career
- **Endpoint:** `DELETE /api/v1/careers/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Career listing deleted successfully"
}
```

---

## Shopping Cart Routes

### Get Cart
- **Endpoint:** `GET /api/v1/cart`
- **Auth Required:** No
- **Response (200):** Cart object with items

### Add to Cart
- **Endpoint:** `POST /api/v1/cart/add`
- **Auth Required:** No
- **Request Body:**
```json
{
  "merchandise_id": "string",
  "quantity": 2
}
```
- **Response (201):** Cart item added

### Update Cart Item
- **Endpoint:** `PUT /api/v1/cart/update`
- **Auth Required:** No
- **Request Body:**
```json
{
  "cart_item_id": "string",
  "quantity": 3
}
```
- **Response (200):** Updated cart item

### Remove from Cart
- **Endpoint:** `DELETE /api/v1/cart/remove`
- **Auth Required:** No
- **Request Body:**
```json
{
  "cart_item_id": "string"
}
```
- **Response (200):**
```json
{
  "message": "Item removed from cart"
}
```

### Clear Cart
- **Endpoint:** `DELETE /api/v1/cart/clear`
- **Auth Required:** No
- **Response (200):**
```json
{
  "message": "Cart cleared successfully"
}
```

---

## Checkout & Orders Routes

### Create Order (Checkout)
- **Endpoint:** `POST /api/v1/checkout`
- **Auth Required:** No
- **Request Body:**
```json
{
  "cart_id": "string",
  "user_id": "string (optional)",
  "shipping_address": {
    "full_name": "string",
    "email": "string",
    "phone": "string",
    "address": "string",
    "city": "string",
    "state": "string",
    "postal_code": "string",
    "country": "string"
  },
  "payment_method": "string"
}
```
- **Response (201):** Created order object

### Get All Orders
- **Endpoint:** `GET /api/v1/orders`
- **Auth Required:** No (should be protected in production)
- **Query Parameters:**
  - `user_id` (optional) - Filter by user ID
- **Response (200):** Array of orders

### Get Order by ID
- **Endpoint:** `GET /api/v1/orders/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):** Order object with items

### Update Order Status
- **Endpoint:** `PUT /api/v1/orders/:id/status`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "status": "pending"  // "pending", "processing", "shipped", "delivered", "cancelled"
}
```
- **Response (200):**
```json
{
  "message": "Order status updated",
  "status": "processing"
}
```

---

## Rooms Routes

### Get All Rooms
- **Endpoint:** `GET /api/v1/rooms`
- **Auth Required:** No
- **Response (200):** Array of rooms

### Create Room
- **Endpoint:** `POST /api/v1/rooms`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "name": "string",
  "genre": "string",
  "description": "string",
  "gradient": "string",
  "text_color": "#ffffff",
  "image": "string (URL)",
  "active": true
}
```
- **Response (201):** Created room

### Get Room by ID
- **Endpoint:** `GET /api/v1/rooms/:id`
- **Auth Required:** No
- **Response (200):** Room object

### Update Room
- **Endpoint:** `PUT /api/v1/rooms/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Room
- **Response (200):** Updated room

### Delete Room
- **Endpoint:** `DELETE /api/v1/rooms/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Room deleted successfully"
}
```

---

## Mixes Routes

### Get All Mixes
- **Endpoint:** `GET /api/v1/mixes`
- **Auth Required:** No
- **Response (200):** Array of mixes

### Create Mix
- **Endpoint:** `POST /api/v1/mixes`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "room_id": "string",
  "title": "string",
  "artist": "string",
  "description": "string",
  "duration": "60:00",
  "tracks": 10,
  "color": "string",
  "text_color": "#ffffff",
  "border_color": "#000000",
  "image": "string (URL)",
  "audio_url": "string (URL)",
  "active": true
}
```
- **Response (201):** Created mix

### Get Mix by ID
- **Endpoint:** `GET /api/v1/mixes/:id`
- **Auth Required:** No
- **Response (200):** Mix object

### Update Mix
- **Endpoint:** `PUT /api/v1/mixes/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Mix
- **Response (200):** Updated mix

### Delete Mix
- **Endpoint:** `DELETE /api/v1/mixes/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Mix deleted successfully"
}
```

### Add Track to Mix
- **Endpoint:** `POST /api/v1/mixes/:id/tracks`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "number": 1,
  "title": "string",
  "artist": "string",
  "duration": "3:45",
  "link": "string (URL)",
  "type": "audio"  // "audio" or "video"
}
```
- **Response (201):** Created track

### Add Multiple Tracks to Mix
- **Endpoint:** `POST /api/v1/mixes/:id/tracks/bulk`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "tracks": [
    {
      "number": 1,
      "title": "string",
      "artist": "string",
      "duration": "3:45",
      "link": "string (URL)",
      "type": "audio"
    }
  ]
}
```
- **Response (201):** Array of created tracks

### Remove Track from Mix
- **Endpoint:** `DELETE /api/v1/mixes/:id/tracks`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "track_id": "string"
}
```
- **Response (200):**
```json
{
  "message": "Track removed successfully"
}
```

---

## Upload Routes

### Upload Single Image
- **Endpoint:** `POST /api/v1/upload`
- **Auth Required:** No (should be protected in production)
- **Content-Type:** `multipart/form-data`
- **Form Data:**
  - `image` (file) - Image file (jpg, jpeg, png, gif, webp)
  - `folder` (optional) - Cloudinary folder (default: "playtz")
- **Response (200):**
```json
{
  "url": "http://res.cloudinary.com/...",
  "public_id": "folder/filename",
  "secure_url": "https://res.cloudinary.com/..."
}
```
- **Error Responses:**
  - `400` - Invalid file type or no file provided
  - `500` - Upload failed

### Upload Multiple Images
- **Endpoint:** `POST /api/v1/upload/multiple`
- **Auth Required:** No (should be protected in production)
- **Content-Type:** `multipart/form-data`
- **Form Data:**
  - `images` (files) - Multiple image files
  - `folder` (optional) - Cloudinary folder (default: "playtz")
- **Response (200):**
```json
{
  "images": [
    {
      "url": "http://res.cloudinary.com/...",
      "public_id": "folder/filename",
      "secure_url": "https://res.cloudinary.com/..."
    }
  ],
  "count": 2
}
```

---

## Users Routes

### Get All Users
- **Endpoint:** `GET /api/v1/users`
- **Auth Required:** No (should be protected in production)
- **Response (200):** Array of users (without passwords)

### Create User
- **Endpoint:** `POST /api/v1/users`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe",
  "role_id": "string"
}
```
- **Response (201):** Created user (without password)

### Get User by ID
- **Endpoint:** `GET /api/v1/users/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):** User object

### Update User
- **Endpoint:** `PUT /api/v1/users/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "first_name": "John",
  "last_name": "Doe",
  "role_id": "string"
}
```
- **Response (200):** Updated user

### Delete User
- **Endpoint:** `DELETE /api/v1/users/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

### Update User Role
- **Endpoint:** `PUT /api/v1/users/:id/role`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "role_id": "string"
}
```
- **Response (200):**
```json
{
  "message": "User role updated",
  "role_id": "string"
}
```

---

## Roles Routes

### Get All Roles
- **Endpoint:** `GET /api/v1/roles`
- **Auth Required:** No (should be protected in production)
- **Response (200):** Array of roles

### Create Role
- **Endpoint:** `POST /api/v1/roles`
- **Auth Required:** No (should be protected in production)
- **Request Body:**
```json
{
  "name": "string",
  "description": "string",
  "permissions": ["read", "write", "delete"],
  "active": true
}
```
- **Response (201):** Created role

### Get Role by ID
- **Endpoint:** `GET /api/v1/roles/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):** Role object

### Update Role
- **Endpoint:** `PUT /api/v1/roles/:id`
- **Auth Required:** No (should be protected in production)
- **Request Body:** Same as Create Role
- **Response (200):** Updated role

### Delete Role
- **Endpoint:** `DELETE /api/v1/roles/:id`
- **Auth Required:** No (should be protected in production)
- **Response (200):**
```json
{
  "message": "Role deleted successfully"
}
```
- **Error:** `400` - Role is assigned to users

---

## Health & Debug Routes

### Health Check
- **Endpoint:** `GET /health`
- **Auth Required:** No
- **Response (200):**
```json
{
  "status": "ok"
}
```

### Cloudinary Debug
- **Endpoint:** `GET /debug/cloudinary`
- **Auth Required:** No
- **Response (200):** Cloudinary configuration status

---

## Authentication

### Session-based Authentication
- Sessions are stored in HTTP-only cookies
- Session ID cookie name: `session_id`
- Session expires after 10 minutes of inactivity
- Include credentials in requests: `credentials: 'include'` (fetch) or `withCredentials: true` (axios)

### Example Request with Authentication
```javascript
// Using fetch
fetch('/api/v1/admin/dashboard', {
  method: 'GET',
  credentials: 'include',  // Important for cookies
  headers: {
    'Content-Type': 'application/json'
  }
})

// Using axios
axios.get('/api/v1/admin/dashboard', {
  withCredentials: true  // Important for cookies
})
```

---

## Error Responses

All endpoints may return these error responses:

- **400 Bad Request** - Invalid request body or parameters
```json
{
  "error": "Invalid request body"
}
```

- **401 Unauthorized** - Authentication required or invalid
```json
{
  "error": "Authentication required"
}
```

- **403 Forbidden** - Insufficient permissions
```json
{
  "error": "Insufficient permissions"
}
```

- **404 Not Found** - Resource not found
```json
{
  "error": "Resource not found"
}
```

- **500 Internal Server Error** - Server error
```json
{
  "error": "Failed to process request"
}
```

---

## Notes

1. **Production Security:** Many endpoints marked as "should be protected in production" are currently public. Add authentication middleware before deploying to production.

2. **CORS:** CORS is configured to allow all origins. Restrict this in production.

3. **Session Timeout:** Sessions expire after 10 minutes of inactivity. The frontend should handle session expiration gracefully.

4. **File Uploads:** Image uploads are automatically converted to WebP format and optimized before uploading to Cloudinary.

5. **IDs:** All IDs are UUIDs (strings).

6. **Dates:** All dates are in RFC3339 format (ISO 8601).

