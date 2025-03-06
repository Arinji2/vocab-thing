-- +goose Up
-- +goose StatementBegin
CREATE TABLE providers (
  id TEXT PRIMARY KEY,
  userId TEXT PRIMARY KEY,
  type VARCHAR(255) NOT NULL,
  refreshToken VARCHAR(255) NOT NULL,
  accessToken VARCHAR(255) NOT NULL,
  expiresAt DATETIME NOT NULL,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE providers;
-- +goose StatementEnd


