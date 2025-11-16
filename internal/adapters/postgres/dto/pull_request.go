package dto

import "time"

type PullRequest struct {
	AssignedReviewers []string          `db:"assigned_reviewers"`
	AuthorId          string            `db:"author_id"`
	CreatedAt         *time.Time        `db:"createdAt"`
	MergedAt          *time.Time        `db:"mergedAt"`
	PullRequestId     string            `db:"pull_request_id"`
	PullRequestName   string            `db:"pull_request_name"`
	Status            PullRequestStatus `db:"status"`
}

type PullRequestStatus string
