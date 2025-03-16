# Phrase API Documentation

## Base URL

```
https://api-vocabthing.arinji.com
```

## Authentication

All endpoints require an authenticated user session. This is taken from the cookies

---

## Endpoints

### Create a Phrase

**Endpoint:**

```
POST /phrase/create
```

**Request Body:**

```json
{
  "phrase": "string",
  "definition": "string",
  "foundIn": "string",
  "public": true
}
```

**Response:**

```json
{
  "id": "string",
  "userID": "string",
  "phrase": "string",
  "definition": "string",
  "foundIn": "string",
  "public": true,
  "createdAt": "timestamp"
}
```

---

### Create a Phrase Tag

**Endpoint:**

```
POST /phrase/tag/create
```

**Request Body:**

```json
{
  "phraseID": "string",
  "tagName": "string",
  "tagColor": "string"
}
```

**Response:**

```json
{
  "id": "string",
  "phraseID": "string",
  "tagName": "string",
  "tagColor": "string",
  "createdAt": "timestamp"
}
```

---

### Get a Phrase by ID

**Endpoint:**

```
GET /phrase/{id}
```

**Response:**

```json
{
  "id": "string",
  "userID": "string",
  "phrase": "string",
  "definition": "string",
  "foundIn": "string",
  "public": true,
  "createdAt": "timestamp"
}
```

---

### Get All Phrases

**Endpoint:**

```
GET /phrases
```

**Query Parameters:**

- `page` (int) - Page number (default: 1)
- `pageSize` (int) - Number of results per page (default: 10, max: 100)
- `sortBy` (string) - Sorting field (`createdAt`, `usageCount`, default: `createdAt`)
- `order` (string) - Sorting order (`ASC`, `DESC`, default: `DESC`)
- `groupBy` (string) - Grouping method (`foundIn`, `public`, default: `foundIn`)

**Response:**

```json
[
  {
    "id": "string",
    "phrase": "string",
    "definition": "string",
    "foundIn": "string",
    "public": true,
    "createdAt": "timestamp"
  }
]
```

---

### Search Phrases

**Endpoint:**

```
GET /phrase/search
```

**Query Parameters:**

- `term` (string) - Search term (required)
- `page` (int) - Page number (default: 1)
- `pageSize` (int) - Number of results per page (default: 10, max: 100)
- `sortBy` (string) - Sorting field (`createdAt`, `usageCount`, default: `createdAt`)
- `order` (string) - Sorting order (`ASC`, `DESC`, default: `DESC`)
- `groupBy` (string) - Grouping method (`foundIn`, `public`, default: `foundIn`)

**Response:**

```json
[
  {
    "id": "string",
    "phrase": "string",
    "definition": "string",
    "foundIn": "string",
    "public": true,
    "createdAt": "timestamp"
  }
]
```

---

### Update a Phrase

**Endpoint:**

```
PUT /phrase/{id}
```

**Request Body:**

```json
{
  "phrase": "string",
  "definition": "string",
  "foundIn": "string",
  "public": true
}
```

**Response:**

```json
{
  "id": "string",
  "phrase": "string",
  "definition": "string",
  "foundIn": "string",
  "public": true,
  "updatedAt": "timestamp"
}
```

---

### Update a Phrase Tag

**Endpoint:**

```
PUT /phrase/{phraseID}/tag/{tagID}
```

**Request Body:**

```json
{
  "tagName": "string",
  "tagColor": "string"
}
```

**Response:**

```json
{
  "id": "string",
  "phraseID": "string",
  "tagName": "string",
  "tagColor": "string",
  "updatedAt": "timestamp"
}
```

---

### Delete a Phrase

**Endpoint:**

```
DELETE /phrase/{id}
```

**Response:**

```
204 No Content
```

---

### Delete a Phrase Tag

**Endpoint:**

```
DELETE /phrase/{phraseID}/tag/{tagID}
```

**Response:**

```
204 No Content
```
