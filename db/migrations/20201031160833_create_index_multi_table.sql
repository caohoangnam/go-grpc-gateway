
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX ON accounts ("owner");

CREATE INDEX ON entries ("account_id");

CREATE INDEX ON transfers ("from_account_id");

CREATE INDEX ON transfers ("to_account_id");

CREATE INDEX ON transfers ("from_account_id", "to_account_id");

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX accounts_owner_idx;
DROP INDEX entries_account_id_idx;
DROP INDEX transfers_from_account_id_idx;
DROP INDEX transfers_to_account_id_idx;
DROP INDEX transfers_from_account_id_to_account_id_idx;
