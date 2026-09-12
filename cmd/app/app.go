package app

import (
	"github.com/dprio/redis-limiter/internal/infrastructure/cache"
	"github.com/dprio/redis-limiter/internal/infrastructure/config"
	"github.com/dprio/redis-limiter/internal/infrastructure/db"
	eventhandlers "github.com/dprio/redis-limiter/internal/infrastructure/event/handlers"
	"github.com/dprio/redis-limiter/internal/infrastructure/graph/graphqlserver"
	"github.com/dprio/redis-limiter/internal/infrastructure/graph/resolvers"
	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/grpcserver"
	"github.com/dprio/redis-limiter/internal/infrastructure/grpc/service"
	"github.com/dprio/redis-limiter/internal/infrastructure/ratelimiter"
	"github.com/dprio/redis-limiter/internal/infrastructure/web/handlers"
	"github.com/dprio/redis-limiter/internal/infrastructure/web/middlewares"
	"github.com/dprio/redis-limiter/internal/infrastructure/web/webserver"
	"github.com/dprio/redis-limiter/internal/usecase"
	"github.com/dprio/redis-limiter/pkg/events"
)

type App struct {
	useCases      *usecase.UseCases
	webServer     *webserver.WebServer
	grpcServer    *grpcserver.GRPCServer
	graphQlServer *graphqlserver.GraphQLServer
}

func New() *App {
	conf := config.New()

	dataBase := db.New(conf.DB)

	eventDispatcher := events.NewEventDispatcher()
	eventhandlers.CreateAndRegisterEventHandlers(eventDispatcher)

	redisClient := cache.NewRedisClient(conf.Redis)

	limiter := ratelimiter.NewLimiterGateway(redisClient, conf.RateLimiter)

	useCases := usecase.New(dataBase, eventDispatcher)

	handlers := handlers.New(*useCases)

	middlewares := middlewares.New(limiter)

	webServer := createWebServer(conf, handlers, middlewares)

	grpcServices := service.NewGRPCServices(useCases)

	grpcServer := grpcserver.New(grpcServices)

	graphQLResolvers := resolvers.NewGraphQLResolvers(useCases)

	graphqlServer := graphqlserver.New("8081", graphQLResolvers)

	return &App{
		useCases:      useCases,
		webServer:     webServer,
		grpcServer:    grpcServer,
		graphQlServer: graphqlServer,
	}

}

func createWebServer(conf *config.Config, handls *handlers.Handlers, middlewares *middlewares.Middlewares) *webserver.WebServer {
	webServer := webserver.New(conf.Web)

	webServer.AddHandler("POST", "/orders", middlewares.RateLimitMiddleware.Handle(handls.CreateOrderHandler.Create))
	webServer.AddHandler("GET", "/orders", middlewares.RateLimitMiddleware.Handle(handls.CreateOrderHandler.GetAll))

	return webServer
}

func (app *App) Start() error {
	go app.webServer.Start()
	go app.grpcServer.Start()

	println("app started successfully...")
	return app.graphQlServer.Start()
}
