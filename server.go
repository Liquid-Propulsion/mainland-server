package main

import (
	"context"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/graph"
	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/session"
	"github.com/Liquid-Propulsion/mainland-server/systems"
	"github.com/go-chi/chi"
)

func main() {
	config.Init()
	canbackend.Init(config.CurrentConfig.CAN.CANType)
	sql.Init(config.CurrentConfig.SQLite.DSN)
	timeseries.Init(false, config.CurrentConfig.TimeSeries.Directory)
	systems.Init()

	configuration := generated.Config{Resolvers: &graph.Resolver{}}

	configuration.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		if _, ok := session.ForContext(ctx); !ok {
			//return nil, fmt.Errorf("access denied")
		}

		// or let it pass through
		return next(ctx)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(configuration))

	router := chi.NewRouter()

	router.Use(session.Middleware())

	router.Handle("/", playground.Handler("LPDT", "/query"))
	router.Handle("/query", srv)

	log.Printf("Listening on %s for GraphQL Connections...", config.CurrentConfig.HTTP.ListenAddr)
	log.Fatal(http.ListenAndServe(config.CurrentConfig.HTTP.ListenAddr, router))
}
