-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- add status
alter table subscription
    add column status text;

-- make awsCustomerId unique
create unique index obp_aws_customer_id on organization_billing_provider ((billing_info ->> 'awsCustomerId'));


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

alter table subscription
    drop column status;

drop index if exists obp_aws_customer_id;
