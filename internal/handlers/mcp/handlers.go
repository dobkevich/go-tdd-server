package mcphandlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/project/go-tdd-server/internal/service"
)

type MCPHandler struct {
	AppSvc service.AppService
	Server *mcp.Server
	SSE    *mcp.SSEHandler
}

func NewMCPHandler(appSvc service.AppService) *MCPHandler {
	// Initialize MCP Server with metadata
	s := mcp.NewServer(&mcp.Implementation{
		Name:    "go-tdd-server",
		Version: "0.0.2",
	}, nil)

	h := &MCPHandler{
		AppSvc: appSvc,
		Server: s,
	}

	h.SSE = mcp.NewSSEHandler(func(req *http.Request) *mcp.Server {
		return h.Server
	}, nil)

	h.registerTools()

	return h
}

// Tool Arguments with JSONSchema tags for automatic documentation
// For github.com/google/jsonschema-go, the tag content itself is the description
type AddArgs struct {
	A int `json:"a" jsonschema:"The first integer to add"`
	B int `json:"b" jsonschema:"The second integer to add"`
}

type EchoArgs struct {
	Message string `json:"message" jsonschema:"The message to echo back"`
}

func (h *MCPHandler) registerTools() {
	// Tool: Add
	// The SDK automatically generates JSON Schema from the AddArgs struct!
	mcp.AddTool(h.Server, &mcp.Tool{
		Name:        "add",
		Description: "Adds two integers together. Use this for any mathematical addition.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args AddArgs) (*mcp.CallToolResult, any, error) {
		result := h.AppSvc.Add(ctx, args.A, args.B)

		// We return a human-readable text result for the LLM to "read"
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("The result of %d + %d is %d", args.A, args.B, result)},
			},
		}, nil, nil
	})

	// Tool: Echo
	mcp.AddTool(h.Server, &mcp.Tool{
		Name:        "echo",
		Description: "Echoes back the message. Useful for testing connectivity.",
	}, func(ctx context.Context, req *mcp.CallToolRequest, args EchoArgs) (*mcp.CallToolResult, any, error) {
		processed := h.AppSvc.Echo(ctx, args.Message)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: processed},
			},
		}, nil, nil
	})
}

func (h *MCPHandler) RegisterRoutes(e *echo.Echo) {
	// SSE endpoint handling
	e.Any("/mcp/sse", echo.WrapHandler(h.SSE))
}
