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
	IsExists(context.Context, string) (bool, error)
}

type UserService struct {
	rep            userRepository
	pullRequestRep pullRequestRepository
	tx             transactor
}

func NewUserService(rep userRepository, pullRequestRep pullRequestRepository, tx transactor) *UserService {
	return &UserService{
		rep:            rep,
		tx:             tx,
		pullRequestRep: pullRequestRep,
	}
}

func (svc *UserService) SetActive(ctx context.Context, userID string, active bool) (domain.User, error) {
	return svc.rep.SetActive(ctx, userID, active)
}

func (svc *UserService) GetReview(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var prs []domain.PullRequest
	return prs, svc.tx.WithTx(ctx, func(ctx context.Context) error {
		var err error
		prs, err = svc.pullRequestRep.GetAllByUserID(ctx, userID)
		return err
	})
}
