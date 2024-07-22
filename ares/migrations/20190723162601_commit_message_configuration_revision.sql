-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- add commit_message
alter table configuration
    add column commit_message text not null default '';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

alter table configuration
    drop column commit_message;
