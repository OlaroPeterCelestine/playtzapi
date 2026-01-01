#!/bin/bash

# Comprehensive Dashboard CRUD Test Script
# Tests Create, Read, Update, Delete for all resources

API_BASE="https://playtzapi-production.up.railway.app/api/v1"
ADMIN_USER="admin"
ADMIN_PASS="admin123"

echo "🧪 Dashboard CRUD Test Suite"
echo "=============================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Login and get session
echo "🔐 Step 1: Login..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -H "Origin: https://playtzadmin.vercel.app" \
  -d "{\"username\":\"$ADMIN_USER\",\"password\":\"$ADMIN_PASS\"}" \
  -c /tmp/dashboard_cookies.txt)

if echo "$LOGIN_RESPONSE" | grep -q "success.*true"; then
  echo -e "${GREEN}✅ Login successful${NC}"
  SESSION_ID=$(echo "$LOGIN_RESPONSE" | grep -o '"session_id":"[^"]*' | cut -d'"' -f4)
else
  echo -e "${RED}❌ Login failed${NC}"
  echo "$LOGIN_RESPONSE"
  exit 1
fi

echo ""

# Function to test CRUD operations
test_crud() {
  local resource=$1
  local create_data=$2
  local update_data=$3
  
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "📦 Testing: $resource"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  
  # CREATE
  echo -n "  CREATE: "
  CREATE_RESPONSE=$(curl -s -X POST "$API_BASE/$resource" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt \
    -d "$create_data")
  
  if echo "$CREATE_RESPONSE" | grep -qE '"id"|"success"|"message"'; then
    echo -e "${GREEN}✅${NC}"
    RESOURCE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
    if [ -z "$RESOURCE_ID" ]; then
      RESOURCE_ID=$(echo "$CREATE_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))" 2>/dev/null)
    fi
    echo "     ID: $RESOURCE_ID"
  else
    echo -e "${RED}❌${NC}"
    echo "     Response: $CREATE_RESPONSE"
    return 1
  fi
  
  if [ -z "$RESOURCE_ID" ]; then
    echo -e "${YELLOW}⚠️  No ID returned, skipping remaining tests${NC}"
    return 1
  fi
  
  # READ (GET by ID)
  echo -n "  READ:   "
  READ_RESPONSE=$(curl -s -X GET "$API_BASE/$resource/$RESOURCE_ID" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt)
  
  if echo "$READ_RESPONSE" | grep -qE '"id"|"title"|"name"'; then
    echo -e "${GREEN}✅${NC}"
  else
    echo -e "${RED}❌${NC}"
    echo "     Response: $READ_RESPONSE"
  fi
  
  # READ ALL
  echo -n "  LIST:   "
  LIST_RESPONSE=$(curl -s -X GET "$API_BASE/$resource" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt)
  
  if echo "$LIST_RESPONSE" | grep -qE '\[|"id"'; then
    echo -e "${GREEN}✅${NC}"
  else
    echo -e "${RED}❌${NC}"
  fi
  
  # UPDATE
  if [ -n "$update_data" ]; then
    echo -n "  UPDATE: "
    UPDATE_RESPONSE=$(curl -s -X PUT "$API_BASE/$resource/$RESOURCE_ID" \
      -H "Content-Type: application/json" \
      -H "Origin: https://playtzadmin.vercel.app" \
      -b /tmp/dashboard_cookies.txt \
      -d "$update_data")
    
    if echo "$UPDATE_RESPONSE" | grep -qE '"id"|"success"|"message"'; then
      echo -e "${GREEN}✅${NC}"
    else
      echo -e "${RED}❌${NC}"
      echo "     Response: $UPDATE_RESPONSE"
    fi
  else
    echo -e "${YELLOW}  UPDATE: ⚠️  Skipped (no update data)${NC}"
  fi
  
  # DELETE
  echo -n "  DELETE: "
  DELETE_RESPONSE=$(curl -s -X DELETE "$API_BASE/$resource/$RESOURCE_ID" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt)
  
  if echo "$DELETE_RESPONSE" | grep -qE '"message"|"success"|200'; then
    echo -e "${GREEN}✅${NC}"
  else
    echo -e "${RED}❌${NC}"
    echo "     Response: $DELETE_RESPONSE"
  fi
  
  echo ""
}

# Test News
test_crud "news" \
  '{"title":"Test News Article","content":"This is a test article","author":"Test Author","published":true}' \
  '{"title":"Updated News Article","content":"Updated content","author":"Test Author","published":true}'

# Test Events
test_crud "events" \
  '{"title":"Test Event","description":"Test event description","date":"2026-12-31","time":"18:00:00","location":"Test Location","active":true}' \
  '{"title":"Updated Event","description":"Updated description","date":"2026-12-31","time":"19:00:00","location":"Updated Location","active":true}'

# Test Merchandise
test_crud "merch" \
  '{"name":"Test Product","description":"Test product description","price":29.99,"stock":100,"active":true}' \
  '{"name":"Updated Product","description":"Updated description","price":39.99,"stock":50,"active":true}'

# Test Careers
test_crud "careers" \
  '{"title":"Test Job","description":"Test job description","department":"Engineering","location":"Remote","type":"full-time","active":true}' \
  '{"title":"Updated Job","description":"Updated description","department":"Engineering","location":"Hybrid","type":"full-time","active":true}'

# Test Rooms
test_crud "rooms" \
  '{"name":"Test Room","genre":"Test Genre","description":"Test room description","gradient":"linear-gradient(45deg, #ff0000, #00ff00)","text_color":"#ffffff","active":true}' \
  '{"name":"Updated Room","genre":"Updated Genre","description":"Updated description","gradient":"linear-gradient(45deg, #0000ff, #ffff00)","text_color":"#000000","active":true}'

# Test Mixes (requires room_id, so we'll create a room first)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Testing: Mixes (requires Room)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Create a room for the mix
ROOM_RESPONSE=$(curl -s -X POST "$API_BASE/rooms" \
  -H "Content-Type: application/json" \
  -H "Origin: https://playtzadmin.vercel.app" \
  -b /tmp/dashboard_cookies.txt \
  -d '{"name":"Test Room for Mix","genre":"Test","description":"Test","gradient":"#ff0000","text_color":"#ffffff","active":true}')

ROOM_ID=$(echo "$ROOM_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
if [ -z "$ROOM_ID" ]; then
  ROOM_ID=$(echo "$ROOM_RESPONSE" | python3 -c "import sys, json; print(json.load(sys.stdin).get('id', ''))" 2>/dev/null)
fi

if [ -n "$ROOM_ID" ]; then
  test_crud "mixes" \
    "{\"room_id\":\"$ROOM_ID\",\"title\":\"Test Mix\",\"artist\":\"Test Artist\",\"description\":\"Test mix description\",\"duration\":\"60:00\",\"tracks\":10,\"active\":true}" \
    "{\"room_id\":\"$ROOM_ID\",\"title\":\"Updated Mix\",\"artist\":\"Updated Artist\",\"description\":\"Updated description\",\"duration\":\"75:00\",\"tracks\":15,\"active\":true}"
  
  # Clean up room
  curl -s -X DELETE "$API_BASE/rooms/$ROOM_ID" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt > /dev/null
else
  echo -e "${RED}❌ Failed to create room for mix test${NC}"
fi

echo ""

# Test Users (requires role_id, so we'll get/create a role first)
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Testing: Users (requires Role)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Get first role or create one
ROLES_RESPONSE=$(curl -s -X GET "$API_BASE/roles" \
  -H "Content-Type: application/json" \
  -H "Origin: https://playtzadmin.vercel.app" \
  -b /tmp/dashboard_cookies.txt)

ROLE_ID=$(echo "$ROLES_RESPONSE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
if [ -z "$ROLE_ID" ]; then
  # Create a role
  ROLE_CREATE=$(curl -s -X POST "$API_BASE/roles" \
    -H "Content-Type: application/json" \
    -H "Origin: https://playtzadmin.vercel.app" \
    -b /tmp/dashboard_cookies.txt \
    -d '{"name":"Test Role","description":"Test role","permissions":["test.read"],"active":true}')
  ROLE_ID=$(echo "$ROLE_CREATE" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
fi

if [ -n "$ROLE_ID" ]; then
  TIMESTAMP=$(date +%s)
  test_crud "users" \
    "{\"email\":\"testuser$TIMESTAMP@example.com\",\"username\":\"testuser$TIMESTAMP\",\"first_name\":\"Test\",\"last_name\":\"User\",\"role_id\":\"$ROLE_ID\"}" \
    "{\"email\":\"testuser$TIMESTAMP@example.com\",\"username\":\"testuser$TIMESTAMP\",\"first_name\":\"Updated\",\"last_name\":\"Name\",\"role_id\":\"$ROLE_ID\"}"
else
  echo -e "${RED}❌ Failed to get/create role for user test${NC}"
fi

echo ""

# Test Roles
test_crud "roles" \
  '{"name":"Test Role CRUD","description":"Test role for CRUD","permissions":["test.read","test.write"],"active":true}' \
  '{"name":"Test Role CRUD","description":"Updated role description","permissions":["test.read","test.write","test.delete"],"active":true}'

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "✅ CRUD Test Complete!"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Cleanup
rm -f /tmp/dashboard_cookies.txt

