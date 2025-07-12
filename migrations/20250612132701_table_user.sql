-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
(
    id              bigserial
        constraint users_pk primary key,
    phone           varchar not null,
    email           varchar default ''
        constraint users_email_unique unique,
    password        varchar not null,
    salt            bytea,
    role_id         smallint,
    created_at      timestamp default now(),
    updated_at      timestamp,
    deleted_at      timestamp,
    last_visited_at timestamp,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
