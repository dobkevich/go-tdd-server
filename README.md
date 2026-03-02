# Go TDD Server Boilerplate

A production-ready Go server boilerplate built with **Echo**, focusing on **Clean Architecture**, **TDD (Test Driven Development)**, and **Kubernetes** readiness.

## 🚀 Features
- **Clean Architecture**: Decoupled layers (Handlers, Services, Models).
- **TDD Mandate**: Built-in testing standards using Table-Driven Tests.
- **K8s Ready**: Includes `/healthz` and `/readyz` probes, `SIGTERM` handling, and `automaxprocs`.
- **Structured Logging**: JSON logging using `log/slog`.
- **API Versioning**: Standardized `/api/v1` prefix.
- **Documentation**: Embedded OpenAPI UI (Scalar) available at `/docs/`.
- **Dockerized**: Secure multi-stage builds running as non-root.

## 🛠 Prerequisites
- Go 1.24+
- Docker (optional)

## 📦 Getting Started

### Local Development
1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd go-tdd-server
   ```
2. **Run tests (mandatory):**
   ```bash
   go test -v ./...
   ```
3. **Build and run:**
   ```bash
   go build -o server main.go
   ./server
   ```

### Docker
1. **Build the image:**
   ```bash
   docker build -t go-tdd-server .
   ```
2. **Run the container:**
   ```bash
   docker run -p 8080:8080 go-tdd-server
   ```

## 📖 API Documentation
Once the server is running, visit [http://localhost:8080/docs/](http://localhost:8080/docs/) for interactive API documentation.

## 🧪 Development Philosophy
This project strictly follows **TDD**. Refer to [GEMINI.md](./GEMINI.md) for detailed engineering standards and contribution guidelines.

## 📄 License
This project is licensed under the MIT License.
