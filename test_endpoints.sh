#!/bin/bash

BASE_URL="http://localhost:8080"
API_URL="${BASE_URL}/api/v1"

echo "🧪 Testing Playtz API Endpoints"
echo "================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables to store created IDs
ROLE_ID=""
USER_ID=""
ROOM_ID=""
MIX_ID=""
NEWS_ID=""
EVENT_ID=""
MERCH_ID=""
CAREER_ID=""

test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    local capture_id=$5  # Variable name to capture ID from response
    
    echo -n "Testing $description... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method "$url" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "${GREEN}✓${NC} (HTTP $http_code)"
        
        # Capture ID if requested
        if [ -n "$capture_id" ]; then
            captured_id=$(echo "$body" | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
            if [ -n "$captured_id" ]; then
                eval "$capture_id='$captured_id'"
                echo -e "  ${BLUE}→ Created ID: $captured_id${NC}"
            fi
        fi
    elif [ "$http_code" -ge 400 ] && [ "$http_code" -lt 500 ]; then
        echo -e "${YELLOW}⚠${NC} (HTTP $http_code) - Client Error"
        echo "  Response: $(echo "$body" | head -c 150)"
        echo ""
    else
        echo -e "${RED}✗${NC} (HTTP $http_code)"
        echo "  Response: $(echo "$body" | head -c 150)"
        echo ""
    fi
}

# Health Check
echo "📊 Health & Debug Endpoints"
test_endpoint "GET" "${BASE_URL}/health" "" "Health Check"
test_endpoint "GET" "${BASE_URL}/debug/cloudinary" "" "Cloudinary Debug"
echo ""

# Roles (needed for users)
echo "👥 Roles Endpoints"
test_endpoint "GET" "${API_URL}/roles" "" "GET all roles"
ROLE_DATA="{\"name\":\"Test Role $(date +%s)\",\"description\":\"Test role description\",\"permissions\":[\"read\",\"write\"]}"
test_endpoint "POST" "${API_URL}/roles" "$ROLE_DATA" "POST create role" "ROLE_ID"
if [ -n "$ROLE_ID" ]; then
    test_endpoint "GET" "${API_URL}/roles/$ROLE_ID" "" "GET role by ID"
fi
echo ""

# Users (requires role)
echo "👤 Users Endpoints"
test_endpoint "GET" "${API_URL}/users" "" "GET all users"
if [ -n "$ROLE_ID" ]; then
    USER_DATA="{\"email\":\"test$(date +%s)@example.com\",\"username\":\"testuser$(date +%s)\",\"password\":\"password123\",\"first_name\":\"Test\",\"last_name\":\"User\",\"role_id\":\"$ROLE_ID\"}"
    test_endpoint "POST" "${API_URL}/users" "$USER_DATA" "POST create user" "USER_ID"
    if [ -n "$USER_ID" ]; then
        test_endpoint "GET" "${API_URL}/users/$USER_ID" "" "GET user by ID"
    fi
else
    echo -e "${YELLOW}⚠ Skipping user creation - no role ID available${NC}"
fi
echo ""

# News
echo "📰 News Endpoints"
test_endpoint "GET" "${API_URL}/news" "" "GET all news"
NEWS_DATA='{"title":"Test News Article","content":"This is test news content for testing the API endpoints.","author":"Test Author"}'
test_endpoint "POST" "${API_URL}/news" "$NEWS_DATA" "POST create news" "NEWS_ID"
if [ -n "$NEWS_ID" ]; then
    test_endpoint "GET" "${API_URL}/news/$NEWS_ID" "" "GET news by ID"
fi
echo ""

# Events
echo "🎉 Events Endpoints"
test_endpoint "GET" "${API_URL}/events" "" "GET all events"
EVENT_DATA='{"title":"Test Event","description":"Test event description for API testing","date":"2024-12-31","time":"18:00:00","location":"Test Location"}'
test_endpoint "POST" "${API_URL}/events" "$EVENT_DATA" "POST create event" "EVENT_ID"
if [ -n "$EVENT_ID" ]; then
    test_endpoint "GET" "${API_URL}/events/$EVENT_ID" "" "GET event by ID"
fi
echo ""

# Merchandise
echo "🛍️  Merchandise Endpoints"
test_endpoint "GET" "${API_URL}/merch" "" "GET all merchandise"
MERCH_DATA='{"name":"Test Product","description":"Test product description for API testing","price":29.99,"stock":100}'
test_endpoint "POST" "${API_URL}/merch" "$MERCH_DATA" "POST create merchandise" "MERCH_ID"
if [ -n "$MERCH_ID" ]; then
    test_endpoint "GET" "${API_URL}/merch/$MERCH_ID" "" "GET merchandise by ID"
fi
echo ""

# Careers
echo "💼 Careers Endpoints"
test_endpoint "GET" "${API_URL}/careers" "" "GET all careers"
CAREER_DATA='{"title":"Test Position","description":"Test job description for API testing","department":"Engineering","location":"Remote","type":"full-time"}'
test_endpoint "POST" "${API_URL}/careers" "$CAREER_DATA" "POST create career" "CAREER_ID"
if [ -n "$CAREER_ID" ]; then
    test_endpoint "GET" "${API_URL}/careers/$CAREER_ID" "" "GET career by ID"
fi
echo ""

# Rooms
echo "🏠 Rooms Endpoints"
test_endpoint "GET" "${API_URL}/rooms" "" "GET all rooms"
ROOM_DATA='{"name":"Test Room","genre":"Electronic","description":"Test room description for API testing","gradient":"linear-gradient(45deg, #ff0000, #00ff00)","text_color":"#ffffff"}'
test_endpoint "POST" "${API_URL}/rooms" "$ROOM_DATA" "POST create room" "ROOM_ID"
if [ -n "$ROOM_ID" ]; then
    test_endpoint "GET" "${API_URL}/rooms/$ROOM_ID" "" "GET room by ID"
fi
echo ""

# Mixes (requires room)
echo "🎵 Mixes Endpoints"
test_endpoint "GET" "${API_URL}/mixes" "" "GET all mixes"
if [ -n "$ROOM_ID" ]; then
    MIX_DATA="{\"room_id\":\"$ROOM_ID\",\"title\":\"Test Mix\",\"artist\":\"Test Artist\",\"description\":\"Test mix description\",\"duration\":\"60:00\",\"tracks\":10}"
    test_endpoint "POST" "${API_URL}/mixes" "$MIX_DATA" "POST create mix" "MIX_ID"
    if [ -n "$MIX_ID" ]; then
        test_endpoint "GET" "${API_URL}/mixes/$MIX_ID" "" "GET mix by ID"
    fi
else
    echo -e "${YELLOW}⚠ Skipping mix creation - no room ID available${NC}"
fi
echo ""

# Cart (requires merchandise)
echo "🛒 Cart Endpoints"
test_endpoint "GET" "${API_URL}/cart" "" "GET cart"
if [ -n "$MERCH_ID" ]; then
    CART_DATA="{\"merchandise_id\":\"$MERCH_ID\",\"quantity\":2}"
    test_endpoint "POST" "${API_URL}/cart/add" "$CART_DATA" "POST add to cart"
else
    echo -e "${YELLOW}⚠ Skipping cart add - no merchandise ID available${NC}"
fi
echo ""

# Orders
echo "📦 Orders Endpoints"
test_endpoint "GET" "${API_URL}/orders" "" "GET all orders"
if [ -n "$USER_ID" ]; then
    ORDER_DATA="{\"user_id\":\"$USER_ID\",\"total\":59.98,\"shipping_address\":\"123 Test St, Test City\"}"
    test_endpoint "POST" "${API_URL}/checkout" "$ORDER_DATA" "POST checkout"
else
    echo -e "${YELLOW}⚠ Skipping checkout - no user ID available${NC}"
fi
echo ""

# Upload (will likely fail without actual image, but testing endpoint)
echo "📤 Upload Endpoints"
echo "Note: Upload endpoints require actual image files, testing endpoint structure only"
test_endpoint "POST" "${API_URL}/upload" '{}' "POST upload (will fail without file)"
echo ""

echo "✅ Endpoint testing complete!"
echo ""
echo "Summary:"
echo "  - Health & Debug: ✓"
echo "  - Roles: Created ID: ${ROLE_ID:-none}"
echo "  - Users: Created ID: ${USER_ID:-none}"
echo "  - News: Created ID: ${NEWS_ID:-none}"
echo "  - Events: Created ID: ${EVENT_ID:-none}"
echo "  - Merchandise: Created ID: ${MERCH_ID:-none}"
echo "  - Careers: Created ID: ${CAREER_ID:-none}"
echo "  - Rooms: Created ID: ${ROOM_ID:-none}"
echo "  - Mixes: Created ID: ${MIX_ID:-none}"
