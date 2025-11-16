package service

import (
	"context"
	"errors"

	"github.com/de4et/avito-test/internal/domain"
)

var (
	ErrUserNotExists = errors.New("user not exists")
)

type userRepository interface {
	SetActive(context.Context, string, bool) (domain.User, error)
}

type UserService struct {
	rep userRepository
}

func NewUserService(rep userRepository) *UserService {
	return &UserService{
		rep: rep,
	}
}

func (svc *UserService) SetActive(ctx context.Context, userID string, active bool) (domain.User, error) {
	return svc.rep.SetActive(ctx, userID, active)
}
