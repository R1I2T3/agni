# Authentication

## hmacLogin

**POST /api/auth/login**

HMAC-based login for client applications; returns a JWT in an HTTP-only cookie

### Parameters

None

### Responses

- **200**: Authentication successful; JWT set in `Agni-auth-token` cookie
- **400**: Invalid request body or missing required fields
- **401**: Invalid credentials or HMAC verification failed
- **500**: Token generation failed


