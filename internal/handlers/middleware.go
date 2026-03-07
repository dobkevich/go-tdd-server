package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates the JWT token against a JWKS endpoint (e.g., Authentik)
func JWTMiddleware(jwksURL string) echo.MiddlewareFunc {
	if jwksURL == "" {
		slog.Warn("JWKS_URL is not set. JWT authentication is disabled.")
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	// Create the keyfunc which will fetch and cache JWKS from the URL
	k, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		// We panic here because if JWKS_URL is provided but invalid, the server shouldn't start in an insecure state
		panic(fmt.Sprintf("Failed to create keyfunc from JWKS URL: %v", err))
	}

	slog.Info("JWT Authentication enabled", "jwks_url", jwksURL)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip auth for health checks, ready checks, root, docs, and MCP SSE
			path := c.Path()
			if path == "/healthz" || path == "/readyz" || path == "/" || strings.HasPrefix(path, "/docs") || strings.HasPrefix(path, "/mcp/sse") {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
			}

			tokenStr := parts[1]

			// Parse and validate the token using the keys fetched from JWKS
			token, err := jwt.Parse(tokenStr, k.Keyfunc)
			if err != nil || !token.Valid {
				slog.Error("JWT validation failed", "error", err)
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired token")
			}

			// Optionally store the token/claims in the context for downstream use
			c.Set("user", token)

			return next(c)
		}
	}
}
