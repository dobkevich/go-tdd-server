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


## 📖 API Documentation
Once the server is running, visit [http://localhost:8080/docs/](http://localhost:8080/docs/) for interactive API documentation.

## 🧪 Development Philosophy
This project strictly follows **Test Driven Development (TDD)**. Refer to [GEMINI.md](./GEMINI.md) for detailed engineering standards and contribution guidelines.

## 📄 License
This project is licensed under the MIT License.
