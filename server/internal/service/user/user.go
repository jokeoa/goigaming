package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
)

type Service struct {
	userRepo ports.UserRepository
}

func NewService(userRepo ports.UserRepository) *Service {
	return &Service{userRepo: userRepo}
}

func (s *Service) GetProfile(ctx context.Context, userID uuid.UUID) (domain.UserProfile, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.UserProfile{}, fmt.Errorf("UserService.GetProfile: %w", err)
	}
	return domain.NewUserProfile(user), nil
}

func (s *Service) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("UserService.GetByID: %w", err)
	}
	return user, nil
}
