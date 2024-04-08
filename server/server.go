package server

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"larrymyers.com/rk9api/api"
	"larrymyers.com/rk9api/graph"
	"larrymyers.com/rk9api/graph/generated"
)

func NewRouter(client *api.Client) http.Handler {
	graphqlServer := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Client: client}}))

	mux := http.NewServeMux()
	mux.Handle("/graphql", playground.Handler("RK9 GraphQL API", "/query"))
	mux.Handle("/query", graphqlServer)

	return mux
}
