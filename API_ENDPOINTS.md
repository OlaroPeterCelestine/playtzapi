# Complete API Endpoints List

**Base URL:** `http://localhost:8080`
**API Base:** `http://localhost:8080/api/v1`

---

## Health Check
- `GET http://localhost:8080/health`

---

## News Endpoints
- `GET http://localhost:8080/api/v1/news` - Get all news articles
- `POST http://localhost:8080/api/v1/news` - Create news article
- `GET http://localhost:8080/api/v1/news/:id` - Get news article by ID
- `PUT http://localhost:8080/api/v1/news/:id` - Update news article
- `DELETE http://localhost:8080/api/v1/news/:id` - Delete news article

---

## Events Endpoints
- `GET http://localhost:8080/api/v1/events` - Get all events
- `POST http://localhost:8080/api/v1/events` - Create event
- `GET http://localhost:8080/api/v1/events/:id` - Get event by ID
- `PUT http://localhost:8080/api/v1/events/:id` - Update event
- `DELETE http://localhost:8080/api/v1/events/:id` - Delete event

---

## Merchandise Endpoints
- `GET http://localhost:8080/api/v1/merch` - Get all merchandise
- `POST http://localhost:8080/api/v1/merch` - Create merchandise item
- `GET http://localhost:8080/api/v1/merch/:id` - Get merchandise by ID
- `PUT http://localhost:8080/api/v1/merch/:id` - Update merchandise
- `DELETE http://localhost:8080/api/v1/merch/:id` - Delete merchandise

---

## Careers Endpoints
- `GET http://localhost:8080/api/v1/careers` - Get all career listings
- `POST http://localhost:8080/api/v1/careers` - Create career listing
- `GET http://localhost:8080/api/v1/careers/:id` - Get career by ID
- `PUT http://localhost:8080/api/v1/careers/:id` - Update career listing
- `DELETE http://localhost:8080/api/v1/careers/:id` - Delete career listing

---

## Shopping Cart Endpoints
- `GET http://localhost:8080/api/v1/cart` - Get cart
- `POST http://localhost:8080/api/v1/cart/add` - Add item to cart
- `PUT http://localhost:8080/api/v1/cart/update` - Update cart item
- `DELETE http://localhost:8080/api/v1/cart/remove` - Remove item from cart
- `DELETE http://localhost:8080/api/v1/cart/clear` - Clear entire cart

---

## Checkout & Orders Endpoints
- `POST http://localhost:8080/api/v1/checkout` - Create order from cart
- `GET http://localhost:8080/api/v1/orders` - Get all orders
- `GET http://localhost:8080/api/v1/orders/:id` - Get order by ID
- `PUT http://localhost:8080/api/v1/orders/:id/status` - Update order status

---

## Rooms Endpoints
- `GET http://localhost:8080/api/v1/rooms` - Get all rooms
- `POST http://localhost:8080/api/v1/rooms` - Create room
- `GET http://localhost:8080/api/v1/rooms/:id` - Get room by ID
- `PUT http://localhost:8080/api/v1/rooms/:id` - Update room
- `DELETE http://localhost:8080/api/v1/rooms/:id` - Delete room

---

## Mixes Endpoints
- `GET http://localhost:8080/api/v1/mixes` - Get all mixes
- `POST http://localhost:8080/api/v1/mixes` - Create mix
- `GET http://localhost:8080/api/v1/mixes/:id` - Get mix by ID
- `PUT http://localhost:8080/api/v1/mixes/:id` - Update mix
- `DELETE http://localhost:8080/api/v1/mixes/:id` - Delete mix
- `POST http://localhost:8080/api/v1/mixes/:id/tracks` - Add track to mix
- `POST http://localhost:8080/api/v1/mixes/:id/tracks/bulk` - Add multiple tracks to mix
- `DELETE http://localhost:8080/api/v1/mixes/:id/tracks` - Remove track from mix

---

## Upload Endpoints
- `POST http://localhost:8080/api/v1/upload` - Upload single image to Cloudinary
- `POST http://localhost:8080/api/v1/upload/multiple` - Upload multiple images to Cloudinary

---

## Users Endpoints
- `GET http://localhost:8080/api/v1/users` - Get all users
- `POST http://localhost:8080/api/v1/users` - Create user
- `GET http://localhost:8080/api/v1/users/:id` - Get user by ID
- `PUT http://localhost:8080/api/v1/users/:id` - Update user
- `DELETE http://localhost:8080/api/v1/users/:id` - Delete user
- `PUT http://localhost:8080/api/v1/users/:id/role` - Update user role

---

## Roles Endpoints
- `GET http://localhost:8080/api/v1/roles` - Get all roles
- `POST http://localhost:8080/api/v1/roles` - Create role
- `GET http://localhost:8080/api/v1/roles/:id` - Get role by ID
- `PUT http://localhost:8080/api/v1/roles/:id` - Update role
- `DELETE http://localhost:8080/api/v1/roles/:id` - Delete role

---

## Debug Endpoints (Development Only)
- `GET http://localhost:8080/debug/cloudinary` - Check Cloudinary configuration

---

## Summary
**Total Endpoints: 50+**

**Base URL for Website:** `http://localhost:8080/api/v1`

