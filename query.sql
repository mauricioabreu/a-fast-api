-- name: CountPeople :one
SELECT COUNT(1) FROM people;

-- name: InsertPerson :exec
INSERT INTO people (
    id, nickname, name, birthdate, stack
) VALUES (
    $1, $2, $3, $4, $5
);
