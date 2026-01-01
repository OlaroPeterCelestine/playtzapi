#!/bin/bash

# Test Login Endpoint Script
# Tests the Railway production login endpoint

BASE_URL="https://playtzapi-production.up.railway.app"
COOKIE_FILE="/tmp/railway_test_cookies.txt"

echo "🧪 Testing Login Endpoint"
echo "========================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Health Check
echo "1️⃣  Testing Server Health..."
HEALTH=$(curl -s "$BASE_URL/health")
if [[ $HEALTH == *"ok"* ]] || [[ $HEALTH == *"status"* ]]; then
    echo -e "${GREEN}✅ Server is running${NC}"
    echo "$HEALTH" | python3 -m json.tool 2>/dev/null || echo "$HEALTH"
else
    echo -e "${RED}❌ Server health check failed${NC}"
    echo "$HEALTH"
    exit 1
fi
echo ""

# Test 2: CORS Preflight
echo "2️⃣  Testing CORS Preflight..."
CORS_RESPONSE=$(curl -s -X OPTIONS "$BASE_URL/api/v1/auth/login" \
  -H "Origin: http://localhost:3001" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$CORS_RESPONSE" | tail -n1)
if [ "$HTTP_CODE" == "204" ] || [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ CORS preflight successful (HTTP $HTTP_CODE)${NC}"
else
    echo -e "${YELLOW}⚠️  CORS preflight returned HTTP $HTTP_CODE${NC}"
fi
echo ""

# Test 3: Valid Login
echo "3️⃣  Testing Login with Valid Credentials..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3001" \
  -d '{"username":"admin","password":"admin123"}' \
  -c "$COOKIE_FILE" \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$LOGIN_RESPONSE" | tail -n1)
BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "200" ]; then
    echo -e "${GREEN}✅ Login successful (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null | head -20 || echo "$BODY"
    
    # Check for session cookie
    if [ -f "$COOKIE_FILE" ] && grep -q "session_id" "$COOKIE_FILE"; then
        echo -e "${GREEN}✅ Session cookie set${NC}"
    else
        echo -e "${YELLOW}⚠️  Session cookie not found${NC}"
    fi
else
    echo -e "${RED}❌ Login failed (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
fi
echo ""

# Test 4: Invalid Credentials
echo "4️⃣  Testing Login with Invalid Credentials..."
INVALID_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3001" \
  -d '{"username":"admin","password":"wrongpassword"}' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$INVALID_RESPONSE" | tail -n1)
BODY=$(echo "$INVALID_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "401" ]; then
    echo -e "${GREEN}✅ Correctly rejected invalid credentials (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
else
    echo -e "${YELLOW}⚠️  Unexpected response (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
fi
echo ""

# Test 5: Missing Fields
echo "5️⃣  Testing Login with Missing Fields..."
MISSING_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3001" \
  -d '{"username":"admin"}' \
  -w "\n%{http_code}")

HTTP_CODE=$(echo "$MISSING_RESPONSE" | tail -n1)
BODY=$(echo "$MISSING_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" == "400" ]; then
    echo -e "${GREEN}✅ Correctly rejected missing fields (HTTP $HTTP_CODE)${NC}"
    echo "$BODY" | python3 -m json.tool 2>/dev/null || echo "$BODY"
else
    echo -e "${YELLOW}⚠️  Unexpected response (HTTP $HTTP_CODE)${NC}"
    echo "$BODY"
fi
echo ""

# Test 6: Session Validation (if login was successful)
if [ -f "$COOKIE_FILE" ] && grep -q "session_id" "$COOKIE_FILE"; then
    echo "6️⃣  Testing Session Validation..."
    SESSION_RESPONSE=$(curl -s -X GET "$BASE_URL/api/v1/auth/me" \
      -H "Origin: http://localhost:3001" \
      -b "$COOKIE_FILE" \
      -w "\n%{http_code}")

    HTTP_CODE=$(echo "$SESSION_RESPONSE" | tail -n1)
    BODY=$(echo "$SESSION_RESPONSE" | sed '$d')

    if [ "$HTTP_CODE" == "200" ]; then
        echo -e "${GREEN}✅ Session valid (HTTP $HTTP_CODE)${NC}"
        USERNAME=$(echo "$BODY" | python3 -c "import sys, json; print(json.load(sys.stdin).get('username', 'N/A'))" 2>/dev/null)
        if [ "$USERNAME" != "N/A" ]; then
            echo "   Username: $USERNAME"
        fi
    else
        echo -e "${YELLOW}⚠️  Session validation failed (HTTP $HTTP_CODE)${NC}"
        echo "$BODY"
    fi
    echo ""
fi

# Test 7: CORS Headers
echo "7️⃣  Checking CORS Headers..."
CORS_HEADERS=$(curl -s -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -H "Origin: http://localhost:3001" \
  -d '{"username":"admin","password":"admin123"}' \
  -v 2>&1 | grep -i "access-control")

if [ -n "$CORS_HEADERS" ]; then
    echo -e "${GREEN}✅ CORS headers present:${NC}"
    echo "$CORS_HEADERS"
else
    echo -e "${YELLOW}⚠️  No CORS headers found${NC}"
fi
echo ""

# Summary
echo "========================="
echo -e "${GREEN}✅ All tests completed!${NC}"
echo ""
echo "📝 Endpoint: $BASE_URL/api/v1/auth/login"
echo "📝 Method: POST"
echo "📝 Content-Type: application/json"
echo ""
echo "Example Request:"
echo 'curl -X POST '"$BASE_URL/api/v1/auth/login"' \'
echo '  -H "Content-Type: application/json" \'
echo '  -H "Origin: http://localhost:3001" \'
echo '  -d '"'"'{"username":"admin","password":"admin123"}'"'"' \'
echo '  -c cookies.txt'
echo ""

# Cleanup
rm -f "$COOKIE_FILE"

