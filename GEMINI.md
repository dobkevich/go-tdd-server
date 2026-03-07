@./instructions.md

## 📝 Commit Guidelines
After completing a task that involves file changes, always provide a draft commit message in the following format (but DO NOT perform the commit itself):
- **Format**: `<type>(<scope>): <short description in lowercase>`
- **Types**: `feat`, `fix`, `chore`, `build`, `refactor`, `docs`, `test`, etc.
- **Scope**: The file or component affected (e.g., `(Dockerfile)`, `(handlers)`).
- **Body**: Include a more detailed description of the changes below the subject line.

---

## 🚀 Correct Server Startup (Background with MCP)

To start the server correctly with MCP enabled and ensure it stays running in the background, use the following command:

```bash
fuser -k 8080/tcp || true; go build -o server main.go && ENABLE_MCP=true nohup ./server > server.log 2>&1 &
```

**Key components:**
- `fuser -k 8080/tcp`: Ensures the port is free before starting.
- `go build -o server main.go`: Compiles the binary (more stable for background execution than `go run`).
- `ENABLE_MCP=true`: Environment variable to activate MCP handlers.
- `nohup ... &`: Isolates the process from the terminal and runs it in the background.
- `> server.log 2>&1`: Redirects all logs (stdout and stderr) to a file for debugging.
