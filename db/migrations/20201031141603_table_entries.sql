
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE entries (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP table entries;
