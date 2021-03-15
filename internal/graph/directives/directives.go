package directives

import (
	"webinar/graphql/server/internal/graph/directives/auth"
	"webinar/graphql/server/internal/graph/generated"
	"webinar/graphql/server/internal/repositories"
	"webinar/graphql/server/internal/services"
)

func Setup(gqlConfig *generated.Config, repoContainer *repositories.Container, serviceContainer *services.Container) {
	authDirective := auth.New(repoContainer.Users, serviceContainer.SessionStorage)
	gqlConfig.Directives.Auth = authDirective.Resolve
}
