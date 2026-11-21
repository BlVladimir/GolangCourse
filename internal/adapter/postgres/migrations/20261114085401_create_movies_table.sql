-- +goose Up
CREATE TABLE movies (
    id bigint PRIMARY KEY,
    title text NOT NULL,
    year integer NOT NULL CHECK (1888 <= year),
    director text NOT NULL
);

-- +goose Down
DROP TABLE movie;
