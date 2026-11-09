package seed

import (
	"context"
	"database/sql"
	"dc-transaction-service/util"
)

func init() {
	Register(Seeder{
		Name:  "users",
		DBKey: "warehouse_db",
		Run: func(ctx context.Context, db *sql.DB) error {
			users := []struct {
				firstName string
				lastName  string
				username  string
				email     string
				roleID    int
				password  string
			}{
				{"Super", "Admin", "superadmin", "superadmin@warehouse.local", 1, "admin123"},
				{"Admin", "User", "admin", "admin@warehouse.local", 2, "admin123"},
				{"Staff", "User", "staff", "staff@warehouse.local", 3, "staff123"},
			}

			for _, u := range users {
				// uid di-generate oleh column DEFAULT di Postgres
				_, err := db.ExecContext(ctx, `
					INSERT INTO "User".users (first_name, last_name, username, email, role_id, password)
					VALUES ($1, $2, $3, $4, $5, $6)
					ON CONFLICT (username) DO NOTHING
				`, u.firstName, u.lastName, u.username, u.email, u.roleID, util.HashSHA256(u.password))
				if err != nil {
					return err
				}
			}
			return nil
		},
	})
}
