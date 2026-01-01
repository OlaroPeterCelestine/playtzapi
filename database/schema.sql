-- Playtz 102.9 Database Schema

-- Roles table
CREATE TABLE IF NOT EXISTS roles (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    permissions TEXT[], -- Array of permission strings
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role_id VARCHAR(50) REFERENCES roles(id),
    active BOOLEAN DEFAULT true,
    password_change_required BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Rooms (Genres) table
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    genre VARCHAR(100),
    description TEXT,
    gradient VARCHAR(100),
    text_color VARCHAR(50),
    image TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Mixes table
CREATE TABLE IF NOT EXISTS mixes (
    id VARCHAR(50) PRIMARY KEY,
    room_id VARCHAR(50) REFERENCES rooms(id),
    title VARCHAR(255) NOT NULL,
    artist VARCHAR(255),
    description TEXT,
    duration VARCHAR(50),
    tracks INTEGER DEFAULT 0,
    color VARCHAR(100),
    text_color VARCHAR(50),
    border_color VARCHAR(50),
    image TEXT,
    audio_url TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tracks table (for mix tracks)
CREATE TABLE IF NOT EXISTS tracks (
    id VARCHAR(50) PRIMARY KEY,
    mix_id VARCHAR(50) REFERENCES mixes(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    title VARCHAR(255),
    artist VARCHAR(255),
    duration VARCHAR(50),
    link TEXT NOT NULL,
    type VARCHAR(20) DEFAULT 'audio', -- 'audio' or 'video'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- News table
CREATE TABLE IF NOT EXISTS news (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    author VARCHAR(255),
    image TEXT,
    published BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Events table
CREATE TABLE IF NOT EXISTS events (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    date DATE,
    time TIME,
    location VARCHAR(255),
    image TEXT,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Merchandise table
CREATE TABLE IF NOT EXISTS merchandise (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    image TEXT,
    stock INTEGER DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Careers table
CREATE TABLE IF NOT EXISTS careers (
    id VARCHAR(50) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    department VARCHAR(100),
    location VARCHAR(255),
    type VARCHAR(50), -- 'full-time', 'part-time', 'contract'
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cart table (for shopping cart)
CREATE TABLE IF NOT EXISTS cart (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50), -- Can be null for guest carts
    session_id VARCHAR(255), -- For guest carts
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cart items table
CREATE TABLE IF NOT EXISTS cart_items (
    id VARCHAR(50) PRIMARY KEY,
    cart_id VARCHAR(50) REFERENCES cart(id) ON DELETE CASCADE,
    merchandise_id VARCHAR(50) REFERENCES merchandise(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(50),
    total DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'processing', 'shipped', 'delivered', 'cancelled'
    shipping_address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Order items table
CREATE TABLE IF NOT EXISTS order_items (
    id VARCHAR(50) PRIMARY KEY,
    order_id VARCHAR(50) REFERENCES orders(id) ON DELETE CASCADE,
    merchandise_id VARCHAR(50) REFERENCES merchandise(id),
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================
-- INDEXES FOR BETTER PERFORMANCE
-- ============================================

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role_id);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);
CREATE INDEX IF NOT EXISTS idx_users_role_active ON users(role_id, active);

-- Roles table indexes
CREATE INDEX IF NOT EXISTS idx_roles_active ON roles(active);
CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);

-- Rooms table indexes
CREATE INDEX IF NOT EXISTS idx_rooms_active ON rooms(active);
CREATE INDEX IF NOT EXISTS idx_rooms_genre ON rooms(genre);
CREATE INDEX IF NOT EXISTS idx_rooms_name ON rooms(name);
CREATE INDEX IF NOT EXISTS idx_rooms_active_genre ON rooms(active, genre);

-- Mixes table indexes
CREATE INDEX IF NOT EXISTS idx_mixes_room ON mixes(room_id);
CREATE INDEX IF NOT EXISTS idx_mixes_active ON mixes(active);
CREATE INDEX IF NOT EXISTS idx_mixes_room_active ON mixes(room_id, active);
CREATE INDEX IF NOT EXISTS idx_mixes_title ON mixes(title);
CREATE INDEX IF NOT EXISTS idx_mixes_artist ON mixes(artist);
CREATE INDEX IF NOT EXISTS idx_mixes_created_at ON mixes(created_at DESC);

-- Tracks table indexes
CREATE INDEX IF NOT EXISTS idx_tracks_mix ON tracks(mix_id);
CREATE INDEX IF NOT EXISTS idx_tracks_mix_number ON tracks(mix_id, number);
CREATE INDEX IF NOT EXISTS idx_tracks_type ON tracks(type);
CREATE INDEX IF NOT EXISTS idx_tracks_artist ON tracks(artist);

-- News table indexes
CREATE INDEX IF NOT EXISTS idx_news_published ON news(published);
CREATE INDEX IF NOT EXISTS idx_news_published_created ON news(published, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_news_title ON news(title);
CREATE INDEX IF NOT EXISTS idx_news_author ON news(author);
CREATE INDEX IF NOT EXISTS idx_news_created_at ON news(created_at DESC);

-- Events table indexes
CREATE INDEX IF NOT EXISTS idx_events_active ON events(active);
CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
CREATE INDEX IF NOT EXISTS idx_events_date_active ON events(date, active);
CREATE INDEX IF NOT EXISTS idx_events_location ON events(location);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at DESC);

-- Merchandise table indexes
CREATE INDEX IF NOT EXISTS idx_merchandise_active ON merchandise(active);
CREATE INDEX IF NOT EXISTS idx_merchandise_name ON merchandise(name);
CREATE INDEX IF NOT EXISTS idx_merchandise_price ON merchandise(price);
CREATE INDEX IF NOT EXISTS idx_merchandise_stock ON merchandise(stock);
CREATE INDEX IF NOT EXISTS idx_merchandise_active_stock ON merchandise(active, stock);
CREATE INDEX IF NOT EXISTS idx_merchandise_created_at ON merchandise(created_at DESC);

-- Careers table indexes
CREATE INDEX IF NOT EXISTS idx_careers_active ON careers(active);
CREATE INDEX IF NOT EXISTS idx_careers_department ON careers(department);
CREATE INDEX IF NOT EXISTS idx_careers_type ON careers(type);
CREATE INDEX IF NOT EXISTS idx_careers_location ON careers(location);
CREATE INDEX IF NOT EXISTS idx_careers_active_department ON careers(active, department);
CREATE INDEX IF NOT EXISTS idx_careers_created_at ON careers(created_at DESC);

-- Cart table indexes
CREATE INDEX IF NOT EXISTS idx_cart_user ON cart(user_id);
CREATE INDEX IF NOT EXISTS idx_cart_session ON cart(session_id);
CREATE INDEX IF NOT EXISTS idx_cart_user_session ON cart(user_id, session_id);
CREATE INDEX IF NOT EXISTS idx_cart_created_at ON cart(created_at DESC);

-- Cart items table indexes
CREATE INDEX IF NOT EXISTS idx_cart_items_cart ON cart_items(cart_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_merchandise ON cart_items(merchandise_id);
CREATE INDEX IF NOT EXISTS idx_cart_items_cart_merchandise ON cart_items(cart_id, merchandise_id);

-- Orders table indexes
CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_user_status ON orders(user_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_orders_status_created ON orders(status, created_at DESC);

-- Order items table indexes
CREATE INDEX IF NOT EXISTS idx_order_items_order ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_merchandise ON order_items(merchandise_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_merchandise ON order_items(order_id, merchandise_id);

-- ============================================
-- MIGRATIONS FOR EXISTING TABLES
-- ============================================

-- Add password_change_required column to users table (if not exists)
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'password_change_required'
    ) THEN
        ALTER TABLE users ADD COLUMN password_change_required BOOLEAN DEFAULT false;
    END IF;
END $$;

-- Full-text search indexes (using GIN for better text search performance)
-- Note: These require the pg_trgm extension for trigram matching
-- CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Text search indexes for common search fields
-- CREATE INDEX IF NOT EXISTS idx_news_title_trgm ON news USING gin(title gin_trgm_ops);
-- CREATE INDEX IF NOT EXISTS idx_news_content_trgm ON news USING gin(content gin_trgm_ops);
-- CREATE INDEX IF NOT EXISTS idx_events_title_trgm ON events USING gin(title gin_trgm_ops);
-- CREATE INDEX IF NOT EXISTS idx_merchandise_name_trgm ON merchandise USING gin(name gin_trgm_ops);

