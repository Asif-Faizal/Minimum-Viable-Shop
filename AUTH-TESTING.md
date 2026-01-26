# Authentication Testing Guide

This document contains curl commands for testing the authentication REST API endpoints.

## Prerequisites

All authentication requests require the following device headers:
- `X-Device-ID`: Unique device identifier
- `X-Device-Type`: Type of device (mobile, desktop, other)
- `X-Device-Model`: Device model (e.g., iPhone 15)
- `X-Device-OS`: Operating system (e.g., iOS, Android)
- `X-Device-OS-Version`: OS version (optional)

## Health Check

```bash
curl -X GET http://localhost:8081/health
```

## Check Email Availability

```bash
curl -X GET "http://localhost:8081/accounts/check-email?email=john@example.com"
```

## Login

Creates a new session or updates existing session for the same device.

```bash
curl -X POST http://localhost:8081/accounts/login \
  -H "Content-Type: application/json" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS" \
  -H "X-Device-OS-Version: 17.0" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

## Refresh Token

Generates new access and refresh tokens. Validates device ID matches session.

```bash
curl -X POST http://localhost:8081/accounts/refresh \
  -H "Content-Type: application/json" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS" \
  -d '{
    "refresh_token": "YOUR_REFRESH_TOKEN"
  }'
```

## Logout

Revokes the session. Validates device ID matches session.

```bash
curl -X POST http://localhost:8081/accounts/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS"
```

## Testing Device Mismatch

Attempt logout with different device ID (should fail):

```bash
curl -X POST http://localhost:8081/accounts/logout \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "X-Device-ID: different-device" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS"
```

## Testing Multiple Devices

Login from different device (creates separate session):

```bash
curl -X POST http://localhost:8081/accounts/login \
  -H "Content-Type: application/json" \
  -H "X-Device-ID: device-456" \
  -H "X-Device-Type: desktop" \
  -H "X-Device-Model: MacBook Pro" \
  -H "X-Device-OS: macOS" \
  -H "X-Device-OS-Version: 14.0" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

## Complete Flow Example

```bash
# 1. Login
RESPONSE=$(curl -s -X POST http://localhost:8081/accounts/login \
  -H "Content-Type: application/json" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS" \
  -d '{"email": "john@example.com", "password": "password123"}')

echo "$RESPONSE"

# Extract tokens (requires jq)
ACCESS_TOKEN=$(echo "$RESPONSE" | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo "$RESPONSE" | jq -r '.data.refresh_token')

# 2. Refresh
curl -X POST http://localhost:8081/accounts/refresh \
  -H "Content-Type: application/json" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS" \
  -d "{\"refresh_token\": \"$REFRESH_TOKEN\"}"

# 3. Logout
curl -X POST http://localhost:8081/accounts/logout \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "X-Device-ID: device-123" \
  -H "X-Device-Type: mobile" \
  -H "X-Device-Model: iPhone 15" \
  -H "X-Device-OS: iOS"
```
