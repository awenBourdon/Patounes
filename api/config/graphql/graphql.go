package graphql

import (
	"net/http"

	"patoune-api/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func NewGraphQLHandler() http.Handler {
	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{},
	}))

	return server
}

func NewPlaygroundHandler() http.Handler {
	return playground.ApolloSandboxHandler("GraphQL Sandbox", "/graphql")
}
