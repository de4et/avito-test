package service

import (
	"context"
	"errors"
	"math/rand"
	"slices"

	"github.com/de4et/avito-test/internal/domain"
)

var (
	ErrPRExists    = errors.New("pull request already exists")
	ErrPRNotExists = errors.New("pull request not exists")
	ErrNotAssigned = errors.New("reviewer is not assigned to this PR")
	ErrNoCandidate = errors.New("no active replacement candidate in team")
	ErrMerged      = errors.New("cannot reassign on merged PR ")
)

type pullRequestRepository interface {
	IsExists(context.Context, string) (bool, error)
	Create(context.Context, domain.PullRequest) (domain.PullRequest, error)
	Merge(context.Context, string) (domain.PullRequest, error)
	Get(context.Context, string) (domain.PullRequest, error)
	UpdateReviewer(context.Context, string, string, string) (domain.PullRequest, error)
	GetAllByUserID(context.Context, string) ([]domain.PullRequest, error)
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
			if team.Members[i].IsActive && team.Members[i].UserId != authorID {
				activeMembers = append(activeMembers, team.Members[i])
			}
		}

		var chosenLen int
		if len(activeMembers) >= 2 {
			chosenLen = 2
		} else if len(activeMembers) == 1 {
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
		pr, err = svc.rep.Create(ctx, pr)
		return err
	})
}

func (svc *PullRequestService) MergePullRequest(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	var prRes domain.PullRequest
	return prRes, svc.tx.WithTx(ctx, func(ctx context.Context) error {
		pr, err := svc.rep.Merge(ctx, pullRequestID)
		prRes = pr
		return err
	})
}

func (svc *PullRequestService) ReassignPullRequest(ctx context.Context, pullRequestID, oldUserID string) (domain.PullRequest, string, error) {
	var prRes domain.PullRequest
	var replacedBy string
	return prRes, replacedBy, svc.tx.WithTx(ctx, func(ctx context.Context) error {
		ok, err := svc.userRep.IsExists(ctx, oldUserID)
		if err != nil {
			return err
		}
		if !ok {
			return ErrUserNotExists
		}

		pr, err := svc.rep.Get(ctx, pullRequestID)
		if err != nil {
			return err
		}
		if pr.Status == domain.PullRequestStatusMERGED {
			return ErrMerged
		}
		if slices.Index(pr.AssignedReviewers, oldUserID) == -1 {
			return ErrNotAssigned
		}

		team, err := svc.teamRep.GetByUserID(ctx, oldUserID)
		if err != nil {
			return err
		}

		freeMembers := make([]domain.TeamMember, 0, len(team.Members))
		for i := range team.Members {
			if team.Members[i].IsActive && team.Members[i].UserId != pr.AuthorId && team.Members[i].UserId != oldUserID && slices.Index(pr.AssignedReviewers, team.Members[i].UserId) == -1 {
				freeMembers = append(freeMembers, team.Members[i])
			}
		}

		if len(freeMembers) == 0 {
			return ErrNoCandidate
		}

		replacedBy = freeMembers[rand.Intn(len(freeMembers))].UserId

		pr, err = svc.rep.UpdateReviewer(ctx, pullRequestID, oldUserID, replacedBy)
		if err != nil {
			return err
		}

		prRes = pr
		return err
	})
}
