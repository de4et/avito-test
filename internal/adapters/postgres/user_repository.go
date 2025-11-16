package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/de4et/avito-test/internal/adapters/postgres/dto"
	"github.com/de4et/avito-test/internal/domain"
	"github.com/de4et/avito-test/internal/service"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/user_set_active.sql
var userSetActiveQuery string

type PostgresqlUserRepository struct {
	client *sqlx.DB
}

func NewPostgresqlUserRepository(client *sqlx.DB) *PostgresqlUserRepository {
	return &PostgresqlUserRepository{
		client: client,
	}
}

func (rep *PostgresqlUserRepository) SetActive(ctx context.Context, userID string, active bool) (domain.User, error) {
	var user dto.User

	err := rep.client.GetContext(
		ctx,
		&user,
		userSetActiveQuery,
		active,
		userID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, service.ErrUserNotExists
		}
		return domain.User{}, err
	}

	return domain.User{
		IsActive: user.IsActive,
		TeamName: user.TeamName,
		UserId:   user.UserId,
		Username: user.Username,
	}, nil
}
