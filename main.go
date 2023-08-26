package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/db"
	"github.com/mauricioabreu/a-fast-api/people"
	"github.com/mauricioabreu/a-fast-api/validators"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func validateStack(fl validator.FieldLevel) bool {
	if fl.Field().IsZero() {
		return true
	}

	elements, ok := fl.Field().Interface().([]interface{})
	if !ok {
		return false
	}

	for _, elem := range elements {
		str, ok := elem.(string)
		if !ok || len(str) > 32 {
			return false
		}
	}

	return true
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	validate := validator.New()
	if err := validate.RegisterValidation("validateStack", validateStack); err != nil {
		log.Fatal().Err(err).Msg("failed to register validation")
	}

	app := fiber.New(fiber.Config{
		Prefork: true,
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
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Location(fmt.Sprintf("/pessoas/%s", uid))
		return c.SendStatus(fiber.StatusCreated)
	})

	app.Get("/pessoas/:id", func(c *fiber.Ctx) error {
		ctx := context.Background()
		p, err := people.FindPerson(c.Params("id"), queries, ctx)

		if err == sql.ErrNoRows {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if err != nil {
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
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(ppl)
	})

	app.Listen(":80")
}
