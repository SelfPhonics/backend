-- +goose Up
CREATE TRIGGER set_timestamp
BEFORE UPDATE ON words
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose Down
DROP TRIGGER IF EXISTS set_timestamp ON words;
