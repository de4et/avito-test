package handlers

import (
	"context"
	"errors"
	"log/slog"

	"github.com/de4et/avito-test/internal/domain"
	"github.com/de4et/avito-test/internal/server/api"
	"github.com/de4et/avito-test/internal/service"
	logger "github.com/de4et/avito-test/pkg"
)

type PullRequestHandler struct {
	svc *service.PullRequestService
}

func NewPullRequestHandler(svc *service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{
		svc: svc,
	}
}

func (h *PullRequestHandler) PostPullRequestCreate(ctx context.Context, request api.PostPullRequestCreateRequestObject) (api.PostPullRequestCreateResponseObject, error) {
	ctx = logger.WithContext(ctx, "author_id", request.Body.AuthorId)
	ctx = logger.WithContext(ctx, "pull_request_id", request.Body.PullRequestId)
	ctx = logger.WithContext(ctx, "pull_request_name", request.Body.PullRequestName)
	slog.InfoContext(ctx, "Creating pull request")

	pr, err := h.svc.CreatePullRequest(
		ctx,
		request.Body.AuthorId,
		request.Body.PullRequestId,
		request.Body.PullRequestName,
	)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			return api.PostPullRequestCreate404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrUserNotFoundMsg,
			)), nil
		}

		if errors.Is(err, service.ErrTeamNotExists) {
			return api.PostPullRequestCreate404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrTeamNotExistsMsg,
			)), nil
		}

		if errors.Is(err, service.ErrPRExists) {
			return api.PostPullRequestCreate409JSONResponse(api.NewError(
				api.PREXISTS,
				api.ErrPRExistsMsg,
			)), nil
		}

		return nil, err
	}

	apiPullRequest := fromDomainPullRequest(pr)
	return api.PostPullRequestCreate201JSONResponse{
		Pr: &apiPullRequest,
	}, nil
}

func (h *PullRequestHandler) PostPullRequestMerge(ctx context.Context, request api.PostPullRequestMergeRequestObject) (api.PostPullRequestMergeResponseObject, error) {
	panic("NOT IMPLEMENTED")

}

func (h *PullRequestHandler) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	panic("NOT IMPLEMENTED")

}

func fromDomainPullRequest(pr domain.PullRequest) api.PullRequest {
	return api.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		Status:            api.PullRequestStatus(pr.Status),
	}
}
