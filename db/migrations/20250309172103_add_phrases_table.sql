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

-- Create an FTS5 virtual table for full-text search on phrases
CREATE VIRTUAL TABLE phrases_fts USING fts5(
  phrase,
  content='phrases',
  content_rowid='id'
);

-- Indexes for better performance
CREATE INDEX idx_phrases_foundIn ON phrases(foundIn);
CREATE INDEX idx_phrase_tags_tagName ON phrase_tags(tagName);

-- Trigger to keep FTS table in sync
CREATE TRIGGER phrases_ai AFTER INSERT ON phrases
BEGIN
  INSERT INTO phrases_fts(rowid, phrase) VALUES (new.id, new.phrase);
end
;

CREATE TRIGGER phrases_ad AFTER DELETE ON phrases
BEGIN
  DELETE FROM phrases_fts WHERE rowid = old.id;
end
;

CREATE TRIGGER phrases_au AFTER UPDATE ON phrases
BEGIN
  UPDATE phrases_fts SET phrase = new.phrase WHERE rowid = old.id;
end
;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER phrases_ai;
DROP TRIGGER phrases_ad;
DROP TRIGGER phrases_au;
DROP INDEX idx_phrases_foundIn;
DROP INDEX idx_phrase_tags_tagName;
DROP TABLE phrases_fts;
DROP TABLE phrase_tags;
DROP TABLE phrases;

-- +goose StatementEnd


