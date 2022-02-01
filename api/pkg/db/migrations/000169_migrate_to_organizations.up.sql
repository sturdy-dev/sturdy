BEGIN;

-- Add temporary columns
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS tmp_github_installation_id TEXT UNIQUE;
ALTER TABLE organizations
    ADD COLUMN IF NOT EXISTS tmp_user_id TEXT UNIQUE;
ALTER TABLE organizations
    ALTER COLUMN created_by DROP NOT NULL;

-- Reset organizations and add unique index/constraint
TRUNCATE organizations;
TRUNCATE organization_members;

UPDATE codebases
SET organization_id = NULL;

CREATE UNIQUE INDEX IF NOT EXISTS organization_members_organization_id_user_id_idx ON organization_members (organization_id, user_id);

-- Create one organization for each github_installation
INSERT INTO organizations (id, name, created_at, created_by, short_id, tmp_github_installation_id)
SELECT md5(random()::text || clock_timestamp()::text)::uuid,
       ghi.owner,
       NOW(),
       NULL, -- created_by is calculated later
       substring(md5(random()::text) from 0 for 7),
       ghi.id
FROM github_installations ghi
WHERE ghi.uninstalled_at IS NULL;

-- Set the organization for all codebases with a github_installation
WITH subquery AS (
    SELECT c.id as codebase_id, o.id as organization_id
    FROM codebases c
             JOIN github_repositories gr on c.id = gr.codebase_id
             JOIN github_installations gi on gi.installation_id = gr.installation_id
             JOIN organizations o on o.tmp_github_installation_id = gi.id
    WHERE c.organization_id IS NULL
)
UPDATE codebases
SET organization_id = subquery.organization_id
FROM subquery
WHERE codebases.id = subquery.codebase_id;

-- Create an organization for all codebases with only one member
-- Named "$USER's Team"
INSERT INTO organizations (id, name, created_at, created_by, short_id, tmp_user_id)
SELECT md5(random()::text || clock_timestamp()::text)::uuid,
       create_orgs.org_name,
       NOW(),
       create_orgs.user_id,
       substring(md5(random()::text) from 0 for 7),
       create_orgs.user_id
FROM (
         SELECT sub_users.*,
                CASE
                    WHEN POSITION(' ' IN sub_users.name) > 0
                        THEN SUBSTR(sub_users.name, 0, POSITION(' ' IN sub_users.name))
                    ELSE sub_users.name
                    END || '''s Team' AS org_name
         FROM (
                  SELECT cu.user_id, u.name
                  FROM (
                           SELECT c.id
                           FROM codebases c
                                    JOIN codebase_users cu on c.id = cu.codebase_id
                           WHERE c.organization_id IS NULL
                           GROUP BY c.id
                           HAVING count(*) = 1) as sub_c
                           JOIN codebase_users cu on cu.codebase_id = sub_c.id
                           JOIN users u on cu.user_id = u.id
                  GROUP BY cu.user_id, u.name
              ) AS sub_users) as create_orgs;

-- Assign codebases to single-user-orgs
WITH subquery AS (
    SELECT sub_c.id cb_id, cu.user_id, org.id org_id, org.name
    FROM (
             SELECT c.id
             FROM codebases c
                      JOIN codebase_users cu on c.id = cu.codebase_id
             WHERE c.organization_id IS NULL
             GROUP BY c.id
             HAVING count(*) = 1) as sub_c
             JOIN codebase_users cu on cu.codebase_id = sub_c.id
             JOIN organizations org on org.tmp_user_id = cu.user_id)
UPDATE codebases
SET organization_id = subquery.org_id
FROM subquery
WHERE id = subquery.cb_id;

-- Make the owner a member of all "single user orgs"
INSERT INTO organization_members (id, user_id, organization_id, created_at, created_by)
SELECT md5(random()::text || clock_timestamp()::text)::uuid,
       tmp_user_id as user_id,
       id          as organization_id,
       NOW(),
       tmp_user_id as created_by
FROM organizations
WHERE tmp_user_id IS NOT NULL;

-- Make the first member of (any) codebase a member of the organization
INSERT INTO organization_members (id, user_id, organization_id, created_at, created_by)
SELECT md5(random()::text || clock_timestamp()::text)::uuid,
       sub.user_id,
       sub.organization_id,
       NOW(),
       sub.user_id
FROM (
         SELECT distinct on (c.id) c.id, c.name, cu.user_id, u.email, c.organization_id
         FROM codebases c
                  JOIN codebase_users cu ON cu.codebase_id = c.id
                  JOIN users u on u.id = cu.user_id
         WHERE c.organization_id IS NOT NULL
         ORDER BY c.id, cu.created_at ASC) as sub
ON CONFLICT (organization_id, user_id) DO NOTHING;

-- Set created_by for organizations that doesn't have it set
WITH subquery AS(
    SELECT o.id as org_id, cu.user_id
    FROM organizations o
             JOIN codebases c on o.id = c.organization_id
             JOIN codebase_users cu on c.id = cu.codebase_id
    WHERE o.created_by IS NULL
)
UPDATE organizations
SET created_by = subquery.user_id
FROM subquery
WHERE id = subquery.org_id;

-- Remove temporary columns
-- ALTER TABLE organizations
--     DROP
--         COLUMN tmp_github_installation_id,
--     DROP
--         COLUMN tmp_user_id;

-- TODO(gustav): add this in a future migration, for now there are still some
--               codebases that are not in an organization.
-- Require all codebases to be in an organization
-- ALTER TABLE codebases
--     ALTER COLUMN organization_id SET NOT NULL;

COMMIT;