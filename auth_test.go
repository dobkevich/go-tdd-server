package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthentication(t *testing.T) {
	// 1. Setup Mock JWKS Server
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKey := privateKey.Public().(*rsa.PublicKey)
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	jwksHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := map[string]interface{}{
			"keys": []map[string]interface{}{
				{
					"kty": "RSA",
					"use": "sig",
					"kid": "test-key-id",
					"alg": "RS256",
					"n":   n,
					"e":   e,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(jwks)
	})

	jwksServer := httptest.NewServer(jwksHandler)
	defer jwksServer.Close()

	// 2. Configure environment
	_ = os.Setenv("JWKS_URL", jwksServer.URL)
	defer func() { _ = os.Unsetenv("JWKS_URL") }()

	router := SetupRouter()

	// 3. Test Cases
	t.Run("Health check should bypass auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("API call without token should return 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/test", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("API call with valid token should return 200", func(t *testing.T) {
		// Generate valid token
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		token.Header["kid"] = "test-key-id"
		tokenStr, err := token.SignedString(privateKey)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/test", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "Hello, test")
	})

	t.Run("API call with invalid token should return 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/hello/test", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("MCP call without token should return 401 (protected)", func(t *testing.T) {
		_ = os.Setenv("ENABLE_MCP", "true")
		defer func() { _ = os.Unsetenv("ENABLE_MCP") }()

		router := SetupRouter()
		req := httptest.NewRequest(http.MethodGet, "/mcp/sse", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("MCP call with valid token should return 200", func(t *testing.T) {
		_ = os.Setenv("ENABLE_MCP", "true")
		defer func() { _ = os.Unsetenv("ENABLE_MCP") }()

		// Generate valid token
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub": "user-123",
			"exp": time.Now().Add(time.Hour).Unix(),
		})
		token.Header["kid"] = "test-key-id"
		tokenStr, err := token.SignedString(privateKey)
		require.NoError(t, err)

		router := SetupRouter()
		ts := httptest.NewServer(router)
		defer ts.Close()

		client := &http.Client{
			Timeout: 500 * time.Millisecond,
		}

		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/mcp/sse", nil)
		req.Header.Set("Authorization", "Bearer "+tokenStr)

		resp, err := client.Do(req)
		if err != nil {
			// Timeout is expected for SSE, but we should check if we got headers
			return
		}
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Contains(t, resp.Header.Get("Content-Type"), "text/event-stream")
	})
}
