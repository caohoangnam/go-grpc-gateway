
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE transfers (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP table transfers;
