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

//go:embed queries/user_exists.sql
var userExistsQuery string

type PostgresqlUserRepository struct {
	client *TxClient
}

func NewPostgresqlUserRepository(client *sqlx.DB) *PostgresqlUserRepository {
	return &PostgresqlUserRepository{
		client: NewTxClient(client),
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

func (rep *PostgresqlUserRepository) IsExists(ctx context.Context, id string) (bool, error) {
	var exists int

	err := rep.client.GetContext(ctx, &exists, userExistsQuery, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
