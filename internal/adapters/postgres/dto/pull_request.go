package dto

import "time"

type PullRequest struct {
	AssignedReviewers []string
	AuthorId          string            `db:"author_id"`
	MergedAt          *time.Time        `db:"merged_at"`
	CreatedAt         *time.Time        `db:"created_at"`
	PullRequestId     string            `db:"pull_request_id"`
	PullRequestName   string            `db:"pull_request_name"`
	Status            PullRequestStatus `db:"status"`
}

type PullRequestStatus string

type PullRequestShort struct {
	AuthorId        string            `db:"author_id"`
	PullRequestId   string            `db:"pull_request_id"`
	PullRequestName string            `db:"pull_request_name"`
	Status          PullRequestStatus `db:"status"`
}
