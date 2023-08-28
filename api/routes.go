package api

import (
	"github.com/gofiber/fiber/v2"
)

func PeopleRouter(app fiber.Router, service *PeopleService) {
	app.Get("/contagem-pessoas", CountPeople(service))
	app.Post("/pessoas", CreatePerson(service))
	app.Get("/pessoas/:id", FindPerson(service))
	app.Get("/pessoas", SearchPeople(service))
}
