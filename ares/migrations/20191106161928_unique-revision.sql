-- +goose Up
-- SQL in this section is executed when the migration is applied.
alter table configuration
    add constraint revision_pipeline_id unique (pipeline_id, revision);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
alter table configuration
    drop constraint revision_pipeline_id;
