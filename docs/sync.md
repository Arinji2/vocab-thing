# Sync API Documentation

## Base URL

```
https://api-vocabthing.arinji.com
```

## Authentication

All endpoints require an authenticated user session. This is taken from the cookies

---

## Endpoints

### Get User Sync Data

**Endpoint:**

```
GET /sync
```

**Response:**

```json
{
  "id": "string",
  "userId": "string",
  "lastUpdatedAt": "timestamp"
}
```

---

### Manually Sync

> Must have a 30 minute delay between syncs

**Endpoint:**

```
POST /phrases
```

**Response:**

```json
{
  "id": "string",
  "userId": "string",
  "lastUpdatedAt": "timestamp"
}
```

---
