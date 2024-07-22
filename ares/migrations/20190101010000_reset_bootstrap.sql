-- +goose Up
-- SQL in this section is executed when the migration upped.

select '' as nothing;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

drop table if exists casbin_rule;
drop table if exists goose_fixture_version;
