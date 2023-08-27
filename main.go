package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	_ "github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/db"
	"github.com/mauricioabreu/a-fast-api/people"
	"github.com/mauricioabreu/a-fast-api/tools"
	"github.com/mauricioabreu/a-fast-api/validators"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	validate := validator.New()

	app := fiber.New(fiber.Config{
		Prefork: false,
	})

	app.Use(logger.New())

	dbc, err := db.NewDB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database")
	}

	queries := db.New(dbc)

	app.Get("/contagem-pessoas", func(c *fiber.Ctx) error {
		ctx := context.Background()
		total, err := people.CountPeople(queries, ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to count people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(fmt.Sprintf("%d", total))
	})

	app.Post("/pessoas", func(c *fiber.Ctx) error {
		p := new(people.PersonDTO)

		if err := c.BodyParser(p); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": "invalid json"})
		}

		if p.Name == nil || p.Nickname == nil {
			return c.SendStatus(fiber.StatusUnprocessableEntity)
		}

		if errs := validators.Validate(validate, p); len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
		}

		ctx := context.Background()
		uid, err := people.InsertPerson(*p, queries, ctx)
		if errors.Is(err, people.ErrUniqueNickname) {
			return c.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{"errors": map[string]string{"nickname": "already exists"}})
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to insert person")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Location(fmt.Sprintf("/pessoas/%s", uid))
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/pessoas/:id", func(c *fiber.Ctx) error {
		ctx := context.Background()
		p, err := people.FindPerson(c.Params("id"), queries, ctx)

		if err == people.ErrNotFound {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to find person")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(p)
	})

	app.Get("/pessoas", func(c *fiber.Ctx) error {
		term := c.Query("t")
		if strings.TrimSpace(term) == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		ppl, err := people.SearchPeople(term, queries, context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to search people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(ppl)
	})

	defer tools.StartProfiling(os.Getenv("PROFILER_MODE")).Stop()

	if err := app.Listen(":80"); err != nil {
		log.Warn().Err(err).Msg("failed to start server")
		os.Exit(1)
	}
}
