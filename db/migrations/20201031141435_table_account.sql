
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE accounts (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
Drop table accounts;
