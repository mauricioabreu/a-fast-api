package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/db"
	"github.com/mauricioabreu/a-fast-api/people"
	"github.com/mauricioabreu/a-fast-api/validators"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	app := fiber.New()

	conn := fmt.Sprintf("host=db user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
	dbc, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("error connecting to the database: ", err)
		os.Exit(1)
	}

	queries := db.New(dbc)

	app.Get("/count-people", func(c *fiber.Ctx) error {
		ctx := context.Background()
		total, err := people.CountPeople(queries, ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to count people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(fmt.Sprintf("%d", total))
	})

	app.Post("/people", func(c *fiber.Ctx) error {
		p := new(people.PersonDTO)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		validate := validator.New()

		if errs := validators.Validate(validate, p); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
		}

		ctx := context.Background()
		if err := people.InsertPerson(*p, queries, ctx); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendString("")
	})

	app.Listen(":80")
}
