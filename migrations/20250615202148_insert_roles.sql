-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (id, name)
VALUES (1, 'user'),
       (2, 'creator'),
       (99, 'moderator'),
       (100, 'admin');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE
FROM roles;

ALTER SEQUENCE IF EXISTS roles_id_seq RESTART;
-- +goose StatementEnd
