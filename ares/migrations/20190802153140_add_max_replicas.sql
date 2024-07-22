-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- add max_replicas column to environment
alter table environment
  add column max_replicas bigint not null default 1;

-- add max_replicas column to default_environment
alter table default_environment
  add column max_replicas bigint not null default 1;

-- update max replicas for existing production environments to 3
update environment
set max_replicas = 3
where name = 'production';

-- update default max replicas for production environments to 3
update default_environment
set max_replicas = 3
where name = 'production';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

alter table environment
  drop column max_replicas;

alter table default_environment
  drop column max_replicas;