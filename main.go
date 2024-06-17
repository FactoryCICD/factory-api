package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/FactoryCICD/factory-api/cmd/server"
	"github.com/FactoryCICD/factory-api/pkg/config"
	dbcontext "github.com/FactoryCICD/factory-api/pkg/db"
	"github.com/FactoryCICD/factory-api/pkg/log"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Create logger
	logger := log.New().With(context.Background(), "version", "1.0.0")

	// Create connection to db
	conn, err := pgx.Connect(context.Background(), config.EnvVars.GetDBURI())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	// Close the connection
	defer conn.Close(context.Background())

	logger.Info("Connection to database established")

	// Create a new dbcontext
	db := dbcontext.New(conn)

	// Start the server
	app := server.Build(logger, db)

	logger.Info("Server is running on port ", config.EnvVars.GetPort())
	http.ListenAndServe(fmt.Sprintf("localhost:%s", config.EnvVars.GetPort()), app)
}
