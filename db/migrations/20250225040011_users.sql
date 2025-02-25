-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id INTEGER PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  createdAt DATETIME NOT NULL
);

CREATE TABLE providers (
  id INTEGER PRIMARY KEY,
  userId INTEGER NOT NULL,
  type VARCHAR(255) NOT NULL,
  refreshToken VARCHAR(255) NOT NULL,
  accessToken VARCHAR(255) NOT NULL,
  expiresAt DATETIME NOT NULL,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE providers;
-- +goose StatementEnd


