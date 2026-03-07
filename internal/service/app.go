package service

import "context"

type AppService interface {
	Add(ctx context.Context, a, b int) int
	Echo(ctx context.Context, message string) string
	GetInternalInfo(ctx context.Context, userID string) string
}

type service struct{}

func NewAppService() AppService {
	return &service{}
}

func (s *service) Add(_ context.Context, a, b int) int {
	return a + b
}

func (s *service) Echo(_ context.Context, message string) string {
	return message
}

func (s *service) GetInternalInfo(_ context.Context, userID string) string {
	// This is an internal-only method, not exposed via MCP
	return "Internal data for user: " + userID
}
