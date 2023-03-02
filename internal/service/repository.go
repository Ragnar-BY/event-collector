package service

import (
	"context"

	"github.com/Ragnar-BY/event-collector/internal/domain"
)

// EventRepository describes repository interface
type EventRepository interface {
	SaveEvents(ctx context.Context, events []domain.Event) error
	Ping(ctx context.Context) error

	Close() error
}

// RepositoryService is service for repository
type RepositoryService struct {
	repo EventRepository
}

// NewRepositoryService creates new service for repository
func NewRepositoryService(repo EventRepository) *RepositoryService {
	s := &RepositoryService{
		repo: repo,
	}

	return s
}

// SaveEvents saves events in repository
func (s *RepositoryService) SaveEvents(ctx context.Context, events []domain.Event) error {
	return s.repo.SaveEvents(ctx, events)
}

// Close closes repository
func (s *RepositoryService) CloseDB() error {
	return s.repo.Close()
}

// Ping pings repository
func (s *RepositoryService) PingDB(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
