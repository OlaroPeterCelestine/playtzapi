# Frontend Integration Guide - JWT Authentication

## Overview

This guide explains how to integrate the Playtz API with your frontend admin dashboard using JWT token authentication.

**Base URL:** `https://playtzapi-production.up.railway.app/api/v1`

---

## Authentication Flow

### 1. Login Flow

1. User submits credentials
2. API returns JWT token in response + sets cookie
3. Store token in memory/localStorage (optional)
4. Use token for all subsequent requests

### 2. Token Storage

- **Cookie (Recommended):** Token is automatically stored in `token` cookie (HttpOnly, Secure)
- **Authorization Header:** Can also use `Authorization: Bearer <token>` header
- **LocalStorage (Optional):** Store token for manual header usage

### 3. Logout Flow

1. Call logout endpoint
2. Clear token cookie
3. Clear any local storage
4. Redirect to login

---

## Implementation Guide

### React/Next.js Example

#### 1. Create Auth Context

```typescript
// contexts/AuthContext.tsx
import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

interface User {
  id: string;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  role_id: string;
  role_name: string;
  active: boolean;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'https://playtzapi-production.up.railway.app/api/v1';

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Check authentication on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const response = await fetch(`${API_BASE}/auth/me`, {
        credentials: 'include', // Important: sends cookies
      });
      
      const data = await response.json();
      
      if (data.authenticated && data.user) {
        setUser(data.user);
      } else {
        setUser(null);
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const login = async (username: string, password: string) => {
    try {
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // Important: sends and receives cookies
        body: JSON.stringify({ username, password }),
      });

      const data = await response.json();

      if (response.ok && data.success) {
        // Token is automatically stored in cookie
        // Optionally store in localStorage for manual header usage
        if (data.token) {
          localStorage.setItem('token', data.token);
        }
        
        setUser(data.user);
        return;
      } else {
        // Handle error - don't log 401 to console (expected for wrong credentials)
        throw new Error(data.error || 'Login failed');
      }
    } catch (error) {
      // Only log unexpected errors
      if (error instanceof Error && !error.message.includes('401')) {
        console.error('Login error:', error);
      }
      throw error;
    }
  };

  const logout = async () => {
    try {
      await fetch(`${API_BASE}/auth/logout`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
      });
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      // Always clear local state
      setUser(null);
      localStorage.removeItem('token');
      // Redirect will be handled by component
    }
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
```

#### 2. Login Page Component

```typescript
// pages/login.tsx or app/login/page.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const { login, isAuthenticated } = useAuth();
  const router = useRouter();

  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated) {
      router.push('/dashboard');
    }
  }, [isAuthenticated, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsLoading(true);

    try {
      await login(username, password);
      router.push('/dashboard');
    } catch (err) {
      // Show error to user, don't log to console
      setError(err instanceof Error ? err.message : 'Login failed. Please check your credentials.');
      setPassword(''); // Clear password for security
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-purple-600 to-purple-800">
      <div className="bg-white p-8 rounded-lg shadow-xl w-full max-w-md">
        <h1 className="text-3xl font-bold text-center mb-6 text-purple-600">
          Playtz 102.9
        </h1>
        <p className="text-center text-gray-600 mb-6">Admin Portal</p>

        {error && (
          <div className="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit}>
          <div className="mb-4">
            <label htmlFor="username" className="block text-gray-700 font-medium mb-2">
              Username or Email
            </label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              autoComplete="username"
            />
          </div>

          <div className="mb-6">
            <label htmlFor="password" className="block text-gray-700 font-medium mb-2">
              Password
            </label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
              autoComplete="current-password"
            />
          </div>

          <button
            type="submit"
            disabled={isLoading}
            className="w-full bg-gradient-to-r from-purple-600 to-purple-800 text-white py-3 rounded-lg font-semibold hover:from-purple-700 hover:to-purple-900 disabled:opacity-50 disabled:cursor-not-allowed transition-all"
          >
            {isLoading ? 'Logging in...' : 'Login'}
          </button>
        </form>
      </div>
    </div>
  );
}
```

#### 3. Protected Route Component

```typescript
// components/ProtectedRoute.tsx
'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push('/login');
    }
  }, [loading, isAuthenticated, router]);

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-600"></div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return null;
  }

  return <>{children}</>;
}
```

#### 4. API Client Utility

```typescript
// utils/api.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'https://playtzapi-production.up.railway.app/api/v1';

// Get token from localStorage (optional - cookie is primary)
function getToken(): string | null {
  if (typeof window !== 'undefined') {
    return localStorage.getItem('token');
  }
  return null;
}

// Make authenticated API request
export async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...options.headers,
  };

  // Add Authorization header if token exists (cookie is primary, this is fallback)
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const response = await fetch(`${API_BASE}${endpoint}`, {
    ...options,
    headers,
    credentials: 'include', // Always include credentials for cookies
  });

  // Handle 401 - token expired or invalid
  if (response.status === 401) {
    // Clear token and redirect to login
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    throw new Error('Authentication required');
  }

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

// Convenience methods
export const api = {
  get: <T>(endpoint: string) => apiRequest<T>(endpoint, { method: 'GET' }),
  post: <T>(endpoint: string, data?: any) =>
    apiRequest<T>(endpoint, {
      method: 'POST',
      body: JSON.stringify(data),
    }),
  put: <T>(endpoint: string, data?: any) =>
    apiRequest<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    }),
  delete: <T>(endpoint: string) => apiRequest<T>(endpoint, { method: 'DELETE' }),
};
```

#### 5. Dashboard Component Example

```typescript
// pages/dashboard.tsx or app/dashboard/page.tsx
'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { api } from '@/utils/api';
import ProtectedRoute from '@/components/ProtectedRoute';

interface DashboardData {
  user: any;
  stats: {
    total_users: number;
    total_news: number;
    total_events: number;
    total_merchandise: number;
    total_orders: number;
    pending_orders: number;
  };
  recent_news?: any[];
  recent_events?: any[];
  recent_orders?: any[];
}

export default function DashboardPage() {
  const { user, logout } = useAuth();
  const [data, setData] = useState<DashboardData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadDashboard();
  }, []);

  const loadDashboard = async () => {
    try {
      const dashboardData = await api.get<DashboardData>('/admin/dashboard');
      setData(dashboardData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dashboard');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    await logout();
    window.location.href = '/login';
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-gray-50">
        <header className="bg-white shadow">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
            <h1 className="text-2xl font-bold text-gray-900">Admin Dashboard</h1>
            <div className="flex items-center gap-4">
              <span className="text-gray-600">Welcome, {user?.username}</span>
              <button
                onClick={handleLogout}
                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
              >
                Logout
              </button>
            </div>
          </div>
        </header>

        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {loading && <div>Loading...</div>}
          {error && <div className="text-red-600">{error}</div>}
          
          {data && (
            <div>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <div className="bg-white p-6 rounded-lg shadow">
                  <h3 className="text-gray-500 text-sm font-medium">Total Users</h3>
                  <p className="text-3xl font-bold text-purple-600">{data.stats.total_users}</p>
                </div>
                <div className="bg-white p-6 rounded-lg shadow">
                  <h3 className="text-gray-500 text-sm font-medium">News Articles</h3>
                  <p className="text-3xl font-bold text-purple-600">{data.stats.total_news}</p>
                </div>
                <div className="bg-white p-6 rounded-lg shadow">
                  <h3 className="text-gray-500 text-sm font-medium">Total Orders</h3>
                  <p className="text-3xl font-bold text-purple-600">{data.stats.total_orders}</p>
                </div>
              </div>

              {/* Add more dashboard content here */}
            </div>
          )}
        </main>
      </div>
    </ProtectedRoute>
  );
}
```

---

## API Endpoints Usage Examples

### Authentication Endpoints

#### Login
```typescript
const response = await fetch(`${API_BASE}/auth/login`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ username: 'admin', password: 'admin123' }),
});

const data = await response.json();
// data.token contains JWT token
// Token is also set in cookie automatically
```

#### Logout
```typescript
await fetch(`${API_BASE}/auth/logout`, {
  method: 'POST',
  credentials: 'include',
  headers: { 'Content-Type': 'application/json' },
});

// Clear local storage
localStorage.removeItem('token');
```

#### Check Current User
```typescript
const response = await fetch(`${API_BASE}/auth/me`, {
  credentials: 'include',
});

const data = await response.json();
// data.authenticated: boolean
// data.user: User | null
```

### Protected Endpoints Examples

#### Get News List
```typescript
const news = await api.get('/news');
```

#### Create News
```typescript
const newArticle = await api.post('/news', {
  title: 'News Title',
  content: 'News content...',
  author: 'Author Name',
  published: true,
});
```

#### Update News
```typescript
const updated = await api.put(`/news/${articleId}`, {
  title: 'Updated Title',
  content: 'Updated content...',
  published: true,
});
```

#### Delete News
```typescript
await api.delete(`/news/${articleId}`);
```

#### Get Users
```typescript
const users = await api.get('/users');
```

#### Create User
```typescript
const newUser = await api.post('/users', {
  email: 'user@example.com',
  username: 'username',
  first_name: 'First',
  last_name: 'Last',
  role_id: 'role-uuid',
  // password is optional - will generate default if not provided
});
```

#### Upload Image
```typescript
const formData = new FormData();
formData.append('file', file);

const response = await fetch(`${API_BASE}/upload`, {
  method: 'POST',
  credentials: 'include',
  body: formData,
  // Don't set Content-Type header - browser will set it with boundary
});

const data = await response.json();
// data.url contains the image URL
```

---

## Error Handling

### Handle 401 Errors

```typescript
try {
  const data = await api.get('/news');
} catch (error) {
  if (error.message === 'Authentication required') {
    // Token expired or invalid
    // Redirect to login (handled by apiRequest)
  } else {
    // Other errors
    console.error('API error:', error);
  }
}
```

### Handle Login Errors

```typescript
try {
  await login(username, password);
} catch (error) {
  // 401 errors are expected for wrong credentials
  // Don't log them to console
  // Show user-friendly message instead
  setError('Invalid username or password');
}
```

---

## Environment Variables

Create `.env.local` in your Next.js project:

```env
NEXT_PUBLIC_API_URL=https://playtzapi-production.up.railway.app/api/v1
```

---

## Complete Example: News Management Component

```typescript
'use client';

import { useState, useEffect } from 'react';
import { api } from '@/utils/api';

interface NewsArticle {
  id: string;
  title: string;
  content: string;
  author: string;
  published: boolean;
  created_at: string;
}

export default function NewsManagement() {
  const [articles, setArticles] = useState<NewsArticle[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    loadNews();
  }, []);

  const loadNews = async () => {
    try {
      const data = await api.get<NewsArticle[]>('/news');
      setArticles(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load news');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this article?')) return;

    try {
      await api.delete(`/news/${id}`);
      setArticles(articles.filter(a => a.id !== id));
    } catch (err) {
      alert('Failed to delete article');
    }
  };

  if (loading) return <div>Loading...</div>;
  if (error) return <div className="text-red-600">{error}</div>;

  return (
    <div>
      <h2 className="text-2xl font-bold mb-4">News Articles</h2>
      <div className="space-y-4">
        {articles.map(article => (
          <div key={article.id} className="bg-white p-4 rounded-lg shadow">
            <h3 className="text-xl font-semibold">{article.title}</h3>
            <p className="text-gray-600">{article.content}</p>
            <div className="mt-2 flex gap-2">
              <button
                onClick={() => handleDelete(article.id)}
                className="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700"
              >
                Delete
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## Important Notes

1. **Always use `credentials: 'include'`** in fetch requests to send/receive cookies
2. **Token is stored in cookie automatically** - no need to manually set it
3. **401 errors on login are expected** - handle gracefully, don't log to console
4. **Token expires in 24 hours** - implement token refresh if needed
5. **CORS is configured** for your frontend domain
6. **All endpoints require authentication** except `/auth/login`, `/auth/logout`, `/auth/me`, and `/health`

---

## Troubleshooting

### Token not being sent
- Ensure `credentials: 'include'` is set
- Check CORS configuration
- Verify cookie settings (SameSite=None, Secure)

### 401 errors on all requests
- Token may be expired (24h)
- User needs to login again
- Check if token cookie is being set

### CORS errors
- Verify your frontend domain is in `CORS_ORIGINS` environment variable
- Check that credentials are included in requests

---

## Quick Start Checklist

- [ ] Create AuthContext with login/logout functions
- [ ] Create Login page component
- [ ] Create ProtectedRoute component
- [ ] Create API utility with error handling
- [ ] Set up environment variables
- [ ] Test login flow
- [ ] Test protected endpoints
- [ ] Implement logout functionality
- [ ] Add error handling for 401 errors
- [ ] Test token expiration handling

---

*Last Updated: 2024*

