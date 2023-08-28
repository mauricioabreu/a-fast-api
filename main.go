package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/api"
	"github.com/mauricioabreu/a-fast-api/db"
	"github.com/mauricioabreu/a-fast-api/tools"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Use(logger.New())

	dbc, err := db.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	queries := db.New(dbc)
	service := api.NewPeopleService(queries)

	api.PeopleRouter(app, service)

	defer tools.StartProfiling(os.Getenv("PROFILER_MODE")).Stop()

	if err := app.Listen(":80"); err != nil {
		log.Warn().Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
