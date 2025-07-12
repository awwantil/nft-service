-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS roles
(
    id   smallserial
        constraint roles_pk primary key,
    name varchar
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS roles;
-- +goose StatementEnd
