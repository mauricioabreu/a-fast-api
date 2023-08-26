-- name: CountPeople :one
SELECT COUNT(1) FROM people;

-- name: InsertPerson :exec
INSERT INTO people (
    id, nickname, name, birthdate, stack
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: FindPerson :one
SELECT id, nickname, name, birthdate, stack FROM people WHERE id = $1;

-- name: SearchPeople :many
SELECT id, nickname, name, birthdate, stack FROM people WHERE term_search LIKE $1 LIMIT 50;
