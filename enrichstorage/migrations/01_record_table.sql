-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS record (
    surname VARCHAR NOT NULL,
    "name" VARCHAR NOT NULL,
    patronymic VARCHAR NOT NULL,
    age INT,
    sex VARCHAR,
    nationality VARCHAR,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    PRIMARY KEY (surname, "name", patronymic)
);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE record;
-- +migrate StatementEnd