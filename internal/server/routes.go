package server

import (
	_ "embed"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	middleware "github.com/oapi-codegen/gin-middleware"

	"github.com/de4et/avito-test/internal/server/api"
	"github.com/de4et/avito-test/internal/server/handlers"
	"github.com/de4et/avito-test/internal/server/middlewares"
	"github.com/de4et/avito-test/internal/service"
	"github.com/gin-gonic/gin"
)

//go:embed api/openapi.yml
var openapiSchema []byte

type server struct {
	*handlers.UserHandler
	*handlers.TeamHandler
	*handlers.PullRequestHandler
}

func RegisterRoutes(teamService *service.TeamService, userService *service.UserService, prService *service.PullRequestService) http.Handler {
	swagger, err := openapi3.NewLoader().LoadFromData(openapiSchema)
	if err != nil {
		panic("couldn't load openapiSchema")
	}

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middlewares.LogHandler())
	r.Use(middleware.OapiRequestValidator(swagger))
	r.Use(middlewares.ErrorHandler())

	api.RegisterHandlers(r, api.NewStrictHandler(
		server{
			UserHandler:        handlers.NewUserHandler(userService),
			TeamHandler:        handlers.NewTeamHandler(teamService),
			PullRequestHandler: handlers.NewPullRequestHandler(prService),
		},
		[]api.StrictMiddlewareFunc{},
	))

	r.GET("/health", healthHandler)

	return r
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
}
