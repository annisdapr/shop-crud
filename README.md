# Observability Integration with Grafana for shop-crud Project

This document provides a comprehensive guide for the `shop-crud` project. It covers the initial setup and demonstrates how to integrate a complete observability stack using Grafana for visualization, Grafana Tempo for distributed tracing, and Grafana Loki with Promtail for structured logging.

## About Project
### Project Structure

The project follows a clean architecture pattern with clear separation of concerns. Here's the detailed structure:

```
codex-apm-o2-local/
├── docker-compose.yaml          # Main Docker Compose file for all services
├── Dockerfile                   # Base Dockerfile
├── README.md                    # This documentation
├── db/
│   └── init.sql                # Database initialization script
│
├── user-service/               # User management microservice
│   ├── docker-compose.yaml    # Service-specific Docker Compose
│   ├── Dockerfile              # Service-specific Dockerfile
│   ├── go.mod                  # Go module dependencies
│   ├── go.sum                  # Go module checksums
│   ├── main.go                 # Service entry point
│   ├── config/
│   │   ├── config.go          # Configuration management
│   │   └── database.go        # Database connection setup
│   ├── docs/
│   │   └── api.md             # API documentation
│   ├── module/
│   │   ├── handlers/          # HTTP request handlers
│   │   │   └── user_handler.go
│   │   ├── models/            # Data models and DTOs
│   │   │   └── user.go
│   │   ├── repositories/      # Data access layer
│   │   │   └── user_repo.go
│   │   └── usecases/          # Business logic layer
│   │       └── user_usecase.go
│   └── pkg/
│       └── tracing/           # OpenTelemetry tracing setup
│           └── tracing.go
│
├── item-service/               # Item/Product management microservice
│   ├── docker-compose.yml     # Service-specific Docker Compose
│   ├── Dockerfile              # Service-specific Dockerfile
│   ├── go.mod                  # Go module dependencies
│   ├── go.sum                  # Go module checksums
│   ├── main.go                 # Service entry point
│   ├── config/
│   │   ├── config.go          # Configuration management
│   │   └── database.go        # Database connection setup
│   ├── middleware/
│   │   └── auth.go            # Authentication middleware
│   ├── modules/
│   │   ├── handlers/          # HTTP request handlers
│   │   │   └── item_handler.go
│   │   ├── models/            # Data models and DTOs
│   │   │   └── item.go
│   │   ├── repositories/      # Data access layer
│   │   │   └── item_repo.go
│   │   └── usecases/          # Business logic layer
│   │       └── item_usecase.go
│   └── tracing/               # OpenTelemetry tracing setup
│       └── tracing.go
│
└── purchase-service/           # Purchase/Transaction management microservice
    ├── docker-compose.yml     # Service-specific Docker Compose
    ├── Dockerfile              # Service-specific Dockerfile
    ├── go.mod                  # Go module dependencies
    ├── go.sum                  # Go module checksums
    ├── main.go                 # Service entry point
    ├── config/
    │   ├── config.go          # Configuration management
    │   └── database.go        # Database connection setup
    ├── db/
    │   └── init.sql           # Service-specific database initialization
    ├── middleware/
    │   └── auth.go            # Authentication middleware
    ├── modules/
    │   ├── handlers/          # HTTP request handlers
    │   │   └── purchase_handler.go
    │   ├── models/            # Data models and DTOs
    │   │   └── purchase.go
    │   ├── repositories/      # Data access layer
    │   │   └── purchase_repo.go
    │   └── usecases/          # Business logic layer
    │       └── purchase_usecase.go
    └── pkg/
        └── tracing/           # OpenTelemetry tracing setup
            └── tracing.go
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

## API Documentation

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

- **User Service**: `8081`
- **Item Service**: `8082`
- **Purchase Service**: `8083`
- **OpenObserve UI**: `5080`
- **OpenObserve Logs**: `5081`
- **OpenObserve Metrics**: `5082`
- **OpenObserve OTLP Traces**: `5083`


## Adding Distributed Tracing with Grafana Tempo


### Step 1: Create Docker Compose Files for Grafana & Tempo
The first step is to define the Grafana (for visualization) and Tempo (for trace storage) services in their respective docker-compose files. This allows us to manage the observability stack separately from the application stack.

Create a file named docker-compose-grafana.yaml inside the tracing-compose folder. This file will run the Grafana service. Ensure its content is as follows:

```yaml
services:
  grafana:
    image: grafana/grafana-oss:latest
    container_name: grafana
    restart: always
    ports:
      - "3001:3000"
    volumes:
      - ./grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
    networks:
      - microservices-net

networks:
  microservices-net:
    external: true
```

Next, create a file named docker-compose-tempo.yaml in the same folder. This file will run the Grafana Tempo service. Ensure its content is as follows:

```yaml
services:
  tempo:
    image: grafana/tempo:latest
    container_name: tempo
    restart: always
    ports:
      - "3200:3200"  # tempo query
      - "4317:4317"  # OTLP gRPC receiver
      - "4318:4318"  # OTLP HTTP receiver
    command: [ "-config.file=/etc/tempo.yaml" ]
    volumes:
      - ./tempo.yaml:/etc/tempo.yaml
      - ./tempo-data:/var/tempo
    networks:
      - microservices-net

networks:
  microservices-net:
    external: true

```

### Step 2: Create the Tempo Configuration File
After defining the services, we need to create a configuration file to tell Tempo how to receive and store trace data.

tempo.yaml File
Create a file named tempo.yaml inside the tracing-compose folder. This is the basic configuration for Tempo. Ensure its content is as follows:

```yaml
  auth_enabled: false

  server:
    http_listen_port: 3200
    log_level: info

  distributor:
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317 
          http:
            endpoint: 0.0.0.0:4318 

  ingester:
    trace_idle_period: 10s
    max_block_bytes: 1_048_576
    max_block_duration: 5m

  compactor:
    compaction:
      block_retention: 1h

  storage:
    trace:
      backend: local
      local:
        path: /tmp/tempo
      wal:
        path: /tmp/tempo/wal
```


## Adding Logging with Grafana Loki and Promtail

### Step 1: Create Docker Compose Files for Loki & Promtail

Just like with Tempo and Grafana, we will define Loki and Promtail in their own respective `docker-compose` files.

Create a file named `docker-compose-loki.yaml` inside the `tracing-compose` folder. This file will run the Grafana Loki service. Ensure its content is as follows:

```yaml
services:
  loki:
    image: grafana/loki:latest
    container_name: loki
    restart: unless-stopped
    command: ["-config.file=/etc/loki/config.yaml"]
    volumes:
      - ./loki-config.yaml:/etc/loki/config.yaml
    ports:
      - "3100:3100" # Loki port
    networks:
      - microservices-net

networks:
  microservices-net:
    external: true
```

Next, create a file named docker-compose-promtail.yaml in the same folder. This file will run Promtail, the log collection agent. Ensure its content is as follows:

```yaml
services:
  promtail:
    image: grafana/promtail:latest
    container_name: promtail
    restart: unless-stopped
    user: root
    group_add:
      - "999" # GANTI 999 DENGAN ID GRUP DOCKER ANDA (jika diperlukan)
    command: ["-config.file=/etc/promtail/config.yaml"]
    volumes:
      - ./promtail-config.yaml:/etc/promtail/config.yaml
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock
    networks:
      - microservices-net

networks:
  microservices-net:
    external: true

```

### Step 2: Create Configuration Files for Loki & Promtail


Create a file named loki-config.yaml inside the tracing-compose folder. This is the basic configuration for Loki. Ensure its content is as follows:

```yaml
auth_enabled: false

server:
  http_listen_port: 3100

common:
  path_prefix: /tmp/loki
  storage:
    filesystem:
      chunks_directory: /tmp/loki/chunks
      rules_directory: /tmp/loki/rules
  replication_factor: 1
  ring:
    instance_addr: 127.0.0.1
    kvstore:
      store: inmemory

schema_config:
  configs:
    - from: 2020-10-24
      store: boltdb-shipper
      object_store: filesystem
      schema: v11
      index:
        prefix: index_
        period: 24h

limits_config:
  allow_structured_metadata: false

ruler:
  alertmanager_url: http://localhost:9093
```

Finally, create a file named promtail-config.yaml in the same folder. This file tells Promtail to automatically discover logs from all running Docker containers. Ensure its content is as follows:

```yaml
server:
  http_listen_port: 9080
  grpc_listen_port: 0

positions:
  filename: /tmp/positions.yaml

clients:
  - url: http://loki:3100/loki/api/v1/push

scrape_configs:
  - job_name: containers
    docker_sd_configs:
      - host: unix:///var/run/docker.sock
        refresh_interval: 5s
    relabel_configs:
      - source_labels: ['__meta_docker_container_name']
        regex: '/(.*)'
        target_label: 'container'
```

## Example Tracing and Logging in a Go Handler


```go
func (h *UserHandler) Register(c echo.Context) error {
   var req models.RegisterRequest
   // start tracing span for RegisterHandler
   tracer := otel.Tracer("user-service-handler")
   ctx, span := tracer.Start(c.Request().Context(), "RegisterHandler")
   defer span.End()
   // set route and email attributes after binding
	
	// 1. Binding request body ke struct.
   if err := c.Bind(&req); err != nil {
		logger.Error(ctx, "Failed to bind register request: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}
   // annotate span with attributes
   span.SetAttributes(
       attribute.String("http.route", c.Path()),
       attribute.String("user.email", req.Email),
   )
	if err := c.Validate(&req); err != nil {
		logger.Error(ctx, "Validation failed on register: "+err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	logger.Info(ctx, "Attempting to register user: "+req.Email)
   user, err := h.userUsecase.Register(ctx, req)
	if err != nil {
		if errors.Is(err, usecases.ErrEmailExists) {
			logger.Warn(ctx, "Registration failed, email already exists: "+req.Email)
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()}) // 409 Conflict
		}
		logger.Error(ctx, "Internal server error on register: "+err.Error())
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to register user"})
	}
	logger.Info(ctx, "✅ Register user success: "+user.Email)
	return c.JSON(http.StatusCreated, user) // 201 Created
}
```
