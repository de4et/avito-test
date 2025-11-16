package domain

import "time"

type PullRequest struct {
	AssignedReviewers []string
	AuthorId          string
	CreatedAt         *time.Time
	MergedAt          *time.Time
	PullRequestId     string
	PullRequestName   string
	Status            PullRequestStatus
}

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequestStatus string
