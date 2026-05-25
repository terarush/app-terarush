# Authentication API

## Overview
Handles user authentication, account creation, and JWT token management. All protected endpoints require a valid JWT token in the Authorization header.

## Authentication Header
```
Authorization: Bearer <jwt_token>
```

## Endpoints

### Register User
**POST** `/api/v1/auth/register`

Creates a new user account.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "name": "John Doe"
}
```

**Response** (201 Created)
```json
{
  "id": "uuid-string",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2025-05-25T10:30:00Z"
}
```

**Error Responses**
- `400 Bad Request`: Invalid input, missing required fields
- `409 Conflict`: Email already registered

---

### Login User
**POST** `/api/v1/auth/login`

Authenticates user and returns JWT token.

**Request Body**
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200 OK)
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid-string",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

**Error Responses**
- `400 Bad Request`: Invalid credentials
- `401 Unauthorized`: User not found or password incorrect

---

### Get Current User
**GET** `/api/v1/auth/me`

Retrieves current authenticated user data. Requires valid JWT token.

**Response** (200 OK)
```json
{
  "id": "uuid-string",
  "email": "user@example.com",
  "name": "John Doe",
  "created_at": "2025-05-25T10:30:00Z"
}
```

**Error Responses**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

---

## Security Notes
- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- Refresh tokens available on request
- All authentication endpoints use HTTPS in production
