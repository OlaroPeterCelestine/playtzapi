# Admin Dashboard API Routes Reference

Complete API routes for building the admin dashboard. All routes use base URL: `/api/v1`

---

## 🔐 Authentication Routes

### Login
```javascript
POST /api/v1/auth/login
Body: { username: string, password: string }
Response: { success: boolean, session_id: string, user: User }
```

### Logout
```javascript
POST /api/v1/auth/logout
Response: { message: string }
```

### Get Current User
```javascript
GET /api/v1/auth/me
Auth: Required
Response: User
```

---

## 📊 Admin Dashboard

### Get Dashboard Data
```javascript
GET /api/v1/admin/dashboard
Auth: Required
Response: {
  user: User,
  stats: {
    total_users: number,
    total_news: number,
    total_events: number,
    total_merchandise: number,
    total_orders: number,
    pending_orders: number
  },
  recent_news: NewsArticle[],
  recent_events: Event[],
  recent_orders: Order[]
}
```

---

## 📰 News Management

### Get All News
```javascript
GET /api/v1/news
Response: NewsArticle[]
```

### Create News
```javascript
POST /api/v1/news
Body: {
  title: string,
  content: string,
  author: string,
  image?: string,
  published?: boolean
}
Response: NewsArticle
```

### Get News by ID
```javascript
GET /api/v1/news/:id
Response: NewsArticle
```

### Update News
```javascript
PUT /api/v1/news/:id
Body: {
  title: string,
  content: string,
  author: string,
  image?: string,
  published?: boolean
}
Response: NewsArticle
```

### Delete News
```javascript
DELETE /api/v1/news/:id
Response: { message: string }
```

---

## 🎉 Events Management

### Get All Events
```javascript
GET /api/v1/events
Response: Event[]
```

### Create Event
```javascript
POST /api/v1/events
Body: {
  title: string,
  description: string,
  date: string,        // "YYYY-MM-DD"
  time: string,        // "HH:MM:SS"
  location: string,
  image?: string,
  active?: boolean
}
Response: Event
```

### Get Event by ID
```javascript
GET /api/v1/events/:id
Response: Event
```

### Update Event
```javascript
PUT /api/v1/events/:id
Body: Same as Create Event
Response: Event
```

### Delete Event
```javascript
DELETE /api/v1/events/:id
Response: { message: string }
```

---

## 🛍️ Merchandise Management

### Get All Merchandise
```javascript
GET /api/v1/merch
Response: Merchandise[]
```

### Create Merchandise
```javascript
POST /api/v1/merch
Body: {
  name: string,
  description: string,
  price: number,
  image?: string,
  stock: number,
  active?: boolean
}
Response: Merchandise
```

### Get Merchandise by ID
```javascript
GET /api/v1/merch/:id
Response: Merchandise
```

### Update Merchandise
```javascript
PUT /api/v1/merch/:id
Body: Same as Create Merchandise
Response: Merchandise
```

### Delete Merchandise
```javascript
DELETE /api/v1/merch/:id
Response: { message: string }
```

---

## 💼 Careers Management

### Get All Careers
```javascript
GET /api/v1/careers
Response: Career[]
```

### Create Career
```javascript
POST /api/v1/careers
Body: {
  title: string,
  description: string,
  department: string,
  location: string,
  type: "full-time" | "part-time" | "contract",
  active?: boolean
}
Response: Career
```

### Get Career by ID
```javascript
GET /api/v1/careers/:id
Response: Career
```

### Update Career
```javascript
PUT /api/v1/careers/:id
Body: Same as Create Career
Response: Career
```

### Delete Career
```javascript
DELETE /api/v1/careers/:id
Response: { message: string }
```

---

## 🏠 Rooms Management

### Get All Rooms
```javascript
GET /api/v1/rooms
Response: Room[]
```

### Create Room
```javascript
POST /api/v1/rooms
Body: {
  name: string,
  genre: string,
  description: string,
  gradient: string,
  text_color: string,
  image?: string,
  active?: boolean
}
Response: Room
```

### Get Room by ID
```javascript
GET /api/v1/rooms/:id
Response: Room
```

### Update Room
```javascript
PUT /api/v1/rooms/:id
Body: Same as Create Room
Response: Room
```

### Delete Room
```javascript
DELETE /api/v1/rooms/:id
Response: { message: string }
```

---

## 🎵 Mixes Management

### Get All Mixes
```javascript
GET /api/v1/mixes
Response: Mix[]
```

### Create Mix
```javascript
POST /api/v1/mixes
Body: {
  room_id: string,
  title: string,
  artist: string,
  description: string,
  duration: string,
  tracks: number,
  color?: string,
  text_color?: string,
  border_color?: string,
  image?: string,
  audio_url?: string,
  active?: boolean
}
Response: Mix
```

### Get Mix by ID
```javascript
GET /api/v1/mixes/:id
Response: Mix
```

### Update Mix
```javascript
PUT /api/v1/mixes/:id
Body: Same as Create Mix
Response: Mix
```

### Delete Mix
```javascript
DELETE /api/v1/mixes/:id
Response: { message: string }
```

### Add Track to Mix
```javascript
POST /api/v1/mixes/:id/tracks
Body: {
  number: number,
  title: string,
  artist: string,
  duration: string,
  link: string,
  type: "audio" | "video"
}
Response: Track
```

### Add Multiple Tracks
```javascript
POST /api/v1/mixes/:id/tracks/bulk
Body: {
  tracks: Track[]
}
Response: Track[]
```

### Remove Track from Mix
```javascript
DELETE /api/v1/mixes/:id/tracks
Body: { track_id: string }
Response: { message: string }
```

---

## 👥 Users Management

### Get All Users
```javascript
GET /api/v1/users
Response: User[]
```

### Create User
```javascript
POST /api/v1/users
Body: {
  email: string,
  username: string,
  password: string,
  first_name: string,
  last_name: string,
  role_id: string
}
Response: User
```

### Get User by ID
```javascript
GET /api/v1/users/:id
Response: User
```

### Update User
```javascript
PUT /api/v1/users/:id
Body: {
  email: string,
  username: string,
  first_name: string,
  last_name: string,
  role_id: string
}
Response: User
```

### Delete User
```javascript
DELETE /api/v1/users/:id
Response: { message: string }
```

### Update User Role
```javascript
PUT /api/v1/users/:id/role
Body: { role_id: string }
Response: { message: string, role_id: string }
```

---

## 🔑 Roles Management

### Get All Roles
```javascript
GET /api/v1/roles
Response: Role[]
```

### Create Role
```javascript
POST /api/v1/roles
Body: {
  name: string,
  description: string,
  permissions: string[],
  active?: boolean
}
Response: Role
```

### Get Role by ID
```javascript
GET /api/v1/roles/:id
Response: Role
```

### Update Role
```javascript
PUT /api/v1/roles/:id
Body: Same as Create Role
Response: Role
```

### Delete Role
```javascript
DELETE /api/v1/roles/:id
Response: { message: string }
```

---

## 📦 Orders Management

### Get All Orders
```javascript
GET /api/v1/orders
Query: ?user_id=string (optional)
Response: Order[]
```

### Get Order by ID
```javascript
GET /api/v1/orders/:id
Response: Order
```

### Update Order Status
```javascript
PUT /api/v1/orders/:id/status
Body: {
  status: "pending" | "processing" | "shipped" | "delivered" | "cancelled"
}
Response: { message: string, status: string }
```

---

## 🛒 Cart Management

### Get Cart
```javascript
GET /api/v1/cart
Response: Cart
```

### Add to Cart
```javascript
POST /api/v1/cart/add
Body: {
  merchandise_id: string,
  quantity: number
}
Response: CartItem
```

### Update Cart Item
```javascript
PUT /api/v1/cart/update
Body: {
  cart_item_id: string,
  quantity: number
}
Response: CartItem
```

### Remove from Cart
```javascript
DELETE /api/v1/cart/remove
Body: { cart_item_id: string }
Response: { message: string }
```

### Clear Cart
```javascript
DELETE /api/v1/cart/clear
Response: { message: string }
```

---

## 💳 Checkout

### Create Order (Checkout)
```javascript
POST /api/v1/checkout
Body: {
  cart_id: string,
  user_id?: string,
  shipping_address: {
    full_name: string,
    email: string,
    phone: string,
    address: string,
    city: string,
    state: string,
    postal_code: string,
    country: string
  },
  payment_method: string
}
Response: Order
```

---

## 📤 Image Upload

### Upload Single Image
```javascript
POST /api/v1/upload
Content-Type: multipart/form-data
Form Data:
  - image: File
  - folder?: string (optional, default: "playtz")
Response: {
  url: string,
  public_id: string,
  secure_url: string
}
```

### Upload Multiple Images
```javascript
POST /api/v1/upload/multiple
Content-Type: multipart/form-data
Form Data:
  - images: File[] (multiple files)
  - folder?: string (optional)
Response: {
  images: UploadResponse[],
  count: number
}
```

---

## 📋 Data Types

### User
```typescript
{
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role_id: string;
  role_name?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### NewsArticle
```typescript
{
  id: string;
  title: string;
  content: string;
  author: string;
  image?: string;
  published: boolean;
  created_at: string;
  updated_at: string;
}
```

### Event
```typescript
{
  id: string;
  title: string;
  description: string;
  date: string;
  time: string;
  location: string;
  image?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Merchandise
```typescript
{
  id: string;
  name: string;
  description: string;
  price: number;
  image?: string;
  stock: number;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Career
```typescript
{
  id: string;
  title: string;
  description: string;
  department: string;
  location: string;
  type: "full-time" | "part-time" | "contract";
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Room
```typescript
{
  id: string;
  name: string;
  genre: string;
  description: string;
  gradient: string;
  text_color: string;
  image?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Mix
```typescript
{
  id: string;
  room_id: string;
  title: string;
  artist: string;
  description: string;
  duration: string;
  tracks: number;
  color?: string;
  text_color?: string;
  border_color?: string;
  image?: string;
  audio_url?: string;
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

### Order
```typescript
{
  id: string;
  user_id?: string;
  items: OrderItem[];
  subtotal: number;
  shipping: number;
  tax: number;
  total: number;
  status: "pending" | "processing" | "shipped" | "delivered" | "cancelled";
  payment_method?: string;
  shipping_address: ShippingAddress;
  created_at: string;
  updated_at: string;
}
```

### Role
```typescript
{
  id: string;
  name: string;
  description: string;
  permissions: string[];
  active: boolean;
  created_at: string;
  updated_at: string;
}
```

---

## 🔧 JavaScript/TypeScript Usage Examples

### Fetch with Authentication
```javascript
// Login
const login = async (username, password) => {
  const response = await fetch('/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ username, password })
  });
  return response.json();
};

// Authenticated Request
const getDashboard = async () => {
  const response = await fetch('/api/v1/admin/dashboard', {
    method: 'GET',
    credentials: 'include',
    headers: { 'Content-Type': 'application/json' }
  });
  return response.json();
};

// Create News
const createNews = async (newsData) => {
  const response = await fetch('/api/v1/news', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(newsData)
  });
  return response.json();
};

// Upload Image
const uploadImage = async (file, folder = 'playtz') => {
  const formData = new FormData();
  formData.append('image', file);
  formData.append('folder', folder);
  
  const response = await fetch('/api/v1/upload', {
    method: 'POST',
    credentials: 'include',
    body: formData
  });
  return response.json();
};

// Update Order Status
const updateOrderStatus = async (orderId, status) => {
  const response = await fetch(`/api/v1/orders/${orderId}/status`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ status })
  });
  return response.json();
};
```

### Axios Usage
```javascript
import axios from 'axios';

// Configure axios to include credentials
axios.defaults.withCredentials = true;

// Login
const login = (username, password) => {
  return axios.post('/api/v1/auth/login', { username, password });
};

// Get Dashboard
const getDashboard = () => {
  return axios.get('/api/v1/admin/dashboard');
};

// Create News
const createNews = (data) => {
  return axios.post('/api/v1/news', data);
};

// Upload Image
const uploadImage = (file, folder = 'playtz') => {
  const formData = new FormData();
  formData.append('image', file);
  formData.append('folder', folder);
  return axios.post('/api/v1/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  });
};
```

---

## 📝 Quick Reference Table

| Resource | GET All | GET One | CREATE | UPDATE | DELETE |
|----------|---------|---------|--------|--------|--------|
| News | `/api/v1/news` | `/api/v1/news/:id` | `POST /api/v1/news` | `PUT /api/v1/news/:id` | `DELETE /api/v1/news/:id` |
| Events | `/api/v1/events` | `/api/v1/events/:id` | `POST /api/v1/events` | `PUT /api/v1/events/:id` | `DELETE /api/v1/events/:id` |
| Merchandise | `/api/v1/merch` | `/api/v1/merch/:id` | `POST /api/v1/merch` | `PUT /api/v1/merch/:id` | `DELETE /api/v1/merch/:id` |
| Careers | `/api/v1/careers` | `/api/v1/careers/:id` | `POST /api/v1/careers` | `PUT /api/v1/careers/:id` | `DELETE /api/v1/careers/:id` |
| Rooms | `/api/v1/rooms` | `/api/v1/rooms/:id` | `POST /api/v1/rooms` | `PUT /api/v1/rooms/:id` | `DELETE /api/v1/rooms/:id` |
| Mixes | `/api/v1/mixes` | `/api/v1/mixes/:id` | `POST /api/v1/mixes` | `PUT /api/v1/mixes/:id` | `DELETE /api/v1/mixes/:id` |
| Users | `/api/v1/users` | `/api/v1/users/:id` | `POST /api/v1/users` | `PUT /api/v1/users/:id` | `DELETE /api/v1/users/:id` |
| Roles | `/api/v1/roles` | `/api/v1/roles/:id` | `POST /api/v1/roles` | `PUT /api/v1/roles/:id` | `DELETE /api/v1/roles/:id` |
| Orders | `/api/v1/orders` | `/api/v1/orders/:id` | `POST /api/v1/checkout` | `PUT /api/v1/orders/:id/status` | - |

---

## ⚠️ Important Notes

1. **Authentication**: Most routes require authentication. Include `credentials: 'include'` in fetch requests.

2. **Session Timeout**: Sessions expire after 10 minutes of inactivity. Handle 401 errors gracefully.

3. **Base URL**: Use relative URLs (`/api/v1/...`) or set base URL:
   ```javascript
   const API_BASE = 'http://localhost:8080/api/v1';
   // or for production
   const API_BASE = 'https://your-railway-url.railway.app/api/v1';
   ```

4. **Error Handling**: Always check response status:
   ```javascript
   if (!response.ok) {
     const error = await response.json();
     throw new Error(error.error || 'Request failed');
   }
   ```

5. **File Uploads**: Use `FormData` for image uploads, not JSON.

---

## 🎯 Dashboard Implementation Tips

1. **Use the dashboard endpoint** for initial data load:
   ```javascript
   GET /api/v1/admin/dashboard
   ```
   This returns user info, stats, and recent items based on role.

2. **Poll for updates** every 30-60 seconds for real-time data.

3. **Handle role-based access** - Different roles see different data in dashboard response.

4. **Use the stats** for dashboard cards/widgets.

5. **Recent items** arrays are already filtered by role - use them directly.

