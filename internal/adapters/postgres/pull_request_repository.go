package postgres

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/de4et/avito-test/internal/adapters/postgres/dto"
	"github.com/de4et/avito-test/internal/domain"
	"github.com/de4et/avito-test/internal/service"
	"github.com/jmoiron/sqlx"
)

//go:embed queries/pull_request_exists.sql
var prExistsQuery string

//go:embed queries/pull_request_create.sql
var prCreateQuery string

//go:embed queries/pull_request_insert_reviewer.sql
var prInsertReviewerQuery string

//go:embed queries/pull_request_merge.sql
var prMergeQuery string

//go:embed queries/pull_request_get.sql
var prGetQuery string

//go:embed queries/pull_request_get_reviewers.sql
var prGetReviewersQuery string

//go:embed queries/pull_request_update_reviewer.sql
var prUpdateReviewerQuery string

//go:embed queries/pull_request_get_all_by_user_id.sql
var prGetAllByUserIDQuery string

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

func (rep *PostgresqlPullRequestRepository) Create(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	var prDTO dto.PullRequest
	err := rep.client.GetContext(
		ctx,
		&prDTO,
		prCreateQuery,
		pr.PullRequestId,
		pr.PullRequestName,
		pr.AuthorId,
		pr.Status,
	)
	fmt.Println(prDTO)
	if err != nil {
		return domain.PullRequest{}, err
	}

	for _, reviewer := range pr.AssignedReviewers {
		_, err := rep.client.ExecContext(ctx,
			prInsertReviewerQuery,
			pr.PullRequestId,
			reviewer,
		)
		prDTO.AssignedReviewers = append(prDTO.AssignedReviewers, reviewer)
		if err != nil {
			return domain.PullRequest{}, err
		}
	}

	return domain.PullRequest{
		AssignedReviewers: prDTO.AssignedReviewers,
		AuthorId:          prDTO.AuthorId,
		CreatedAt:         prDTO.CreatedAt,
		MergedAt:          prDTO.MergedAt,
		PullRequestId:     prDTO.PullRequestId,
		PullRequestName:   prDTO.PullRequestName,
		Status:            domain.PullRequestStatus(prDTO.Status),
	}, nil
}

func (rep *PostgresqlPullRequestRepository) Merge(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	var pr dto.PullRequest

	err := rep.client.GetContext(
		ctx,
		&pr,
		prMergeQuery,
		domain.PullRequestStatusMERGED,
		pullRequestID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, service.ErrPRNotExists
		}
		return domain.PullRequest{}, err
	}

	err = rep.client.SelectContext(
		ctx,
		&pr.AssignedReviewers,
		prGetReviewersQuery,
		pullRequestID,
	)

	return domain.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		MergedAt:          pr.MergedAt,
		CreatedAt:         pr.CreatedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            domain.PullRequestStatus(pr.Status),
	}, err
}

func (rep *PostgresqlPullRequestRepository) Get(ctx context.Context, pullRequestID string) (domain.PullRequest, error) {
	var pr dto.PullRequest

	err := rep.client.GetContext(
		ctx,
		&pr,
		prGetQuery,
		pullRequestID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.PullRequest{}, service.ErrPRNotExists
		}
		return domain.PullRequest{}, err
	}

	err = rep.client.SelectContext(
		ctx,
		&pr.AssignedReviewers,
		prGetReviewersQuery,
		pullRequestID,
	)

	return domain.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		MergedAt:          pr.MergedAt,
		CreatedAt:         pr.CreatedAt,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            domain.PullRequestStatus(pr.Status),
	}, err
}

func (rep *PostgresqlPullRequestRepository) UpdateReviewer(ctx context.Context, pullRequestID, from, to string) (domain.PullRequest, error) {
	_, err := rep.client.ExecContext(
		ctx,
		prUpdateReviewerQuery,
		to,
		from,
		pullRequestID,
	)

	if err != nil {
		return domain.PullRequest{}, err
	}

	fmt.Println("updated from ", from, "to", to)
	return rep.Get(ctx, pullRequestID)
}

func (rep *PostgresqlPullRequestRepository) GetAllByUserID(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	var prs []dto.PullRequestShort

	err := rep.client.SelectContext(
		ctx,
		&prs,
		prGetAllByUserIDQuery,
		userID,
	)

	if err != nil {
		return []domain.PullRequest{}, err
	}

	arr := make([]domain.PullRequest, len(prs))
	for i := range prs {
		arr[i] = domain.PullRequest{
			AuthorId:        prs[i].AuthorId,
			PullRequestId:   prs[i].PullRequestId,
			PullRequestName: prs[i].PullRequestName,
			Status:          domain.PullRequestStatus(prs[i].Status),
		}
	}

	return arr, err
}
