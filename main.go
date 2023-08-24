package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/people"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := fiber.New()

	conn := fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("error connecting to the database: ", err)
		os.Exit(1)
	}

	queries := people.New(db)

	app.Get("/contagem-pessoas", func(c *fiber.Ctx) error {
		ctx := context.Background()
		total, err := queries.CountPeople(ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to count people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(fmt.Sprintf("%d", total))
	})

	app.Listen(":80")
}
