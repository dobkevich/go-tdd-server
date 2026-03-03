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

## 📝 Commit Guidelines
After completing a task that involves file changes, always provide a draft commit message in the following format (but DO NOT perform the commit itself):
- **Format**: `<type>(<scope>): <short description in lowercase>`
- **Types**: `feat`, `fix`, `chore`, `build`, `refactor`, `docs`, `test`, etc.
- **Scope**: The file or component affected (e.g., `(Dockerfile)`, `(handlers)`).
- **Body**: Include a more detailed description of the changes below the subject line.

