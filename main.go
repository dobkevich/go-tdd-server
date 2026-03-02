package main

import (
	"github.com/project/go-tdd-server/internal/handlers"
	"github.com/project/go-tdd-server/internal/service"
	"context"
	"embed"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "go.uber.org/automaxprocs"
)

//go:embed docs/*
var docsContents embed.FS

const defaultPort = "8080"
const appVersion = "1.1.0"

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	e := SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      e,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info("Starting server", "port", port, "version", appVersion)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}

func SetupRouter() *echo.Echo {
	e := echo.New()

	// Setup Validator
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	appSvc := service.NewAppService()
	h := handlers.NewHandler(appSvc, appVersion)

	e.GET("/healthz", h.Healthz)
	e.GET("/readyz", h.Readyz)

	api := e.Group("/api/v1")
	api.GET("/ping", h.Ping)
	api.GET("/hello/:name", h.Hello)
	api.GET("/status", h.Status)
	api.GET("/add", h.Add)
	api.POST("/echo", h.Echo)
	api.GET("/time", h.Time)

	e.StaticFS("/docs/", echo.MustSubFS(docsContents, "docs"))
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/docs/")
	})

	return e
}
