-- +goose Up
-- SQL in this section is executed when the migration is applied.
alter table users alter column email drop not null;

update users set email = null where email like 'archived@@%';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
update users set email = 'archived@@' || id where email is null;
alter table users
    alter column email set not null;
