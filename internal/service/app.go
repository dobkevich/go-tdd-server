package service

import "context"

type AppService interface {
	Add(ctx context.Context, a, b int) int
	Echo(ctx context.Context, message string) string
}

type appService struct{}

func NewAppService() AppService {
	return &appService{}
}

func (s *appService) Add(ctx context.Context, a, b int) int {
	return a + b
}

func (s *appService) Echo(ctx context.Context, message string) string {
	return message
}
