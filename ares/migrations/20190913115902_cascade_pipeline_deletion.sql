-- +goose Up
-- SQL in this section is executed when the migration is applied.
alter table configuration drop constraint configuration_pipeline_id_fkey;
alter table configuration
    add constraint configuration_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id)
    on delete cascade;

alter table base_configuration drop constraint base_configuration_configuration_id_fkey;
alter table base_configuration
    add constraint base_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id)
    on delete cascade;

alter table trigger_configuration drop constraint trigger_configuration_configuration_id_fkey;
alter table trigger_configuration
    add constraint trigger_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id)
    on delete cascade;

alter table action_configuration drop constraint action_configuration_configuration_id_fkey;
alter table action_configuration
    add constraint action_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id)
    on delete cascade;

alter table deployment drop constraint deployment_configuration_id_fkey;
alter table deployment
    add constraint deployment_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id)
    on delete cascade;

alter table deployment drop constraint deployment_pipeline_id_fkey;
alter table deployment
    add constraint deployment_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id)
    on delete cascade;

alter table environment_variable drop constraint environment_variable_pipeline_id_fkey;
alter table environment_variable
    add constraint environment_variable_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id)
    on delete cascade;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
alter table configuration drop constraint configuration_pipeline_id_fkey;
alter table configuration
    add constraint configuration_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id);

alter table base_configuration drop constraint base_configuration_configuration_id_fkey;
alter table base_configuration
    add constraint base_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id);

alter table trigger_configuration drop constraint trigger_configuration_configuration_id_fkey;
alter table trigger_configuration
    add constraint trigger_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id);

alter table action_configuration drop constraint action_configuration_configuration_id_fkey;
alter table action_configuration
    add constraint action_configuration_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id);

alter table deployment drop constraint deployment_configuration_id_fkey;
alter table deployment
    add constraint deployment_configuration_id_fkey
    foreign key (configuration_id)
    references configuration (id);

alter table deployment drop constraint deployment_pipeline_id_fkey;
alter table deployment
    add constraint deployment_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id);

alter table environment_variable drop constraint environment_variable_pipeline_id_fkey;
alter table environment_variable
    add constraint environment_variable_pipeline_id_fkey
    foreign key (pipeline_id)
    references pipeline (id);

