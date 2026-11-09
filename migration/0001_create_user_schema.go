package migration

import "database/sql"

func init() {
	Register(Migration{
		Version: "202611090001",
		Name:    "create_user_schema_tables",
		DBKey:   "warehouse_db",
		Up: func(tx *sql.Tx) error {
			stmts := []string{
				`CREATE SCHEMA IF NOT EXISTS "User"`,
				`CREATE TABLE IF NOT EXISTS "User".roles (
					id   INT PRIMARY KEY,
					name VARCHAR(32) NOT NULL
				)`,
				`CREATE TABLE IF NOT EXISTS "User".users (
					uid        VARCHAR(16) PRIMARY KEY DEFAULT upper(substr(md5(random()::text), 1, 16)),
					first_name VARCHAR(64) NOT NULL,
					last_name  VARCHAR(32),
					username   VARCHAR(32) NOT NULL UNIQUE,
					email      VARCHAR(255) NOT NULL UNIQUE,
					role_id    INT         NOT NULL REFERENCES "User".roles (id),
					password   VARCHAR(64) NOT NULL
				)`,
			}
			for _, stmt := range stmts {
				if _, err := tx.Exec(stmt); err != nil {
					return err
				}
			}
			return nil
		},
		Down: func(tx *sql.Tx) error {
			stmts := []string{
				`DROP TABLE IF EXISTS "User".users`,
				`DROP TABLE IF EXISTS "User".roles`,
				`DROP SCHEMA IF EXISTS "User"`,
			}
			for _, stmt := range stmts {
				if _, err := tx.Exec(stmt); err != nil {
					return err
				}
			}
			return nil
		},
	})
}
