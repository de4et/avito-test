package service

import (
	"context"
	"errors"
	"math/rand"

	"github.com/de4et/avito-test/internal/domain"
)

var (
	ErrPRExists = errors.New("pull request already exists")
)

type pullRequestRepository interface {
	IsExists(context.Context, string) (bool, error)
	Create(context.Context, domain.PullRequest) error
}

type PullRequestService struct {
	rep     pullRequestRepository
	userRep userRepository
	teamRep teamRepository
	tx      transactor
}

func NewPullRequestService(rep pullRequestRepository, userRepository userRepository, teamRepository teamRepository, tx transactor) *PullRequestService {
	return &PullRequestService{
		rep:     rep,
		userRep: userRepository,
		teamRep: teamRepository,
		tx:      tx,
	}
}

func (svc *PullRequestService) CreatePullRequest(ctx context.Context, authorID, pullRequestID, pullRequestName string) (domain.PullRequest, error) {
	var pr domain.PullRequest
	return pr, svc.tx.WithTx(ctx, func(ctx context.Context) error {
		ok, err := svc.rep.IsExists(ctx, pullRequestID)
		if err != nil {
			return err
		}
		if ok {
			return ErrPRExists
		}

		ok, err = svc.userRep.IsExists(ctx, authorID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrUserNotExists
		}

		team, err := svc.teamRep.GetByUserID(ctx, authorID)
		if err != nil {
			return err
		}

		activeMembers := make([]domain.TeamMember, 0, len(team.Members))
		for i := range team.Members {
			if team.Members[i].IsActive {
				activeMembers = append(activeMembers, team.Members[i])
			}
		}

		var chosenLen int
		if len(activeMembers) > 2 {
			chosenLen = 2
		} else if len(activeMembers) == 2 {
			chosenLen = 1
		} else {
			chosenLen = 0
		}

		chosen := make([]string, chosenLen)
		rand.Shuffle(len(activeMembers), func(i int, j int) {
			activeMembers[i], activeMembers[j] = activeMembers[j], activeMembers[i]
		})
		for i := 0; i < chosenLen; i++ {
			chosen[i] = activeMembers[i].UserId
		}

		pr = domain.PullRequest{
			AssignedReviewers: chosen,
			AuthorId:          authorID,
			PullRequestId:     pullRequestID,
			PullRequestName:   pullRequestName,
			Status:            domain.PullRequestStatusOPEN,
		}
		return svc.rep.Create(ctx, pr)
	})
}
