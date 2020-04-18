ALTER TABLE managed_code
ADD COLUMN user_id BIGINT REFERENCES users(id);

-- Just assign it to the first admin in the organization.
-- This shouldn't really do anything in production but just in case.
UPDATE managed_code AS mc
SET user_id = sub.id
FROM (
    SELECT u.id, ur.org_id
    FROM users AS u
    INNER JOIN user_roles AS ur
        ON ur.user_id = u.id
    INNER JOIN organization_available_roles AS oar
        ON oar.id = ur.role_id
    WHERE oar.is_admin_role = true
    LIMIT 1
) AS sub
WHERE sub.org_id = mc.org_id;

ALTER TABLE managed_code
ALTER COLUMN user_id SET NOT NULL;
