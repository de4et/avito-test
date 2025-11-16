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
}

type TeamService struct {
	rep teamRepository
}

func NewTeamService(rep teamRepository) *TeamService {
	return &TeamService{
		rep: rep,
	}
}

func (svc *TeamService) AddTeam(ctx context.Context, name string, members []domain.TeamMember) (domain.Team, error) {
	slog.DebugContext(ctx, "Creating team in service")

	team := domain.Team{
		Members:  members,
		TeamName: name,
	}
	return team, svc.rep.Add(ctx, team)
}

func (svc *TeamService) GetTeam(ctx context.Context, name string) (domain.Team, error) {
	slog.DebugContext(ctx, "Getting team in service")
	return svc.rep.Get(ctx, name)
}
