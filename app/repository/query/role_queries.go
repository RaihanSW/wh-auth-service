package query

const InsertRole = `
INSERT INTO "User".roles (name)
VALUES ($1)
RETURNING id;
`

const GetRoleByID = `
SELECT id, name
FROM "User".roles
WHERE id = $1;
`

const GetAllRoles = `
SELECT id, name
FROM "User".roles
ORDER BY id;
`

const UpdateRoleByID = `
UPDATE "User".roles
SET name = $1
WHERE id = $2;
`

const DeleteRoleByID = `
DELETE FROM "User".roles
WHERE id = $1;
`
