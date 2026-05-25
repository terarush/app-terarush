# Blogs API

## Overview
Manages blog creation, retrieval, updating, and deletion. Supports rich markdown content, image uploads, and slug-based URL generation.

## Authentication
Protected endpoints (marked with 🔒) require valid JWT token in Authorization header.

## Endpoints

### Get All Blogs
**GET** `/api/v1/blogs`

Retrieve paginated list of published blogs.

**Query Parameters**
- `page` (int, default: 1): Page number
- `limit` (int, default: 10): Items per page
- `search` (string, optional): Search in title or content

**Response** (200 OK)
```json
{
  "data": [
    {
      "id": "uuid-string",
      "title": "Getting Started with Go",
      "slug": "getting-started-with-go",
      "excerpt": "A beginner's guide to Go programming...",
      "content": "# Getting Started with Go\n\n...",
      "author": {
        "id": "uuid-string",
        "name": "John Doe",
        "email": "john@example.com"
      },
      "image_url": "https://example.com/images/blog-1.jpg",
      "tags": ["go", "backend", "tutorial"],
      "views_count": 150,
      "created_at": "2025-05-20T10:30:00Z",
      "updated_at": "2025-05-20T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total_items": 42,
    "total_pages": 5
  }
}
```

**Error Responses**
- `400 Bad Request`: Invalid pagination parameters

---

### Get Single Blog by Slug
**GET** `/api/v1/blogs/:slug`

Retrieve a single blog by its unique slug.

**Path Parameters**
- `slug` (string): Blog slug (e.g., "getting-started-with-go")

**Response** (200 OK)
```json
{
  "id": "uuid-string",
  "title": "Getting Started with Go",
  "slug": "getting-started-with-go",
  "excerpt": "A beginner's guide to Go programming...",
  "content": "# Getting Started with Go\n\n...",
  "author": {
    "id": "uuid-string",
    "name": "John Doe",
    "email": "john@example.com"
  },
  "image_url": "https://example.com/images/blog-1.jpg",
  "tags": ["go", "backend", "tutorial"],
  "views_count": 150,
  "created_at": "2025-05-20T10:30:00Z",
  "updated_at": "2025-05-20T10:30:00Z"
}
```

**Error Responses**
- `404 Not Found`: Blog not found

---

### Create Blog 🔒
**POST** `/api/v1/blogs`

Create a new blog post. Requires authentication.

**Request Body**
```json
{
  "title": "Getting Started with Go",
  "excerpt": "A beginner's guide to Go programming...",
  "content": "# Getting Started with Go\n\nMarkdown formatted content...",
  "image_url": "https://example.com/images/blog-1.jpg",
  "tags": ["go", "backend", "tutorial"]
}
```

**Response** (201 Created)
```json
{
  "id": "uuid-string",
  "title": "Getting Started with Go",
  "slug": "getting-started-with-go",
  "excerpt": "A beginner's guide to Go programming...",
  "content": "# Getting Started with Go\n\n...",
  "author": {
    "id": "uuid-string",
    "name": "John Doe"
  },
  "image_url": "https://example.com/images/blog-1.jpg",
  "tags": ["go", "backend", "tutorial"],
  "created_at": "2025-05-25T10:30:00Z"
}
```

**Error Responses**
- `400 Bad Request`: Invalid input, missing required fields
- `401 Unauthorized`: Missing or invalid token

---

### Update Blog 🔒
**PUT** `/api/v1/blogs/:id`

Update an existing blog. Only author can update their own blog.

**Path Parameters**
- `id` (string): Blog UUID

**Request Body**
```json
{
  "title": "Updated Title",
  "excerpt": "Updated excerpt...",
  "content": "Updated markdown content...",
  "image_url": "https://example.com/images/updated.jpg",
  "tags": ["updated", "tags"]
}
```

**Response** (200 OK)
```json
{
  "id": "uuid-string",
  "title": "Updated Title",
  "slug": "updated-title",
  "excerpt": "Updated excerpt...",
  "content": "Updated markdown content...",
  "author": {
    "id": "uuid-string",
    "name": "John Doe"
  },
  "image_url": "https://example.com/images/updated.jpg",
  "tags": ["updated", "tags"],
  "updated_at": "2025-05-25T15:45:00Z"
}
```

**Error Responses**
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Not authorized to update this blog
- `404 Not Found`: Blog not found

---

### Delete Blog 🔒
**DELETE** `/api/v1/blogs/:id`

Delete a blog. Only author can delete their own blog.

**Path Parameters**
- `id` (string): Blog UUID

**Response** (204 No Content)

**Error Responses**
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Not authorized to delete this blog
- `404 Not Found`: Blog not found

---

### Upload Blog Image 🔒
**POST** `/api/v1/blogs/upload`

Upload an image for blog post. Supports JPG, PNG, WebP (max 5MB).

**Request Body** (multipart/form-data)
- `file` (file): Image file

**Response** (200 OK)
```json
{
  "url": "https://example.com/images/uploads/blog-image-123.jpg",
  "size": 245120,
  "mime_type": "image/jpeg"
}
```

**Error Responses**
- `400 Bad Request`: No file provided or invalid format
- `413 Payload Too Large`: File exceeds 5MB limit
- `401 Unauthorized`: Missing or invalid token

---

## Content Requirements

### Title
- Required
- Min: 3 characters, Max: 200 characters
- Auto-generates slug (lowercase, kebab-case)

### Excerpt
- Required
- Min: 10 characters, Max: 500 characters
- Brief summary for listing

### Content
- Required (Markdown format)
- Min: 50 characters
- Supports all markdown features

### Tags
- Optional
- Max: 5 tags per blog
- Each tag: 2-30 characters

---

## Response Codes Summary

| Code | Meaning |
|------|---------|
| 200  | Success |
| 201  | Created |
| 204  | No Content (deleted) |
| 400  | Bad Request |
| 401  | Unauthorized |
| 403  | Forbidden |
| 404  | Not Found |
| 409  | Conflict |
| 413  | Payload Too Large |
| 500  | Server Error |
