-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_tokens
(
    id                 bigserial
        constraint refresh_tokens_pk primary key,
    user_id            bigserial,
    token              varchar,
    expired_at         timestamp,
    refresh_token      varchar,
    refresh_expired_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_tokens;
-- +goose StatementEnd
