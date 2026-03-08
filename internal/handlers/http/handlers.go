package httphandlers

import (
	"github.com/project/go-tdd-server/internal/models"
	"github.com/project/go-tdd-server/internal/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	StartTime time.Time
	AppSvc    service.AppService
	Version   string
}

func NewHandler(appSvc service.AppService, version string) *Handler {
	return &Handler{
		AppSvc:    appSvc,
		Version:   version,
		StartTime: time.Now(),
	}
}

// Health checks
func (h *Handler) Healthz(c echo.Context) error {
	uptime := time.Since(h.StartTime).Truncate(time.Second).String()
	return c.JSON(http.StatusOK, models.HealthResponse{
		Status:    "ok",
		Version:   h.Version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Uptime:    uptime,
	})
}

func (h *Handler) Readyz(c echo.Context) error {
	// Add health checks for external dependencies (e.g., DB connection) here
	return c.JSON(http.StatusOK, models.ReadyResponse{
		Status:  "ready",
		Version: h.Version,
	})
}

// API methods
func (h *Handler) Hello(c echo.Context) error {
	name := c.Param("name")
	return c.JSON(http.StatusOK, models.HelloResponse{Message: "Hello, " + name})
}

func (h *Handler) Add(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.AddRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_FORMAT",
			Message: "invalid query parameters format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "VALIDATION_FAILED",
			Message: err.Error(),
		})
	}

	result := h.AppSvc.Add(ctx, req.A, req.B)
	return c.JSON(http.StatusOK, models.AddResponse{
		Result: result,
	})
}

func (h *Handler) Echo(c echo.Context) error {
	ctx := c.Request().Context()
	req := new(models.EchoRequest)

	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "INVALID_FORMAT",
			Message: "invalid JSON format",
		})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "VALIDATION_FAILED",
			Message: err.Error(),
		})
	}

	processed := h.AppSvc.Echo(ctx, req.Message)
	return c.JSON(http.StatusOK, models.EchoResponse{
		Message:   processed,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Internal(c echo.Context) error {
	ctx := c.Request().Context()
	return c.String(http.StatusOK, h.AppSvc.GetInternalInfo(ctx, "internal-system"))
}
