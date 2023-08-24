CREATE EXTENSION IF NOT EXISTS pg_trgm;

DROP TABLE IF EXISTS people;

CREATE TABLE IF NOT EXISTS people (
    id VARCHAR(36) NOT NULL,
    nickname VARCHAR(32) CONSTRAINT nickname_pk PRIMARY KEY,
    fullname VARCHAR(100) NOT NULL,
    birthdate DATE NOT NULL,
    stack VARCHAR(1024),
    term_search VARCHAR(1158) GENERATED ALWAYS AS (LOWER(fullname) || ' ' || LOWER(nickname) || ' ' || LOWER(stack)) STORED
);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_term_search ON people USING GIN (term_search gin_trgm_ops);
