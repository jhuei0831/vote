package config

import (
	"context"
	"fmt"
	"strconv"
	"vote/app/database"
	"vote/app/service"
	graph "vote/graph/generated"
	resolver "vote/graph/resolver"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/vektah/gqlparser/v2/ast"
)

func graphqlHandler() gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := graph.Config{Resolvers: &resolver.Resolver{}}
	graphqlService := service.NewGraphqlService()

	c.Directives.HasPermission = func(ctx context.Context, obj any, next graphql.Resolver, resource string, action string) (interface{}, error) {
		userInfo, err := graphqlService.GetUserInfoFromContext(ctx)
		if err != nil {
			return nil, err
		}
		ok, err := database.Enforcer.Enforce(strconv.FormatUint(userInfo.UserID, 10), resource, action)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("forbidden")
		}

		// or let it pass through
		return next(ctx)
	}
	srv := handler.New(graph.NewExecutableSchema(c))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	srv := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		srv.ServeHTTP(c.Writer, c.Request)
	}
}
