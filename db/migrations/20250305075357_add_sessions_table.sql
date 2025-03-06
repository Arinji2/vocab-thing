-- +goose Up
-- +goose StatementBegin
CREATE TABLE sessions (
  id VARCHAR(255) PRIMARY KEY,
  userId VARCHAR(255) NOT NULL,
  providerId VARCHAR(255) NOT NULL,
  fingerprint VARCHAR(255) NOT NULL,
  ip VARCHAR(255) NOT NULL,
  expiresAt DATETIME NOT NULL,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (providerId) REFERENCES providers(id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE sessions;
-- +goose StatementEnd


