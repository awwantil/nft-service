-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS nft_data
(
    id              bigserial
        constraint nft_pk primary key,
    token_id        bigint  default 0,
    content         text    default '',
    cidv0           varchar not null,
    cidv1           varchar not null,
    file_size       varchar default '',
    file_name       varchar default '',
    created_at      timestamp default now(),
    updated_at      timestamp,
    deleted_at      timestamp,
    last_visited_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS nft_data;
-- +goose StatementEnd
