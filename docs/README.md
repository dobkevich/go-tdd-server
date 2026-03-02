# API Documentation

This directory contains the OpenAPI 3.0 specification for the Go TDD Server.

## 📄 Contents
- **`openapi.json`**: The core API specification.
- **`index.html`**: A static page using [Scalar](https://github.com/scalar/scalar) to render the `openapi.json` into a beautiful, interactive documentation UI.

## 🚀 Serving Documentation
The documentation is automatically served by the application at the `/docs/` path.
Whenever you add or modify endpoints, update `openapi.json` to keep the UI in sync.
