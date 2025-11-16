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

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		svc: svc,
	}
}

// TODO:
func (h *UserHandler) GetUsersGetReview(ctx context.Context, request api.GetUsersGetReviewRequestObject) (api.GetUsersGetReviewResponseObject, error) {
	ctx = logger.WithContext(ctx, "user_id", request.Params.UserId)
	slog.InfoContext(ctx, "Getting review")

	prs, err := h.svc.GetReview(ctx, request.Params.UserId)
	if err != nil {
		// no error in specs (only 200 status)
		return nil, err
	}

	apiPRs := toPullRequestShorts(prs)
	return api.GetUsersGetReview200JSONResponse{
		PullRequests: apiPRs,
		UserId:       request.Params.UserId,
	}, nil
}

func (h *UserHandler) PostUsersSetIsActive(ctx context.Context, request api.PostUsersSetIsActiveRequestObject) (api.PostUsersSetIsActiveResponseObject, error) {
	ctx = logger.WithContext(ctx, "user_id", request.Body.UserId)
	ctx = logger.WithContext(ctx, "is_active", request.Body.IsActive)
	slog.InfoContext(ctx, "Setting active")

	user, err := h.svc.SetActive(ctx, request.Body.UserId, request.Body.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrUserNotExists) {
			return api.PostUsersSetIsActive404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrUserNotFoundMsg,
			)), nil
		}
		return nil, err
	}

	apiUser := fromDomainUser(user)
	return api.PostUsersSetIsActive200JSONResponse{
		User: &apiUser,
	}, nil
}

func fromDomainUser(user domain.User) api.User {
	return api.User{
		IsActive: user.IsActive,
		TeamName: user.TeamName,
		UserId:   user.UserId,
		Username: user.Username,
	}
}

func toPullRequestShorts(prs []domain.PullRequest) []api.PullRequestShort {
	a := make([]api.PullRequestShort, len(prs))
	for i := range prs {
		a[i] = api.PullRequestShort{
			AuthorId:        prs[i].AuthorId,
			PullRequestId:   prs[i].PullRequestId,
			PullRequestName: prs[i].PullRequestName,
			Status:          api.PullRequestShortStatus(prs[i].Status),
		}
	}
	return a
}
