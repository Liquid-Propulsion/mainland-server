package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/Liquid-Propulsion/mainland-server/canbackend"
	"github.com/Liquid-Propulsion/mainland-server/config"
	"github.com/Liquid-Propulsion/mainland-server/database/sql"
	"github.com/Liquid-Propulsion/mainland-server/database/timeseries"
	"github.com/Liquid-Propulsion/mainland-server/graph"
	"github.com/Liquid-Propulsion/mainland-server/graph/generated"
	"github.com/Liquid-Propulsion/mainland-server/systems"
)

const defaultPort = "8080"

func main() {
	config.Init()
	canbackend.Init(config.CurrentConfig.CAN.CANType)
	sql.Init(config.CurrentConfig.SQLite.DSN)
	timeseries.Init(false, config.CurrentConfig.TimeSeries.Directory)
	systems.Init()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("Listening on %s for GraphQL Connections...", config.CurrentConfig.HTTP.ListenAddr)
	log.Fatal(http.ListenAndServe(config.CurrentConfig.HTTP.ListenAddr, nil))
}
