## A. Introduction to the shop-crud Project

`shop-crud` is a Go-based application designed with a microservices architecture. Its purpose is to provide basic CRUD (Create, Read, Update, Delete) functionality for an e-commerce system.

The project is composed of several independent services:

- **User Service**: Manages user data, registration, and authentication.
- **Item Service**: Manages product or item data.
- **Purchase Service**: Manages purchase transactions.

### Project Structure

The project follows a clean architecture pattern with clear separation of concerns. Here's the detailed structure:

```
├───db                            # Shared SQL initialization or migration files
│       init.sql
│
├───Doc                           # Documentation files (general project documentation)
│
├───item-service                 # Microservice handling item-related operations
│   │   docker-compose.yml       # Docker Compose file for item service
│   │   Dockerfile               # Dockerfile to build the item service container
│   │   go.mod                   # Go module file for dependencies
│   │   go.sum                   # Dependency checksums
│   │   main.go                  # Entry point for item service
│   │
│   ├───config                   # Configuration handling (env, DB setup, etc.)
│   │       config.go
│   │       database.go
│   │
│   ├───db                       # Item service specific database initialization
│   │   └───init.sql
│   ├───doc                      # Swagger or additional API documentation (optional)
│   ├───middleware               # Middleware functions (e.g., authentication)
│   │       auth.go
│   │
│   └───modules                  # Business logic for item service
│       ├───handlers             # HTTP handlers for item endpoints
│       │       item_handler.go
│       │
│       ├───models               # Data models for item entities
│       │       item.go
│       │
│       ├───repositories         # Data access layer (database operations)
│       │       item_repo.go
│       │
│       └───usecases             # Application logic and use case implementations
│               item_usecase.go
│
├───purchase-service             # Microservice for managing purchases
│   │   docker-compose.yml       # Docker Compose file for purchase service
│   │   Dockerfile               # Dockerfile to build the purchase service container
│   │   go.mod
│   │   go.sum
│   │   main.go
│   │
│   ├───config                   # Configuration files
│   │       config.go
│   │       database.go
│   │
│   ├───db                       # Purchase-specific database init script
│   │       init.sql
│   │
│   ├───middleware               # Middleware (auth, etc.)
│   │       auth.go
│   │
│   ├───modules
│   │   ├───handlers             # HTTP handlers for purchase endpoints
│   │   │       purchase_handler.go
│   │   │
│   │   ├───models               # Data models for purchase entities
│   │   │       purchase.go
│   │   │
│   │   ├───repositories         # Data access layer
│   │   │       purchase_repo.go
│   │   │
│   │   └───usecases             # Business logic and use cases
│   │           purchase_usecase.go
│   │
│   └───pkg                      # Utility packages (common helper functions)
│
└───user-service                 # Microservice handling user registration and authentication
    │   .env.user.example        # Example environment config
    │   docker-compose.yaml      # Docker Compose file for user service
    │   Dockerfile               # Dockerfile for user service
    │   go.mod
    │   go.sum
    │   main.go
    │
    ├───config                   # Configuration files
    │       config.go
    │       database.go
    │
    ├───docs                     # API documentation for user service
    │       api.md
    │
    ├───module
    │   ├───handlers             # HTTP handlers for user-related endpoints
    │   │       user_handler.go
    │   │
    │   ├───models               # Data models for user entities
    │   │       user.go
    │   │
    │   ├───repositories         # Data access layer for users
    │   │       user_repo.go
    │   │
    │   └───usecases             # User-related business logic
    │           user_usecase.go
    │
    └───pkg                      # Utility or shared helper code for user service

```

### Architecture Overview

Each microservice follows the **Clean Architecture** pattern with the following layers:

1. **Handlers Layer** (`handlers/`): HTTP request/response handling, input validation, and routing
2. **Use Cases Layer** (`usecases/`): Business logic and orchestration
3. **Repository Layer** (`repositories/`): Data access and persistence
4. **Models Layer** (`models/`): Data structures, DTOs, and domain entities

### Key Features

- **Microservices Architecture**: Independent, loosely-coupled services
- **Clean Architecture**: Separation of concerns with clear layer boundaries
- **JWT Authentication**: Secure token-based authentication system
- **Database Integration**: PostgreSQL with proper connection pooling
- **Docker Support**: Containerized deployment with Docker Compose
- **Observability Ready**: OpenTelemetry tracing integration
- **Input Validation**: Request validation using struct tags
- **Error Handling**: Proper HTTP status codes and error responses

## D. API Documentation

This section provides comprehensive documentation for all available API endpoints across the three microservices.

### Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### User Service API

The User Service handles user registration, authentication, and user management.

**Base URL**: `http://localhost:8081/api/v1`

#### POST /users/register
Register a new user account.

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Validation Rules:**
- `name`: Required
- `email`: Required, must be valid email format
- `password`: Required, minimum 8 characters

**Responses:**
- `201 Created`: User successfully registered
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-01T10:00:00Z"
}
```

- `400 Bad Request`: Validation error
- `409 Conflict`: Email already exists
- `500 Internal Server Error`: Server error

#### POST /users/login
Authenticate user and get access token.

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Responses:**
- `200 OK`: Login successful
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

- `400 Bad Request`: Validation error
- `401 Unauthorized`: Invalid credentials
- `500 Internal Server Error`: Server error

### Item Service API

The Item Service manages product/item data with full CRUD operations.

**Base URL**: `http://localhost:8082/api/v1`

#### GET /items
Get all items (public endpoint).

**Query Parameters:**
- `limit` (optional): Number of items per page
- `offset` (optional): Number of items to skip

**Responses:**
- `200 OK`: List of items
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Laptop Gaming",
    "description": "High-performance gaming laptop",
    "price": 1500.00,
    "stock": 10,
    "created_at": "2025-01-01T10:00:00Z",
    "updated_at": "2025-01-01T10:00:00Z"
  }
]
```

#### GET /items/:id
Get item by ID (public endpoint).

**Path Parameters:**
- `id`: Item UUID

**Responses:**
- `200 OK`: Item details
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "name": "Laptop Gaming",
  "description": "High-performance gaming laptop",
  "price": 1500.00,
  "stock": 10,
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-01T10:00:00Z"
}
```

- `400 Bad Request`: Invalid item ID format
- `404 Not Found`: Item not found
- `500 Internal Server Error`: Server error

#### POST /items
Create a new item (requires authentication).

**Request Body:**
```json
{
  "name": "Laptop Gaming",
  "description": "High-performance gaming laptop",
  "price": 1500.00,
  "stock": 10
}
```

**Validation Rules:**
- `name`: Required, minimum 3 characters
- `description`: Optional
- `price`: Required, must be >= 0
- `stock`: Required, must be >= 0

**Responses:**
- `201 Created`: Item successfully created
- `400 Bad Request`: Validation error
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

#### PUT /items/:id
Update an existing item (requires authentication).

**Path Parameters:**
- `id`: Item UUID

**Request Body:**
```json
{
  "name": "Updated Laptop Gaming",
  "description": "Updated high-performance gaming laptop",
  "price": 1600.00,
  "stock": 8
}
```

**Responses:**
- `200 OK`: Item successfully updated
- `400 Bad Request`: Validation error or invalid ID format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Item not found
- `500 Internal Server Error`: Server error

#### DELETE /items/:id
Delete an item (requires authentication).

**Path Parameters:**
- `id`: Item UUID

**Responses:**
- `204 No Content`: Item successfully deleted
- `400 Bad Request`: Invalid item ID format
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: Item not found
- `500 Internal Server Error`: Server error

### Purchase Service API

The Purchase Service handles transaction creation and management.

**Base URL**: `http://localhost:8083/api/v1`

#### POST /purchases
Create a new purchase (requires authentication).

**Request Body:**
```json
{
  "items": [
    {
      "item_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 2
    },
    {
      "item_id": "550e8400-e29b-41d4-a716-446655440002",
      "quantity": 1
    }
  ]
}
```

**Validation Rules:**
- `items`: Required, must have at least 1 item
- `item_id`: Required, must be valid UUID
- `quantity`: Required, must be > 0

**Responses:**
- `201 Created`: Purchase successfully created
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "total_amount": 3100.00,
  "created_at": "2025-01-01T10:00:00Z",
  "items": [
    {
      "item_id": "550e8400-e29b-41d4-a716-446655440001",
      "quantity": 2,
      "name": "Laptop Gaming",
      "price": 1500.00
    },
    {
      "item_id": "550e8400-e29b-41d4-a716-446655440002",
      "quantity": 1,
      "name": "Mouse Gaming",
      "price": 100.00
    }
  ]
}
```
**Base URL**: `http://localhost:8083/api/v1`

#### GET /purchases
Get user's purchase history (requires authentication).

**Responses:**
- `200 OK`: Successfully retrieved purchase history
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "total_amount": 3100.00,
    "created_at": "2025-01-01T10:00:00Z",
    "items": [
      {
        "item_id": "550e8400-e29b-41d4-a716-446655440001",
        "quantity": 2,
        "name": "Laptop Gaming",
        "price": 1500.00
      }
    ]
  }
]
```

- `400 Bad Request`: Validation error
- `401 Unauthorized`: Missing or invalid token
- `409 Conflict`: Item not found or insufficient stock
- `500 Internal Server Error`: Server error

### Error Response Format

All endpoints return errors in a consistent format:

```json
{
  "error": "Descriptive error message"
}
```

### HTTP Status Codes

- `200 OK`: Request successful
- `201 Created`: Resource created successfully
- `204 No Content`: Resource deleted successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required or invalid
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., email already exists)
- `500 Internal Server Error`: Server error

### Service Ports

- **User Service**: `5000`
- **Item Service**: `5001`
- **Purchase Service**: `5002`