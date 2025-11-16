package handlers

import (
	"context"

	"github.com/de4et/avito-test/internal/server/api"
)

type pullRequestRepository interface {
}

type PullRequestHandler struct {
	rep pullRequestRepository
}

func NewPullRequestHandler(rep pullRequestRepository) *PullRequestHandler {
	return &PullRequestHandler{rep: rep}
}

func (h *PullRequestHandler) PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error) {
	return api.PostPullRequestCreate201JSONResponse{}, nil
}

func (h *PullRequestHandler) PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error) {
	return api.PostPullRequestMerge200JSONResponse{}, nil
}

func (h *PullRequestHandler) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	return api.PostPullRequestReassign404JSONResponse{}, nil
}
