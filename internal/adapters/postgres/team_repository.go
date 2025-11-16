package postgres

import (
	"context"
	_ "embed"

	"github.com/de4et/avito-test/internal/adapters/postgres/dto"
	"github.com/de4et/avito-test/internal/domain"
	"github.com/de4et/avito-test/internal/service"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

//go:embed queries/team_create.sql
var teamCreateQuery string

//go:embed queries/team_get.sql
var teamGetQuery string

//go:embed queries/user_add.sql
var userAddQuery string

//go:embed queries/user_get_team_name.sql
var userGetTeamQuery string

type PostgresqlTeamRepository struct {
	client *TxClient
}

func NewPostgresqlTeamRepository(client *sqlx.DB) *PostgresqlTeamRepository {
	return &PostgresqlTeamRepository{
		client: NewTxClient(client),
	}
}

func (rep *PostgresqlTeamRepository) Add(ctx context.Context, team domain.Team) error {
	_, err := rep.client.ExecContext(ctx, teamCreateQuery, team.TeamName)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && IsDuplicate(pgErr) {
			return service.ErrTeamExists
		}
		return err
	}

	for _, m := range team.Members {
		_, err := rep.client.ExecContext(ctx,
			userAddQuery,
			m.UserId,
			m.Username,
			team.TeamName,
			m.IsActive,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rep *PostgresqlTeamRepository) Get(ctx context.Context, name string) (domain.Team, error) {
	var team dto.Team
	team.TeamName = name

	err := rep.client.SelectContext(
		ctx,
		&team.Members,
		teamGetQuery,
		name,
	)

	if len(team.Members) == 0 {
		return domain.Team{}, service.ErrTeamNotExists
	}

	if err != nil {
		return domain.Team{}, err
	}

	members := make([]domain.TeamMember, len(team.Members))
	for i := range team.Members {
		members[i] = domain.TeamMember{
			IsActive: team.Members[i].IsActive,
			UserId:   team.Members[i].UserId,
			Username: team.Members[i].Username,
		}
	}

	return domain.Team{
		Members:  members,
		TeamName: team.TeamName,
	}, nil
}

func (rep *PostgresqlTeamRepository) GetByUserID(ctx context.Context, id string) (domain.Team, error) {
	var name string
	err := rep.client.GetContext(ctx, &name, userGetTeamQuery, id)

	if err != nil {
		return domain.Team{}, err
	}

	var team dto.Team
	team.TeamName = name

	err = rep.client.SelectContext(
		ctx,
		&team.Members,
		teamGetQuery,
		name,
	)

	if len(team.Members) == 0 {
		return domain.Team{}, service.ErrTeamNotExists
	}

	if err != nil {
		return domain.Team{}, err
	}

	members := make([]domain.TeamMember, len(team.Members))
	for i := range team.Members {
		members[i] = domain.TeamMember{
			IsActive: team.Members[i].IsActive,
			UserId:   team.Members[i].UserId,
			Username: team.Members[i].Username,
		}
	}

	return domain.Team{
		Members:  members,
		TeamName: team.TeamName,
	}, nil
}
