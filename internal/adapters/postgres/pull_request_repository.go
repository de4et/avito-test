package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"

	"github.com/de4et/avito-test/internal/domain"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/pull_request_exists.sql
var prExistsQuery string

//go:embed queries/pull_request_create.sql
var prCreateQuery string

//go:embed queries/pull_request_insert_reviewer.sql
var prInsertReviewerQuery string

type PostgresqlPullRequestRepository struct {
	client *TxClient
}

func NewPostgresqlPullRequestRepository(client *sqlx.DB) *PostgresqlPullRequestRepository {
	return &PostgresqlPullRequestRepository{
		client: NewTxClient(client),
	}
}

func (rep *PostgresqlPullRequestRepository) IsExists(ctx context.Context, pullRequestID string) (bool, error) {
	var exists int

	err := rep.client.GetContext(ctx, &exists, prExistsQuery, pullRequestID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (rep *PostgresqlPullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) error {
	_, err := rep.client.ExecContext(
		ctx,
		prCreateQuery,
		pr.PullRequestId,
		pr.PullRequestName,
		pr.AuthorId,
		pr.Status,
	)
	if err != nil {
		return err
	}

	for _, reviewer := range pr.AssignedReviewers {
		_, err := rep.client.ExecContext(ctx,
			prInsertReviewerQuery,
			pr.PullRequestId,
			reviewer,
		)
		if err != nil {
			return err
		}
	}

	return nil
}
