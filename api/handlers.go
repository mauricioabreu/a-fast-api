package api

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/mauricioabreu/a-fast-api/api/people"
	"github.com/mauricioabreu/a-fast-api/validators"
	"github.com/rs/zerolog/log"
)

func CountPeople(service *PeopleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		total, err := people.CountPeople(service.dbc, ctx)
		if err != nil {
			log.Error().Err(err).Msg("failed to count people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendString(fmt.Sprintf("%d", total))
	}
}

func CreatePerson(service *PeopleService) fiber.Handler {
	validate := validator.New()
	return func(c *fiber.Ctx) error {
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
		uid, err := people.CreatePerson(*p, service.dbc, ctx)
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
	}
}

func FindPerson(service *PeopleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.Background()
		p, err := people.FindPerson(c.Params("id"), service.dbc, ctx)

		if err == people.ErrNotFound {
			return c.SendStatus(fiber.StatusNotFound)
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to find person")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(p)
	}
}

func SearchPeople(service *PeopleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		term := c.Query("t")
		if strings.TrimSpace(term) == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		ppl, err := people.SearchPeople(term, service.dbc, context.Background())
		if err != nil {
			log.Error().Err(err).Msg("failed to search people")
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(ppl)
	}
}
