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
	ctx = logger.WithContext(ctx, "pull_request_id", request.Body.PullRequestId)
	slog.InfoContext(ctx, "Merging pull request")

	pr, err := h.svc.MergePullRequest(ctx, request.Body.PullRequestId)
	if err != nil {
		if errors.Is(err, service.ErrPRNotExists) {
			return api.PostPullRequestMerge404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrPRNotExistsMsg,
			)), nil
		}
		return nil, err
	}

	apiPullRequest := fromDomainPullRequest(pr)
	return api.PostPullRequestMerge200JSONResponse{
		Pr: &apiPullRequest,
	}, nil
}

func (h *PullRequestHandler) PostPullRequestReassign(ctx context.Context, request api.PostPullRequestReassignRequestObject) (api.PostPullRequestReassignResponseObject, error) {
	ctx = logger.WithContext(ctx, "pull_request_id", request.Body.PullRequestId)
	ctx = logger.WithContext(ctx, "old_reviewer_id", request.Body.OldUserId)
	slog.InfoContext(ctx, "Reassign pull request")

	pr, replacedBy, err := h.svc.ReassignPullRequest(ctx, request.Body.PullRequestId, request.Body.OldUserId)
	if err != nil {
		if errors.Is(err, service.ErrPRNotExists) {
			return api.PostPullRequestReassign404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrPRNotExistsMsg,
			)), nil
		}

		if errors.Is(err, service.ErrUserNotExists) {
			return api.PostPullRequestReassign404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrUserNotFoundMsg,
			)), nil
		}

		if errors.Is(err, service.ErrMerged) {
			return api.PostPullRequestReassign409JSONResponse(api.NewError(
				api.PRMERGED,
				api.ErrMergedMsg,
			)), nil
		}

		if errors.Is(err, service.ErrNotAssigned) {
			return api.PostPullRequestReassign409JSONResponse(api.NewError(
				api.NOTASSIGNED,
				api.ErrNotAssignedMsg,
			)), nil
		}

		if errors.Is(err, service.ErrNoCandidate) {
			return api.PostPullRequestReassign409JSONResponse(api.NewError(
				api.NOCANDIDATE,
				api.ErrNoCandidateMsg,
			)), nil
		}
		return nil, err
	}

	apiPullRequest := fromDomainPullRequest(pr)
	return api.PostPullRequestReassign200JSONResponse{
		Pr:         apiPullRequest,
		ReplacedBy: replacedBy,
	}, nil
}

func fromDomainPullRequest(pr domain.PullRequest) api.PullRequest {
	return api.PullRequest{
		AssignedReviewers: pr.AssignedReviewers,
		AuthorId:          pr.AuthorId,
		PullRequestId:     pr.PullRequestId,
		PullRequestName:   pr.PullRequestName,
		MergedAt:          pr.MergedAt,
		CreatedAt:         pr.CreatedAt,
		Status:            api.PullRequestStatus(pr.Status),
	}
}
