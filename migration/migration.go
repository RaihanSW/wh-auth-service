package migration

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"time"
)

// Migration is a single schema change bound to one database key from
// application_<env>.json (e.g. "auth_user", "dc21_transaction").
// Register new migrations from init() in files named <version>_<name>.go.
type Migration struct {
	Version string // timestamp-style, e.g. 202611090001 — applied ascending
	Name    string // e.g. "create_users_table"
	DBKey   string // cfg.Database key — always the writer, never the replica
	Up      func(tx *sql.Tx) error
	Down    func(tx *sql.Tx) error
}

var registry []Migration

func Register(m Migration) {
	registry = append(registry, m)
}

// All returns registered migrations sorted by version ascending.
func All() []Migration {
	ms := make([]Migration, len(registry))
	copy(ms, registry)
	sort.Slice(ms, func(i, j int) bool { return ms[i].Version < ms[j].Version })
	return ms
}

const createTrackingTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version    VARCHAR(32) NOT NULL,
	name       VARCHAR(255) NOT NULL,
	applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`

// RunUp applies every pending migration in version order.
func RunUp(ctx context.Context, dbs map[string]*sql.DB) error {
	applied := 0
	for _, m := range All() {
		done, err := isApplied(ctx, dbs, m)
		if err != nil {
			return err
		}
		if done {
			continue
		}

		db, err := dbFor(dbs, m.DBKey)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("migrate up %s (%s): begin: %w", m.Version, m.Name, err)
		}
		if err := m.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate up %s (%s): %w", m.Version, m.Name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`, m.Version, m.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate up %s (%s): track: %w", m.Version, m.Name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migrate up %s (%s): commit: %w", m.Version, m.Name, err)
		}
		fmt.Printf("up   %s %s [%s]\n", m.Version, m.Name, m.DBKey)
		applied++
	}
	if applied == 0 {
		fmt.Println("nothing to migrate")
	}
	return nil
}

// RunDown rolls back the latest `steps` applied migrations, newest first.
func RunDown(ctx context.Context, dbs map[string]*sql.DB, steps int) error {
	if steps < 1 {
		steps = 1
	}

	var applied []Migration
	for _, m := range All() {
		done, err := isApplied(ctx, dbs, m)
		if err != nil {
			return err
		}
		if done {
			applied = append(applied, m)
		}
	}
	sort.Slice(applied, func(i, j int) bool { return applied[i].Version > applied[j].Version })

	if len(applied) == 0 {
		fmt.Println("nothing to roll back")
		return nil
	}
	if steps > len(applied) {
		steps = len(applied)
	}

	for _, m := range applied[:steps] {
		db, err := dbFor(dbs, m.DBKey)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("migrate down %s (%s): begin: %w", m.Version, m.Name, err)
		}
		if err := m.Down(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate down %s (%s): %w", m.Version, m.Name, err)
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM schema_migrations WHERE version = $1`, m.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("migrate down %s (%s): untrack: %w", m.Version, m.Name, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("migrate down %s (%s): commit: %w", m.Version, m.Name, err)
		}
		fmt.Printf("down %s %s [%s]\n", m.Version, m.Name, m.DBKey)
	}
	return nil
}

// RunStatus prints every registered migration and whether it is applied.
func RunStatus(ctx context.Context, dbs map[string]*sql.DB) error {
	for _, m := range All() {
		status := "pending"
		at := "-"
		done, appliedAt, err := appliedAt(ctx, dbs, m)
		if err != nil {
			return err
		}
		if done {
			status = "applied"
			at = appliedAt.Format(time.DateTime)
		}
		fmt.Printf("%-8s %s %-40s %-20s %s\n", status, m.Version, m.Name, m.DBKey, at)
	}
	return nil
}

func dbFor(dbs map[string]*sql.DB, key string) (*sql.DB, error) {
	db, ok := dbs[key]
	if !ok || db == nil {
		return nil, fmt.Errorf("no database connection for key %q", key)
	}
	return db, nil
}

func ensureTrackingTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, createTrackingTable)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	return nil
}

func isApplied(ctx context.Context, dbs map[string]*sql.DB, m Migration) (bool, error) {
	done, _, err := appliedAt(ctx, dbs, m)
	return done, err
}

func appliedAt(ctx context.Context, dbs map[string]*sql.DB, m Migration) (bool, time.Time, error) {
	db, err := dbFor(dbs, m.DBKey)
	if err != nil {
		return false, time.Time{}, err
	}
	if err := ensureTrackingTable(ctx, db); err != nil {
		return false, time.Time{}, err
	}

	var at time.Time
	err = db.QueryRowContext(ctx, `SELECT applied_at FROM schema_migrations WHERE version = $1`, m.Version).Scan(&at)
	if err == sql.ErrNoRows {
		return false, time.Time{}, nil
	}
	if err != nil {
		return false, time.Time{}, fmt.Errorf("check migration %s on %s: %w", m.Version, m.DBKey, err)
	}
	return true, at, nil
}
