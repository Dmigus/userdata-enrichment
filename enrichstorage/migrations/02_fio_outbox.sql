-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE IF NOT EXISTS fio_outbox
(
    id BIGSERIAL NOT NULL primary key,
    payload BYTEA NOT NULL
);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS fio_outbox;
-- +migrate StatementEnd