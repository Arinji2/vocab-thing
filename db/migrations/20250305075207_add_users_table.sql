-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id TEXT PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  createdAt DATETIME NOT NULL
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd


