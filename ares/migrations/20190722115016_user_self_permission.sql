-- This migration resolves the issue solved in PR 141
-- https://bitbucket.org/unchain/ares/pull-requests/141/tbg-326-add-enforcermakeuser-to-invited/diff
-- The issue related to the role not being created when a user was invited.

-- +goose Up
-- +goose StatementBegin
CREATE or REPLACE function casbin_role() RETURNS bool as $$

BEGIN
    IF EXISTS (select 1 from pg_class where relname='casbin_rule') THEN
        INSERT into casbin_rule ( p_type, v0, v1, v2, v3, v4, v5 )
        SELECT 'g', concat('*::', id) "fmt", 'User', '', '', '', '' 
        FROM users
        WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE v0=concat('*::', id));
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

SELECT '' as nothing;