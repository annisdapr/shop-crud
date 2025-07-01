# Observability Integration with Codex CLI and OpenObserve for shop-crud Project

This document will guide you through installing the Codex CLI, provide a brief overview of the `shop-crud` project, and demonstrate how to use Codex to add observability functionality using OpenObserve.

## A. Codex CLI Installation and Configuration

This section covers all the necessary steps to install and set up Codex in your local environment.

### 1. Install Codex via NPM

```bash
npm install -g @openai/codex@latest
```
### 2. Create the Configuration File

Create the configuration directory and file by running the following command:

```bash
mkdir -p ~/.codex && cat > ~/.codex/config.json <<'EOF'
{
  "providers": {
    "azure": {
      "name": "AzureOpenAI",
      "baseURL": "https://dpe-open-ai.openai.azure.com/openai",
      "envKey": "AZURE_OPENAI_API_KEY"
    }
  },
  "provider": "azure",
  "model": "o4-mini"
}
```

### 3. Set Up Environment Variables

Add your API key as environment variables. Codex will use these for authentication.

Replace `(your_provided_key)` with your actual API key:

```bash
export AZURE_OPENAI_API_KEY="(your_provided_key)"
export OPENAI_API_KEY="(your_provided_key)"
```

### 4. Run the Codex CLI

Verify your installation by running the following command.  
The `--auto-edit` flag allows Codex to modify your files directly:

```bash
codex -p azure --auto-edit
```
## B. Introduction to the shop-crud Project

`shop-crud` is a Go-based application designed with a microservices architecture. Its purpose is to provide basic CRUD (Create, Read, Update, Delete) functionality for an e-commerce system.

The project is composed of several independent services:

- **User Service**: Manages user data, registration, and authentication.
- **Item Service**: Manages product or item data.
- **Purchase Service**: Manages purchase transactions.

### Main Objective of this Tutorial

The goal of this guide is to demonstrate how the Codex CLI can be used as an AI assistant to automatically add observability features to the `shop-crud` project, using **OpenObserve** as the backend for logs, metrics, and distributed tracing.

## C. Using Codex to Add OpenObserve

Here are example prompts we will use to ask Codex to modify our project and integrate OpenObserve.

---

### 🧠 Prompt 1: Add the OpenObserve Service to Docker Compose

**Prompt Goal:**  
To ask Codex to add the OpenObserve service definition to the main `docker-compose.yaml` file.

**The Prompt:**

Add OpenObserve to my existing Docker Compose file using the latest image, with the service name openobserve. The UI should be on port 5080, logs on 5081, metrics on 5082, and OTLP traces on 5083. Set a root email and password for the dashboard UI. Include restart: unless-stopped, and add a volume that maps ./openobserve_data to /app/data. Do not include any depends_on configuration.


**Expected Outcome:**

- **File to be Modified:** `docker-compose.yaml`  
- **Location:** Root directory of your project  
- **Changes:** Codex will add a new service named `openobserve` to the Docker Compose file with:
  - UI on port `5080`
  - Logs on port `5081`
  - Metrics on port `5082`
  - Traces via OTLP on port `5083`
  - Volume mount: `./openobserve_data:/app/data`
  - Environment variables: root email & password for dashboard
  - `restart: unless-stopped` policy

---

### 🧠 Prompt 2: Create the Tracing Configuration File

**Prompt Goal:**  
To ask Codex to create a new Go file that initializes OpenTelemetry and sends traces to OpenObserve.

**The Prompt:**

Create a file at internal/tracing/tracing.go that contains a function InitTracerProvider(serviceName, collectorHost string) to initialize OpenTelemetry tracing using an OTLP HTTP exporter. The exporter should send traces to /api/default/v1/traces and include an Authorization header using basic auth read from the environment variable OTEL_AUTH_TOKEN, with a default fallback if the variable is not set. Add resource attributes for the service name and environment=development, and use AlwaysSample for sampling. Also include a helper function getEnvOrDefault(key, fallback) to handle the environment variable with fallback logic.

**Expected Outcome:**

- **File to be Created:** `internal/tracing/tracing.go`  
- **Location:** Inside each microservice folder that will use tracing (e.g., `user-service/pkg/tracing/`, `item-service/pkg/tracing/`, etc.)

**File Contents Will Include:**

- ✅ A function: `InitTracerProvider(serviceName, collectorHost string)`
- ✅ A helper: `getEnvOrDefault(key, fallback string)`
- ✅ Full configuration to:
  - Use OTLP HTTP exporter
  - Send data to `/api/default/v1/traces`
  - Attach Authorization header using `OTEL_AUTH_TOKEN` or fallback
  - Include resource attributes like `service.name` and `environment=development`
  - Use `AlwaysSample` to sample all traces

### 🧠 Prompt 3: Add Tracing to The Main Function

**Prompt Goal:**  
Initialize the tracer provider in the main service entry point.

**The Prompt:**
Add tracing initialization to the main function in user-service's main.go using InitTracerProvider("user-service", "openobserve:5080"). If there's an error, log.Fatal it. Also add a deferred shutdown call with context.Background() and log error if it fails.

**Expected Outcome:**

```go
tp, err := tracing.InitTracerProvider("user-service", "openobserve:5080")
if err != nil {
	log.Fatalf("❌ Failed to initialize tracer: %v", err)
}
defer func() {
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Printf("❌ Error shutting down tracer: %v", err)
	}
}()
```
### 🧠 Prompt 4: Add Tracing to RegisterHandler (User Handler)

**Prompt Goal:**  
To instrument the `RegisterHandler` in `user-service`.

**The Prompt:**
Add OpenTelemetry tracing to RegisterHandler in user-service modules. The handler is in user-service/bin/modules/handlers/user_handler.go. Use otel.Tracer("user-service-handler"), name the span "RegisterHandler", set attributes http.route and user.email.

**Expected Outcome:**

```go
ctx := c.Request().Context()
tr := otel.Tracer("user-service-handler")
ctx, span := tr.Start(ctx, "RegisterHandler")
defer span.End()

span.SetAttributes(
	attribute.String("http.route", "/users/register"),
	attribute.String("user.email", req.Email),
)
```

### 🧠 Prompt 5: Add Tracing to Register Usecase

**Prompt Goal:**  
Instrument the core business logic for user registration

**The Prompt:**
Add OpenTelemetry tracing to the Register function in user_usecase.go in user-service. Use otel.Tracer("user-service-usecase"), span name "UserUsecase.Register", and set user.email and user.name as span attributes from the request.

**Expected Outcome:**

```go
tr := otel.Tracer("user-service-usecase")
ctx, span := tr.Start(ctx, "UserUsecase.Register")
defer span.End()

span.SetAttributes(
	attribute.String("user.email", req.Email),
	attribute.String("user.name", req.Name),
)
```
### 🧠 Prompt 6: Add Tracing to Repository Layer (UserRepository.Create)

**Prompt Goal:**  
Trace database interaction at the repository level in the `user-service`.

**The Prompt:**
Add OpenTelemetry tracing to the Create method in user_repo.go of user-service. Use otel.Tracer("user-service-repo"), span name "UserRepository.Create", and add attributes user.id and user.email.

**Expected Outcome:**

```go
tr := otel.Tracer("user-service-repo")
ctx, span := tr.Start(ctx, "UserRepository.Create")
defer span.End()

span.SetAttributes(
	attribute.String("user.id", user.ID.String()),
	attribute.String("user.email", user.Email),
)
```

## 🧪 Tracing Setup and Test Instructions

Follow these steps to start OpenObserve, configure authentication, and run the user service to test tracing.

---

### 🔧 Step 1: Start OpenObserve

Run the following command to build and start the OpenObserve service:

```bash
docker compose up openobserve --build
```
### 🌐 Step 2: Open the Dashboard

Once OpenObserve is running, open your browser and navigate to:

http://localhost:5080

### 🔐 Step 3: Get Basic Auth Token

After logging in with your root email and password, obtain your **Basic Auth token**.  

**Example:**
Basic YWRtaW5AZ21haWwuY29tOlE3U21YNmJ2ZXRiY1V4Y0E=

### 📝 Step 4: Update Your `.env` File

Copy only the value after `Basic` and add it to your `.env` file:

```env
OTEL_AUTH_TOKEN=YWRtaW5AZ21haWwuY29tOlE3U21YNmJ2ZXRiY1V4Y0E=
```


### 🚀 Step 5: Start User Service and Database

Open a **new terminal** and run the following command to start both the user service and PostgreSQL database:

```bash
docker compose up user-service db --build
```
### 🚀 Test the Register Endpoint

**URL:**

http://localhost:{{USER_SERVICE_PORT}}/api/v1/users/register

***Request Payload (JSON):***
```json
{
  "name": "example",
  "email": "example@gmail.com",
  "password": "example123"
}
```
If tracing has been configured correctly, you should see a new trace appear in the OpenObserve dashboard at http://localhost:5080