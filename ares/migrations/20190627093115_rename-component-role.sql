-- Rename casbin rules.

-- +goose Up
-- +goose StatementBegin
CREATE or REPLACE function casbin_role() RETURNS bool as $$

BEGIN
    IF EXISTS (select 1 from pg_class where relname='casbin_rule') THEN
        update casbin_rule set v0='PipelineOperator' where v0='RolePipelineOperator';
        update casbin_rule set v1='PipelineOperator' where v1='RolePipelineOperator';        
        RETURN true;
    ELSE 
        RETURN false;
    END IF;
END; $$

LANGUAGE PLPGSQL;

SELECT casbin_role();
DROP FUNCTION IF EXISTS casbin_role;
-- +goose StatementEnd


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- +goose StatementBegin
CREATE or REPLACE function casbin_role() RETURNS bool as $$

BEGIN
    IF EXISTS (select 1 from pg_class where relname='casbin_rule') THEN
        update casbin_rule set v0='RolePipelineOperator' where v0='PipelineOperator';
        update casbin_rule set v1='RolePipelineOperator' where v1='PipelineOperator';
        RETURN true;
    ELSE 
        RETURN false;
    END IF;
END; $$

LANGUAGE PLPGSQL;
-- +goose StatementEnd

SELECT casbin_role();
DROP FUNCTION IF EXISTS casbin_role;



