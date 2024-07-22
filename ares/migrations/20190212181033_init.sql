-- +goose Up
create type component_type as enum ('action', 'trigger', 'base');

create table users
(
    id                      bigserial   not null primary key,

    created_at              timestamp   not null,
    created_by_id           bigint,
    updated_at              timestamp   not null,
    updated_by_id           bigint,

    full_name               text        not null,
    password_hash           text        not null,
    email                   text unique not null,
    default_organization_id bigint,
    status                  text,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

create table account_confirmation_token
(
    id              bigserial   not null primary key,

    user_id         bigint      not null,
    token           text unique not null,
    expiration_time timestamp   not null,

    foreign key (user_id) references users (id)
);

create table organization
(
    id            bigserial   not null primary key,

    created_at    timestamp   not null,
    created_by_id bigint      not null,
    updated_at    timestamp   not null,
    updated_by_id bigint      not null,

    display_name  text        not null,
    name          text unique not null,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

alter table users
    add foreign key (default_organization_id) references organization (id);


create table subscription_plan
(
    id             bigserial   not null primary key,

    created_at     timestamp   not null,
    created_by_id  bigint      not null,
    updated_at     timestamp   not null,
    updated_by_id  bigint      not null,

    name           text unique not null,
    pipeline_limit bigint      not null,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

create table subscription
(
    id                   bigserial not null primary key,
    subscription_plan_id bigint    not null,
    organization_id      bigint    not null,


    created_at           timestamp not null,
    created_by_id        bigint    not null,
    updated_at           timestamp not null,
    updated_by_id        bigint    not null,

    start_date           timestamp not null,
    end_date             timestamp not null,

    foreign key (subscription_plan_id) references subscription_plan (id),
    foreign key (organization_id) references organization (id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

create table organization_billing_provider
(
    id              bigserial not null primary key,
    organization_id bigint    not null,

    created_at      timestamp not null,
    created_by_id   bigint    not null,
    updated_at      timestamp not null,
    updated_by_id   bigint    not null,

    provider_name   text      not null,
    billing_info    jsonb,

    foreign key (organization_id) references organization (id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

create table user_organization
(
    user_id         bigint not null,
    organization_id bigint not null,

    primary key (user_id, organization_id),
    foreign key (organization_id) references organization (id),
    foreign key (user_id) references users (id)
);

create table base
(
    id            bigserial   not null primary key,
    developer_id  bigint      not null,

    created_at    timestamp   not null,
    created_by_id bigint      not null,
    updated_at    timestamp   not null,
    updated_by_id bigint      not null,

    name          text unique not null,
    display_name  text        not null,
    description   text        not null,
    public        boolean     not null default false,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (developer_id) references organization (id)
);

create table base_version
(
    id               bigserial not null primary key,
    base_id          bigint    not null,

    created_at       timestamp not null,
    created_by_id    bigint    not null,
    updated_at       timestamp not null,
    updated_by_id    bigint    not null,

    version          text      not null,
    description      text      not null,
    readme           text      not null,
    docker_image_ref text      not null,
    entrypoint       text      not null,

    public           boolean   not null default false,

    unique (base_id, version),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (base_id) references base (id)
);

create table trigger
(
    id            bigserial   not null primary key,
    developer_id  bigint      not null,

    created_at    timestamp   not null,
    created_by_id bigint      not null,
    updated_at    timestamp   not null,
    updated_by_id bigint      not null,

    name          text unique not null,
    display_name  text        not null,
    description   text        not null,
    public        boolean     not null default false,


    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (developer_id) references organization (id)
);

create table trigger_version
(
    id             bigserial not null primary key,
    trigger_id     bigint    not null,

    created_at     timestamp not null,
    created_by_id  bigint    not null,
    updated_at     timestamp not null,
    updated_by_id  bigint    not null,

    version        text      not null,
    description    text      not null,
    public         boolean   not null default false,
    example_config text      not null,
    readme         text      not null,
    file_name      text      not null,
    file_id        text      not null,
    input_schema   jsonb,
    output_schema  jsonb,


    unique (trigger_id, version),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (trigger_id) references trigger (id)
);

create table trigger_version_supported_bases
(
    trigger_version_id bigint not null,
    supported_id       bigint not null,

    primary key (trigger_version_id, supported_id),
    foreign key (trigger_version_id) references trigger_version (id),
    foreign key (supported_id) references base_version (id)
);

create table action
(
    id            bigserial   not null primary key,
    developer_id  bigint      not null,

    created_at    timestamp   not null,
    created_by_id bigint      not null,
    updated_at    timestamp   not null,
    updated_by_id bigint      not null,

    name          text unique not null,
    display_name  text        not null,
    description   text        not null,
    public        boolean     not null default false,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (developer_id) references organization (id)
);


create table action_version
(
    id             bigserial not null primary key,
    action_id      bigint    not null,

    created_at     timestamp not null,
    created_by_id  bigint    not null,
    updated_at     timestamp not null,
    updated_by_id  bigint    not null,

    version        text      not null,
    public         boolean   not null default false,
    example_config text      not null,
    description    text      not null,
    readme         text      not null,
    file_name      text      not null,
    file_id        text      not null,
    input_schema   jsonb,
    output_schema  jsonb,

    unique (action_id, version),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (action_id) references action (id)
);

create table action_version_supported_bases
(
    action_version_id bigint not null,
    supported_id      bigint not null,

    primary key (action_version_id, supported_id),
    foreign key (action_version_id) references action_version (id),
    foreign key (supported_id) references base_version (id)
);

create table draft_configuration
(
    id              bigserial not null primary key,

    created_at      timestamp not null,
    created_by_id   bigint    not null,
    updated_at      timestamp not null,
    updated_by_id   bigint    not null,

    organization_id bigint    not null,

    revision        bigint    not null,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (organization_id) references organization (id)
);

create table pipeline
(
    id                     bigserial not null primary key,
    organization_id        bigint    not null,

    created_at             timestamp not null,
    created_by_id          bigint    not null,
    updated_at             timestamp not null,
    updated_by_id          bigint    not null,

    display_name           text      not null,
    name                   text      not null check ( char_length(name) <= 50 ),
    status                 text      not null,
    description            text      not null,
    draft_configuration_id bigint,

    unique (name, organization_id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (organization_id) references organization (id),
    foreign key (draft_configuration_id) references draft_configuration (id)
);

create table configuration
(
    id              bigserial not null primary key,

    created_at      timestamp not null,
    created_by_id   bigint    not null,
    updated_at      timestamp not null,
    updated_by_id   bigint    not null,

    pipeline_id     bigint    not null,
    organization_id bigint    not null,

    revision        bigint    not null,

    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    foreign key (pipeline_id) references pipeline (id),
    foreign key (organization_id) references organization (id)
);

create table action_configuration
(
    id               bigserial not null primary key,
    configuration_id bigint    not null,
    version_id       bigint    not null,

    index            bigint    not null,
    name             text      not null,
    config           text      not null,
    message_config   jsonb,

    foreign key (configuration_id) references configuration (id),
    foreign key (version_id) references action_version (id),
    unique (configuration_id, name)
);

create table base_configuration
(
    id               bigserial     not null primary key,
    configuration_id bigint unique not null,
    version_id       bigint,

    config           text          not null,

    foreign key (configuration_id) references configuration (id),
    foreign key (version_id) references base_version (id)
);

create table trigger_configuration
(
    id               bigserial     not null primary key,
    configuration_id bigint unique not null,
    version_id       bigint,

    name             text          not null,
    config           text          not null,
    message_config   jsonb,

    foreign key (configuration_id) references configuration (id),
    foreign key (version_id) references trigger_version (id),
    unique (configuration_id, name)
);


create table action_draft_configuration
(
    id                     bigserial not null primary key,
    draft_configuration_id bigint    not null,
    version_id             bigint    not null,

    index                  bigint    not null,
    name                   text      not null,
    config                 text      not null,
    message_config         jsonb,

    foreign key (draft_configuration_id) references draft_configuration (id),
    foreign key (version_id) references action_version (id),
    unique (draft_configuration_id, name)
);

create table base_draft_configuration
(
    id                     bigserial     not null primary key,
    draft_configuration_id bigint unique not null,
    version_id             bigint,

    config                 text          not null,

    foreign key (draft_configuration_id) references draft_configuration (id),
    foreign key (version_id) references base_version (id)
);

create table trigger_draft_configuration
(
    id                     bigserial     not null primary key,
    draft_configuration_id bigint unique not null,
    version_id             bigint,

    name                   text          not null,
    config                 text          not null,
    message_config         jsonb,

    foreign key (draft_configuration_id) references draft_configuration (id),
    foreign key (version_id) references trigger_version (id),
    unique (draft_configuration_id, name)
);

create table default_environment
(
    id    bigserial not null primary key,
    index bigint    not null,
    name  text      not null,

    unique (name)
);

create table environment
(
    id              bigserial not null primary key,
    index           bigint    not null,
    name            text      not null,
    organization_id bigint    not null,

    created_at      timestamp not null,
    created_by_id   bigint    not null,
    updated_at      timestamp not null,
    updated_by_id   bigint    not null,

    unique (name, organization_id),
    foreign key (organization_id) references organization (id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);

create table environment_variable
(
    id              bigserial not null primary key,
    organization_id bigint    not null,
    pipeline_id     bigint    not null,
    environment_id  bigint    not null,

    created_at      timestamp not null,
    created_by_id   bigint    not null,
    updated_at      timestamp not null,
    updated_by_id   bigint    not null,

    index           bigint    not null,
    key             text      not null,
    value           text      not null,
    secret          boolean   not null,
    deployed        boolean   not null,

    unique (key, pipeline_id, environment_id),
    foreign key (organization_id) references organization (id),
    foreign key (pipeline_id) references pipeline (id),
    foreign key (environment_id) references environment (id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id)
);
--
-- create table organization_environment_secret
-- (
--   id                          bigserial not null primary key,
--   organization_environment_id bigint    not null,
--
--   name                        text   not null,
--   value                       text   not null,
--
--   foreign key (organization_environment_id) references organization_environment (id)
-- );

create table deployment
(
    id               bigserial primary key,
    pipeline_id      bigint    not null,
    configuration_id bigint    not null,
    environment_id   bigint    not null,

    created_at       timestamp not null,
    created_by_id    bigint    not null,
    updated_at       timestamp not null,
    updated_by_id    bigint    not null,

    replicas         bigint    not null,
    url              text      not null,
    image            text      not null,
    host             text      not null,
    path             text      not null,
    rewrite_target   text      not null,
    full_name        text      not null,
    dirty            boolean   not null,

    foreign key (configuration_id) references configuration (id),
    foreign key (environment_id) references environment (id),
    foreign key (pipeline_id) references pipeline (id),
    foreign key (created_by_id) references users (id),
    foreign key (updated_by_id) references users (id),
    unique (pipeline_id, environment_id)
);
--
-- create table environment_deployment
-- (
--
--   configuration_id bigint not null,
--   environment_id     bigint not null,
--   id                             bigserial not null primary key,
--   configuration_environment_id bigint    not null,
--
--   replicas                       int    not null,
--   url                            text   not null,
--   image                          text   not null,
--   host                           text   not null,
--   path                           text   not null,
--   status                         text   not null,
--
--   --   start_time timestamp,
--   --   stop_time timestamp
--
--   foreign key (configuration_environment_id) references configuration_environment (id)
-- );

--
-- -- +goose StatementBegin
-- create or replace function create_role_for_organization()
--   returns trigger as
-- $BODY$
-- begin
--   execute 'create role ' || quote_ident(NEW.name);
--   return NEW;
-- end;
-- $BODY$
--   language plpgsql;
-- -- +goose StatementEnd
--
-- create trigger insert_organization
--   after insert
--   on organization
--   for each row
-- execute procedure create_role_for_organization();
--
-- -- +goose StatementBegin
-- create or replace function drop_role_for_organization()
--   returns trigger as
-- $BODY$
-- begin
--   execute 'drop role ' || quote_ident(NEW.name);
--   return NEW;
-- end;
-- $BODY$
--   language plpgsql;
-- -- +goose StatementEnd
--
-- create trigger delete_organization
--   after delete
--   on organization
--   for each row
-- execute procedure drop_role_for_organization();
;
-- +goose Down
drop table if exists deployment;
drop table if exists environment_variable;
drop table if exists environment;
drop table if exists default_environment;
drop table if exists trigger_draft_configuration;
drop table if exists base_draft_configuration;
drop table if exists action_draft_configuration;
drop table if exists trigger_configuration;
drop table if exists base_configuration;
drop table if exists action_configuration;
drop table if exists configuration;
drop table if exists pipeline;
drop table if exists draft_configuration;
drop table if exists action_version_supported_bases;
drop table if exists action_version;
drop table if exists action;
drop table if exists trigger_version_supported_bases;
drop table if exists trigger_version;
drop table if exists trigger;
drop table if exists base_version;
drop table if exists base;
drop table if exists user_organization;
drop table if exists organization_billing_provider;
drop table if exists subscription;
drop table if exists subscription_plan;

alter table users
    drop constraint users_default_organization_id_fkey;
drop table if exists organization;
drop table if exists account_confirmation_token;
drop table if exists users;
drop type if exists component_type;
