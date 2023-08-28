package api

import "github.com/mauricioabreu/a-fast-api/db"

type PeopleService struct {
	dbc *db.Queries
}

func NewPeopleService(queries *db.Queries) *PeopleService {
	return &PeopleService{dbc: queries}
}
