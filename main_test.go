package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/project/go-tdd-server/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	e := SetupRouter()

	t.Run("Healthz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var res models.HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "ok", res.Status)
		assert.NotEmpty(t, res.Timestamp)
		assert.NotEmpty(t, res.Uptime)
	})

	t.Run("Readyz", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		var res models.HealthResponse
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.Equal(t, "ready", res.Status)
	})
}

func TestV1Endpoints(t *testing.T) {
	e := SetupRouter()

	t.Run("Ping", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "pong", rec.Body.String())
	})

	t.Run("Add", func(t *testing.T) {
		tests := []struct {
			name           string
			url            string
			expectedStatus int
			expectedResult int
		}{
			{"Valid", "/api/v1/add?a=10&b=5", http.StatusOK, 15},
			{"Invalid Param", "/api/v1/add?a=abc&b=5", http.StatusBadRequest, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodGet, tt.url, nil)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				if tt.expectedStatus == http.StatusOK {
					var res models.AddResponse
					err := json.Unmarshal(rec.Body.Bytes(), &res)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedResult, res.Result)
				}
			})
		}
	})

	t.Run("Echo", func(t *testing.T) {
		tests := []struct {
			name           string
			payload        string
			expectedMsg    string
			expectedStatus int
		}{
			{
				"Success",
				`{"message": "Hello TDD!"}`,
				"Hello TDD!",
				http.StatusOK,
			},
			{
				"Validation Empty",
				`{"message": ""}`,
				"",
				http.StatusBadRequest,
			},
			{
				"Validation Too Long",
				`{"message": "This message is definitely longer than one hundred characters, which is the maximum limit we have set for our validation rules in the model structure."}`,
				"",
				http.StatusBadRequest,
			},
			{
				"Invalid JSON",
				`{"message": "unclosed json}`,
				"",
				http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodPost, "/api/v1/echo", strings.NewReader(tt.payload))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				e.ServeHTTP(rec, req)

				assert.Equal(t, tt.expectedStatus, rec.Code)

				if tt.expectedStatus == http.StatusOK {
					var res models.EchoResponse
					err := json.Unmarshal(rec.Body.Bytes(), &res)
					assert.NoError(t, err)
					assert.Equal(t, tt.expectedMsg, res.Message)
					assert.NotEmpty(t, res.Timestamp)
				}
			})
		}
	})
}
