# Project: Go TDD Server Boilerplate

This project is a high-quality template for Go applications using TDD (Test Driven Development) and Clean Architecture, optimized for Kubernetes and professional CI/CD workflows.

## 🛠 Project Architecture & Module
- **Module Name**: `github.com/project/go-tdd-server`
- **Layers**:
  - `internal/models`: Data structures and validation tags.
  - `internal/service`: Business logic interfaces and implementations.
  - `internal/handlers`: HTTP handlers (Echo) with Dependency Injection.
  - `main.go`: Application bootstrapping and Graceful Shutdown.
  - `docs/`: OpenAPI 3.0 specs and UI (Scalar).

## 🧪 Testing Standards (TDD Mandate)
**RULE #1: Red -> Green -> Refactor.** Any new feature or bug fix MUST start with a test.
- Use **Table-Driven Tests** for comprehensive coverage.
- Test through the router (`e.ServeHTTP`) to verify middlewares and routing.
- Mock external dependencies if necessary (use `testify/mock`).

## 🏗 Engineering & Quality Standards
1. **Automation**: Use the **Makefile** for all standard tasks (`make test`, `make build`, `make lint`).
2. **Linting**: Code must pass `golangci-lint` (config in `.golangci.yml`). Run `make lint` before committing.
3. **CI/CD**: GitHub Actions (`.github/workflows/ci.yml`) automatically runs tests and linting on push.
4. **Dependencies**: Maintain the `vendor/` directory. Run `go mod tidy && go mod vendor` after adding libraries.
5. **Logging**: Use structured JSON logging via `log/slog`.
6. **Configuration**: Environment variables (e.g., `PORT`) must have sane defaults.
7. **Graceful Shutdown**: Always handle `SIGINT` and `SIGTERM` with a 10s timeout.

## 🐳 Docker & Kubernetes
1. **Security**: Run as non-root (`USER 1001`) in Docker.
2. **Health**: Probes at `/healthz` (liveness) and `/readyz` (readiness) are mandatory.
3. **Runtime**: Use `go.uber.org/automaxprocs` to respect container CPU limits.

## 🌐 API Guidelines
1. **Versioning**: Use `/api/v1` prefix for all business logic endpoints.
2. **Validation**: Use declarative validation (`github.com/go-playground/validator/v10`) in models.
3. **Documentation**: Keep `docs/openapi.json` synchronized with code changes.

## 🚀 Key Commands
- `make help` — list all available commands.
- `make test` — run all tests with race detection.
- `make lint` — run code quality checks.
- `make build` — compile the production binary.
- `make run` — build and start the server locally.
- `make docker-build` — build the production Docker image.

---

# 🛠 How-To: Adding New Functionality (REST + MCP)

Follow these steps to add a new endpoint and its MCP tool companion.

## Step 1: Define Business Logic (Service Layer)
Add the method to the `AppService` interface in `internal/service/app.go` and implement it.
```go
// internal/service/app.go
type AppService interface {
    Multiply(ctx context.Context, a, b int) int
}
```

## Step 2: Create HTTP Handler (REST)
Create a handler in `internal/handlers/http/handlers.go` and write a test in `main_test.go`.
```go
// internal/handlers/http/handlers.go
func (h *Handler) Multiply(c echo.Context) error {
    // 1. Bind & Validate input
    // 2. Call h.AppSvc.Multiply
    // 3. Return JSON response
}
```
Register it in `main.go`: `api.GET("/multiply", h.Multiply)`.

## Step 3: Create MCP Tool Companion (AI)
If the feature should be available to AI agents, register it in `internal/handlers/mcp/handlers.go`.
```go
// internal/handlers/mcp/handlers.go

// 1. Define typed arguments with descriptions for LLM
type MultiplyArgs struct {
    A int `json:"a" jsonschema:"First number"`
    B int `json:"b" jsonschema:"Second number"`
}

// 2. Register the tool in registerTools()
mcp.AddTool(h.Server, &mcp.Tool{
    Name:        "multiply",
    Description: "Multiplies two numbers. Use this for math.",
}, func(ctx context.Context, req *mcp.CallToolRequest, args MultiplyArgs) (*mcp.CallToolResult, any, error) {
    res := h.AppSvc.Multiply(ctx, args.A, args.B)
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprint(res)}},
    }, nil, nil
})
```

---

## 🔐 Authentication Management (JWT/OIDC)

### How to Enable
Set the `JWKS_URL` environment variable to your provider's JWKS endpoint.
```bash
export JWKS_URL="https://auth.example.com/jwks"
```
The middleware will automatically protect both `/api/v1/*` and `/mcp/*` routes.

### Client Configuration (MCP)
To connect a protected agent (like Gemini CLI), you must provide the JWT token in the headers. Update your `.gemini/settings.json`:
```json
"mcpServers": {
  "go-tdd-server": {
    "url": "http://localhost:8080/mcp/sse",
    "type": "sse",
    "headers": {
      "Authorization": "Bearer YOUR_JWT_TOKEN"
    }
  }
}
```

### How to Disable
Unset or leave `JWKS_URL` empty. The server will allow all requests (useful for internal networks).

---

## 🤖 MCP (AI Tools) Management

### How to Enable
Set `ENABLE_MCP=true`. This will initialize the MCP Server and open the `/mcp/sse` endpoint.

### Internal-Only Features
To keep a feature internal (not visible to AI), simply **do not** add it to `internal/handlers/mcp/handlers.go`. It will remain a standard REST-only endpoint.
