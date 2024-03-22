-- +goose Up
CREATE TABLE IF NOT EXISTS "words" (
  "id" SERIAL PRIMARY KEY NOT NULL,
  "sections" JSONB,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS "words";
