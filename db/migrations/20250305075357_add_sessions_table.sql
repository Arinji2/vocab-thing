-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  id TEXT PRIMARY KEY,
  userId TEXT PRIMARY KEY,
  providerId TEXT PRIMARY KEY,
  fingerprint VARCHAR(255) NOT NULL,
  ip VARCHAR(255) NOT NULL,
  expiresAt DATETIME NOT NULL,
  createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (providerId) REFERENCES providers(id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd


