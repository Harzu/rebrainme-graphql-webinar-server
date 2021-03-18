package internal

import (
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"

	"webinar/graphql/server/internal/config"
	"webinar/graphql/server/internal/constants"
	"webinar/graphql/server/internal/graph/directives"
	"webinar/graphql/server/internal/graph/generated"
	"webinar/graphql/server/internal/graph/resolvers"
	"webinar/graphql/server/internal/repositories"
	"webinar/graphql/server/internal/services"
	"webinar/graphql/server/internal/services/middlewares"
	"webinar/graphql/server/internal/system/database/psql"
	"webinar/graphql/server/internal/system/database/redis"
)

type Service struct {
	config      *config.Config
	router      *gin.Engine
	gqlConfig   generated.Config
	psqlClient  psql.Repository
	redisClient redis.Repository
}

func (s *Service) Setup() error {
	cfg, err := config.InitConfig(constants.ServiceName)
	if err != nil {
		return fmt.Errorf("failed to init config: %w", err)
	}
	s.config = cfg

	psqlClient, err := psql.New(cfg.DSN)
	if err != nil {
		return fmt.Errorf("failed to init psql connection: %w", err)
	}
	s.psqlClient = psqlClient

	redisClient, err := redis.New(cfg.Redis)
	if err != nil {
		return fmt.Errorf("failed to init redis connection: %w", err)
	}
	s.redisClient = redisClient

	repoContainer := repositories.New(psqlClient)
	serviceContainer := services.Setup(redisClient)

	s.gqlConfig = generated.Config{Resolvers: resolvers.New(serviceContainer, repoContainer)}
	directives.Setup(&s.gqlConfig, repoContainer, serviceContainer)

	s.router = gin.Default()

	return nil
}

func (s *Service) ListenAndServe() error {
	gqlHandler := handler.NewDefaultServer(generated.NewExecutableSchema(s.gqlConfig))
	graph := s.router.Group("/graph")
	graph.Use(middlewares.EnrichContextWithSession)
	graph.POST("/query", func(ctx *gin.Context) {
		gqlHandler.ServeHTTP(ctx.Writer, ctx.Request)
	})

	if s.config.PlaygroundEnable {
		gqlPlayground := playground.Handler("GraphQL", "/graph/query")
		s.router.GET("/", func(ctx *gin.Context) {
			gqlPlayground.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", s.config.Port)
	if err := s.router.Run(fmt.Sprintf(":%s", s.config.Port)); err != nil {
		return fmt.Errorf("failed to run server: %w", err)
	}

	return nil
}

func (s *Service) Shutdown() error {
	if err := s.psqlClient.GetConnection().Close(); err != nil {
		return fmt.Errorf("failed to close psql connection: %w", err)
	}
	if err := s.redisClient.GetConnection().Close(); err != nil {
		return fmt.Errorf("failed to close redis connection: %w", err)
	}

	return nil
}
