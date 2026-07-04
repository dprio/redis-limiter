package graphqlserver

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/dprio/redis-limiter/internal/infrastructure/graph"
	"github.com/dprio/redis-limiter/internal/infrastructure/graph/resolvers"
	"github.com/vektah/gqlparser/v2/ast"
)

type GraphQLServer struct {
	port   string
	server *handler.Server
}

func New(port string, resolvers *resolvers.Resolvers) *GraphQLServer {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers.OrderResolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return &GraphQLServer{
		port:   port,
		server: srv,
	}
}

func (g *GraphQLServer) Start() error {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", g.server)

	return http.ListenAndServe(":"+g.port, nil)
}
