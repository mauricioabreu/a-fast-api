package people

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mauricioabreu/a-fast-api/db"
)

func InsertPerson(person PersonDTO, q *db.Queries, ctx context.Context) (string, error) {
	birthDate, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		return "", err
	}
	uid := uuid.New().String()
	if err := q.InsertPerson(ctx, db.InsertPersonParams{
		ID:        uid,
		Nickname:  person.Nickname,
		Name:      person.Name,
		Birthdate: birthDate,
		Stack:     sql.NullString{String: strings.Join(person.Stack, ","), Valid: true},
	}); err != nil {
		return "", err
	}

	return uid, nil
}

func CountPeople(q *db.Queries, ctx context.Context) (int64, error) {
	return q.CountPeople(ctx)
}

func FindPerson(uid string, q *db.Queries, ctx context.Context) (*PersonDTO, error) {
	p, err := q.FindPerson(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &PersonDTO{
		ID:        p.ID,
		Nickname:  p.Nickname,
		Name:      p.Name,
		Birthdate: p.Birthdate.Format("2006-01-02"),
		Stack:     strings.Split(p.Stack.String, ","),
	}, nil
}
