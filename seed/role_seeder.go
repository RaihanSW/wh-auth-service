package seed

import (
	"context"
	"database/sql"
)

func init() {
	Register(Seeder{
		Name:  "roles",
		DBKey: "warehouse_db",
		Run: func(ctx context.Context, db *sql.DB) error {
			roles := []struct {
				id   int
				name string
			}{
				{1, "Superadmin"},
				{2, "Admin"},
				{3, "Staff"},
			}

			for _, r := range roles {
				_, err := db.ExecContext(ctx, `
					INSERT INTO "User".roles (id, name)
					VALUES ($1, $2)
					ON CONFLICT (id) DO NOTHING
				`, r.id, r.name)
				if err != nil {
					return err
				}
			}
			return nil
		},
	})
}
