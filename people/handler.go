package people

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mauricioabreu/a-fast-api/db"
)

var (
	ErrUniqueNickname = fmt.Errorf("nickname already exists")
)

func InsertPerson(person PersonDTO, q *db.Queries, ctx context.Context) (string, error) {
	birthDate, err := time.Parse("2006-01-02", person.Birthdate)
	if err != nil {
		return "", err
	}
	uid := uuid.New().String()
	if err := q.InsertPerson(ctx, db.InsertPersonParams{
		ID:        uid,
		Nickname:  *person.Nickname,
		Name:      *person.Name,
		Birthdate: birthDate,
		Stack:     sql.NullString{String: strings.Join(person.Stack, ","), Valid: true},
	}); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" && pqErr.Constraint == "nickname_pk" {
				return "", ErrUniqueNickname
			}
		}
		return "", err
	}

	return uid, nil
}

func CountPeople(q *db.Queries, ctx context.Context) (int64, error) {
	return q.CountPeople(ctx)
}

func FindPerson(uid string, q *db.Queries, ctx context.Context) (*PersonDTO, error) {
	person, err := q.FindPerson(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &PersonDTO{
		ID:        person.ID,
		Nickname:  &person.Nickname,
		Name:      &person.Name,
		Birthdate: person.Birthdate.Format("2006-01-02"),
		Stack:     strings.Split(person.Stack.String, ","),
	}, nil
}

func SearchPeople(term string, q *db.Queries, ctx context.Context) ([]*PersonDTO, error) {
	p, err := q.SearchPeople(ctx, sql.NullString{String: fmt.Sprintf("%%%s%%", term), Valid: true})
	if err != nil {
		return nil, err
	}

	ppl := make([]*PersonDTO, 0, len(p))

	for _, person := range p {
		ppl = append(ppl, &PersonDTO{
			ID:        person.ID,
			Nickname:  &person.Nickname,
			Name:      &person.Name,
			Birthdate: person.Birthdate.Format("2006-01-02"),
			Stack:     strings.Split(person.Stack.String, ","),
		})
	}

	return ppl, nil
}
