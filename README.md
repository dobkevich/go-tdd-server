# Go TDD Server Boilerplate (AI-Ready with MCP)

A production-ready Go server boilerplate built with **Echo**, focusing on **Clean Architecture**, **TDD (Test Driven Development)**, and **AI Readiness (MCP)**. Optimized for Kubernetes and professional CI/CD.

## 🚀 Features
- **Clean Architecture**: Decoupled layers (Handlers, Services, Models).
- **AI-Ready (MCP)**: Built-in support for **Model Context Protocol (MCP)** via SSE.
- **Secure by Default**: JWT Authentication with OIDC/JWKS integration (supports Authentik, Okta, Auth0, etc.).
- **TDD Mandate**: Strict testing standards using Table-Driven Tests.
- **K8s Ready**: Includes `/healthz` and `/readyz` probes, `SIGTERM` handling, and `automaxprocs`.
- **Dockerized**: Secure multi-stage builds running as non-root.
- **Structured Logging**: JSON logging using `log/slog`.
- **API Versioning**: Standardized `/api/v1` prefix.
- **Documentation**: Embedded OpenAPI UI (Scalar) available at `/docs/`.

## 🛠 Configuration (Environment Variables)

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server listening port | `8080` |
| `ENABLE_MCP` | Enable Model Context Protocol (AI tools) | `false` |
| `JWKS_URL` | JWKS endpoint for JWT validation (OIDC provider) | (Empty = Disabled) |

## 📦 Getting Started

### Local Development
1. **Clone the repository:**
   ```bash
   git clone https://github.com/dobkevich/go-tdd-server.git
   cd go-tdd-server
   ```
2. **Run tests (mandatory):**
   ```bash
   make test
   ```
3. **Run code quality checks:**
   ```bash
   make lint
   ```
4. **Build and run:**
   ```bash
   make run
   ```

### Docker
1. **Build the image:**
   ```bash
   make docker-build
   ```
2. **Run the container:**
   ```bash
   make docker-run
   ```

### Using MCP (AI Tools)
To enable MCP for AI agents (like Claude Desktop or Cursor):
```bash
export ENABLE_MCP=true
make run
```
The MCP endpoint will be available at `http://localhost:8080/mcp/sse`.

### Enabling Authentication
Set the JWKS URL from your OIDC provider (e.g., Authentik):
```bash
export JWKS_URL="https://your-auth-provider.com/application/o/app/jwks/"
make run
```

## 📖 API Documentation
Once the server is running, visit [http://localhost:8080/docs/](http://localhost:8080/docs/) for interactive API documentation.

## 🧪 Development Philosophy
This project strictly follows **Test Driven Development (TDD)**. Refer to [instructions.md](./instructions.md) for detailed engineering standards and contribution guidelines.

### 📖 Developer Guides
- **[How-To: Add REST + MCP Features](./instructions.md#how-to-adding-new-functionality-rest-mcp)**
- **[Auth & Security Management](./instructions.md#authentication-management-jwtoidc)**
- **[MCP Server Management](./instructions.md#mcp-ai-tools-management)**

## 📄 License
This project is licensed under the MIT License.
