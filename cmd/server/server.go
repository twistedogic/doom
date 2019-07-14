package server

import (
	"log"
	"net/http"

	"github.com/twistedogic/doom/gql/generated"
	"github.com/vektah/gqlgen/handler"
)

func Start(port int) {
	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &generated.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	return http.ListenAndServe(":"+port, nil)

}
