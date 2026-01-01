# Admin Login Credentials

## Default Admin Account

The admin user is **automatically created** when the server starts (via `database/seed.go`).

### Login Credentials

```
Username: admin
Email:    admin@playtz.com
Password: admin123
```

### User Details

- **First Name:** Admin
- **Last Name:** User
- **Role:** Admin (full system access)
- **Status:** Active
- **Password Change Required:** No

### Role Permissions

The admin role has full access to all features:
- ✅ Users (read, write, delete)
- ✅ Roles (read, write, delete)
- ✅ News (read, write, delete)
- ✅ Events (read, write, delete)
- ✅ Mixes (read, write, delete)
- ✅ Merchandise (read, write, delete)
- ✅ Orders (read, write, delete)
- ✅ Careers (read, write, delete)
- ✅ Rooms (read, write, delete)
- ✅ Admin Dashboard
- ✅ Admin Settings

---

## How It Works

1. **On Server Startup:** The `SeedAdmin()` function runs automatically
2. **Checks for Admin Role:** Creates "admin" role if it doesn't exist
3. **Checks for Admin User:** Creates admin user if it doesn't exist
4. **Password Reset:** If admin user exists, password is reset to `admin123`

This ensures the admin account is always available for login.

---

## Login Endpoints

### Production
```
POST https://playtzapi-production.up.railway.app/api/v1/auth/login
```

### Local Development
```
POST http://localhost:8080/api/v1/auth/login
```

### Request Body
```json
{
  "username": "admin",
  "password": "admin123"
}
```

### Response (Success)
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "uuid",
    "email": "admin@playtz.com",
    "username": "admin",
    "first_name": "Admin",
    "last_name": "User",
    "role_id": "uuid",
    "role_name": "Admin",
    "active": true
  }
}
```

---

## Security Notes

⚠️ **Important:**
- The default password `admin123` is for development/testing
- **Change the password in production** using the change password endpoint
- The admin user is automatically created - this is intentional for easy setup
- Consider disabling auto-creation in production if needed

---

## Testing Login

### Using cURL
```bash
curl -X POST https://playtzapi-production.up.railway.app/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c cookies.txt
```

### Using JavaScript/TypeScript
```typescript
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
console.log('Token:', data.token);
console.log('User:', data.user);
```

---

## Other Users

To see all users in the database, use the protected endpoint:

```bash
GET /api/v1/users
```

**Authentication Required:** Yes (JWT token)

---

*Last Updated: 2024*

