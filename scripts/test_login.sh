#!/bin/bash

# Test Login and Dashboard Access Script
# Usage: ./scripts/test_login.sh

BASE_URL="http://localhost:8080"
COOKIE_FILE="/tmp/test_login_cookies.txt"

echo "🧪 Testing Login and Dashboard Access"
echo "===================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo "1️⃣  Testing Server Health..."
HEALTH=$(curl -s "$BASE_URL/health")
if [[ $HEALTH == *"ok"* ]]; then
    echo -e "${GREEN}✅ Server is running${NC}"
else
    echo -e "${RED}❌ Server is not responding${NC}"
    exit 1
fi
echo ""

# Test 2: Login
echo "2️⃣  Testing Login API..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  -c "$COOKIE_FILE" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ Login successful${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null | head -15
else
    echo -e "${RED}❌ Login failed (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
    exit 1
fi
echo ""

# Test 3: Get Current User
echo "3️⃣  Testing Session (Get Current User)..."
USER_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/me" \
  -b "$COOKIE_FILE" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$USER_RESPONSE" | tail -n1)
BODY=$(echo "$USER_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ Session valid${NC}"
    USERNAME=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin)['username'])" 2>/dev/null)
    ROLE=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('role_name', 'N/A'))" 2>/dev/null)
    echo "   Username: $USERNAME"
    echo "   Role: $ROLE"
else
    echo -e "${RED}❌ Session invalid (HTTP $HTTP_CODE)${NC}"
fi
echo ""

# Test 4: Access Dashboard API
echo "4️⃣  Testing Admin Dashboard API..."
DASHBOARD_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/admin/dashboard" \
  -b "$COOKIE_FILE" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$DASHBOARD_RESPONSE" | tail -n1)
BODY=$(echo "$DASHBOARD_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ Dashboard access granted${NC}"
    echo "$BODY" | python3 -c "import sys, json; data=json.load(sys.stdin); print(f\"   User: {data['user']['username']} ({data['user']['role_name']})\"); print(f\"   Stats: {data['stats']}\")" 2>/dev/null
else
    echo -e "${RED}❌ Dashboard access denied (HTTP $HTTP_CODE)${NC}"
fi
echo ""

# Test 5: Logout
echo "5️⃣  Testing Logout..."
LOGOUT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/logout" \
  -b "$COOKIE_FILE" \
  -c "$COOKIE_FILE" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$LOGOUT_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ Logout successful${NC}"
else
    echo -e "${YELLOW}⚠️  Logout response: HTTP $HTTP_CODE${NC}"
fi
echo ""

# Summary
echo "===================================="
echo -e "${GREEN}✅ All tests completed!${NC}"
echo ""
echo "📝 Dashboard URLs:"
echo "   Login:    $BASE_URL/admin/login"
echo "   Dashboard: $BASE_URL/admin/dashboard"
echo ""
echo "🔐 Test Credentials:"
echo "   Username: admin"
echo "   Password: admin123"
echo ""

# Cleanup
rm -f "$COOKIE_FILE"

