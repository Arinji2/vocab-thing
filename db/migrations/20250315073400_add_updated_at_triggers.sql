-- +goose Up
-- +goose StatementBegin
CREATE TABLE sync_metadata (
  id TEXT PRIMARY KEY,
  userId TEXT NOT NULL,
  last_updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TRIGGER update_phrases_timestamp
AFTER INSERT ON phrases
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = NEW.userId;
end
;

CREATE TRIGGER update_phrases_timestamp_update
AFTER UPDATE ON phrases
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = NEW.userId;
end
;

CREATE TRIGGER update_phrases_timestamp_delete
AFTER DELETE ON phrases
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = OLD.userId;
end
;

CREATE TRIGGER update_phrase_tags_timestamp
AFTER INSERT ON phrase_tags
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = (SELECT userId FROM phrases WHERE id = NEW.phraseId);
end
;

CREATE TRIGGER update_phrase_tags_timestamp_update
AFTER UPDATE ON phrase_tags
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = (SELECT userId FROM phrases WHERE id = NEW.phraseId);
end
;

CREATE TRIGGER update_phrase_tags_timestamp_delete
AFTER DELETE ON phrase_tags
BEGIN
    UPDATE sync_metadata SET last_updated_at = CURRENT_TIMESTAMP WHERE userId = (SELECT userId FROM phrases WHERE id = OLD.phraseId);
end
;
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_phrases_timestamp;
DROP TRIGGER IF EXISTS update_phrases_timestamp_update;
DROP TRIGGER IF EXISTS update_phrases_timestamp_delete;
DROP TRIGGER IF EXISTS update_phrase_tags_timestamp;
DROP TRIGGER IF EXISTS update_phrase_tags_timestamp_update;
DROP TRIGGER IF EXISTS update_phrase_tags_timestamp_delete;
DROP TABLE IF EXISTS sync_metadata;
-- +goose StatementEnd


