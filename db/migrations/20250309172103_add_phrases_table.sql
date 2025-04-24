-- +goose Up
-- +goose StatementBegin
CREATE TABLE phrases (
  id TEXT PRIMARY KEY,
  userId TEXT NOT NULL,
  phrase VARCHAR(255) NOT NULL,
  phraseDefinition VARCHAR(255) NOT NULL,
  pinned BOOLEAN NOT NULL DEFAULT FALSE,
  foundIn VARCHAR(255),
  public BOOLEAN NOT NULL DEFAULT FALSE,
  usageCount INT NOT NULL DEFAULT 0,
  createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updatedAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deletedAt DATETIME,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE phrase_tags (
  id TEXT PRIMARY KEY,
  phraseId TEXT NOT NULL,
  tagName TEXT NOT NULL,
  tagColor VARCHAR(255) NOT NULL,
  createdAt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (phraseId) REFERENCES phrases(id) ON DELETE CASCADE
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE phrase_tags;
DROP TABLE phrases;

-- +goose StatementEnd


