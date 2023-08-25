package people

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/mauricioabreu/a-fast-api/db"
)

func InsertPerson(person PersonDTO, q *db.Queries, ctx context.Context) error {
	birthDate, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		return err
	}
	if err := q.InsertPerson(ctx, db.InsertPersonParams{
		ID:        "foo",
		Nickname:  person.Nickname,
		Name:      person.Name,
		Birthdate: birthDate,
		Stack:     sql.NullString{String: strings.Join(person.Stack, ","), Valid: true},
	}); err != nil {
		return err
	}

	return nil
}

func CountPeople(q *db.Queries, ctx context.Context) (int64, error) {
	return q.CountPeople(ctx)
}
