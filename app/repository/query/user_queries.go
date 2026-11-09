package query

const CheckUsernameExist = `
SELECT uid
FROM "User".users
WHERE username = $1;
`

const CheckEmailExist = `
SELECT uid
FROM "User".users
WHERE email = $1;
`

const InsertUser = `
INSERT INTO "User".users
	(first_name, last_name, username, email, role_id, password)
VALUES
	($1, $2, $3, $4, $5, $6)
RETURNING uid;
`

const GetUserByUsername = `
SELECT uid, first_name, last_name, username, email, role_id, password
FROM "User".users
WHERE username = $1;
`

const GetUserByUID = `
SELECT uid, first_name, last_name, username, email, role_id, password
FROM "User".users
WHERE uid = $1;
`

const UpdateUserByUID = `
UPDATE "User".users
SET first_name = $1,
	last_name = $2,
	username = $3,
	email = $4,
	role_id = $5,
	password = $6
WHERE uid = $7;
`

const DeleteUserByUID = `
DELETE FROM "User".users
WHERE uid = $1;
`
