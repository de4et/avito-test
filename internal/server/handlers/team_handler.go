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

type TeamHandler struct {
	svc *service.TeamService
}

func NewTeamHandler(svc *service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

func (h *TeamHandler) PostTeamAdd(
	ctx context.Context,
	request api.PostTeamAddRequestObject,
) (api.PostTeamAddResponseObject, error) {
	ctx = logger.WithContext(ctx, "name", request.Body.TeamName)
	slog.InfoContext(ctx, "Creating team")

	team, err := h.svc.AddTeam(ctx, request.Body.TeamName, toDomainMembers(request.Body.Members))
	if err != nil {
		if errors.Is(err, service.ErrTeamExists) {
			return api.PostTeamAdd400JSONResponse(api.NewError(
				api.TEAMEXISTS,
				api.ErrTeamExistsMsg,
			)), nil
		}
		return nil, err
	}

	slog.InfoContext(ctx, "Successfully created team")
	return api.PostTeamAdd201JSONResponse{
		Team: &api.Team{
			TeamName: team.TeamName,
			Members:  fromDomainMembers(team.Members),
		},
	}, nil
}

func (h *TeamHandler) GetTeamGet(ctx context.Context, request api.GetTeamGetRequestObject) (api.GetTeamGetResponseObject, error) {
	ctx = logger.WithContext(ctx, "name", request.Params.TeamName)
	slog.InfoContext(ctx, "Getting team")

	team, err := h.svc.GetTeam(ctx, request.Params.TeamName)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotExists) {
			return api.GetTeamGet404JSONResponse(api.NewError(
				api.NOTFOUND,
				api.ErrTeamNotExistsMsg,
			)), nil
		}
		return nil, err
	}
	return api.GetTeamGet200JSONResponse{
		TeamName: team.TeamName,
		Members:  fromDomainMembers(team.Members),
	}, nil
}

func toDomainMembers(members []api.TeamMember) []domain.TeamMember {
	m := make([]domain.TeamMember, len(members))
	for i := range members {
		member := domain.TeamMember{
			IsActive: members[i].IsActive,
			UserId:   members[i].UserId,
			Username: members[i].Username,
		}
		m[i] = member
	}
	return m
}

func fromDomainMembers(members []domain.TeamMember) []api.TeamMember {
	m := make([]api.TeamMember, len(members))
	for i := range members {
		member := api.TeamMember{
			IsActive: members[i].IsActive,
			UserId:   members[i].UserId,
			Username: members[i].Username,
		}
		m[i] = member
	}
	return m
}
