package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMCPModularity(t *testing.T) {
	t.Run("MCP Disabled by Default", func(t *testing.T) {
		_ = os.Unsetenv("ENABLE_MCP")
		e := SetupRouter()
		req := httptest.NewRequest(http.MethodGet, "/mcp/sse", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("MCP Tool Echo", func(t *testing.T) {
		_ = os.Setenv("ENABLE_MCP", "true")
		defer func() { _ = os.Unsetenv("ENABLE_MCP") }()

		e := SetupRouter()
		ts := httptest.NewServer(e)
		defer ts.Close()

		// 1. GET request for SSE
		resp, err := http.Get(ts.URL + "/mcp/sse")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		defer func() { _ = resp.Body.Close() }()

		scanner := bufio.NewScanner(resp.Body)
		var sessionID string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "sessionid=") {
				parts := strings.Split(line, "sessionid=")
				sessionID = parts[1]
				break
			}
		}
		assert.NotEmpty(t, sessionID)

		// 2. Initialize MCP
		initPayload := `{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}},"id":1}`
		_, _ = http.Post(ts.URL+"/mcp/sse?sessionid="+sessionID, "application/json", strings.NewReader(initPayload))

		// Skip initialize response
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), `"result":`) {
				break
			}
		}

		// 3. Call Echo Tool
		testMessage := "Hello MCP Test"
		payload := fmt.Sprintf(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"echo","arguments":{"message":"%s"}},"id":2}`, testMessage)
		postResp, err := http.Post(ts.URL+"/mcp/sse?sessionid="+sessionID, "application/json", strings.NewReader(payload))
		assert.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, postResp.StatusCode)

		// 4. Verify result in stream
		found := false
		done := make(chan bool)
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, testMessage) {
					found = true
					done <- true
					return
				}
			}
			done <- false
		}()

		select {
		case <-done:
			assert.True(t, found, "Echo result should be in SSE stream")
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for echo response")
		}
	})
}
