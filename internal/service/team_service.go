package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/de4et/avito-test/internal/domain"
)

var (
	ErrTeamExists    = errors.New("team already exists")
	ErrTeamNotExists = errors.New("team not exists")
)

type teamRepository interface {
	Add(context.Context, domain.Team) error
	Get(context.Context, string) (domain.Team, error)
	GetByUserID(context.Context, string) (domain.Team, error)
}

type TeamService struct {
	rep teamRepository
	tr  transactor
}

func NewTeamService(rep teamRepository, tr transactor) *TeamService {
	return &TeamService{
		rep: rep,
		tr:  tr,
	}
}

func (svc *TeamService) AddTeam(ctx context.Context, name string, members []domain.TeamMember) (domain.Team, error) {
	slog.DebugContext(ctx, "Creating team in service")

	team := domain.Team{
		Members:  members,
		TeamName: name,
	}

	return team, svc.tr.WithTx(ctx, func(ctx context.Context) error {
		return svc.rep.Add(ctx, team)
	})
}

func (svc *TeamService) GetTeam(ctx context.Context, name string) (domain.Team, error) {
	slog.DebugContext(ctx, "Getting team in service")
	return svc.rep.Get(ctx, name)
}
