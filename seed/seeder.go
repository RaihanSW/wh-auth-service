package seed

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
)

// Seeder inserts initial/demo data into one database key from
// application_<env>.json. Register new seeders from init() in
// files named <name>_seeder.go. Seeders must be idempotent —
// safe to run repeatedly without duplicating rows.
type Seeder struct {
	Name  string // e.g. "users"
	DBKey string // cfg.Database key — always the writer, never the replica
	Run   func(ctx context.Context, db *sql.DB) error
}

var registry []Seeder

func Register(s Seeder) {
	registry = append(registry, s)
}

// All returns registered seeders sorted by name.
func All() []Seeder {
	ss := make([]Seeder, len(registry))
	copy(ss, registry)
	sort.Slice(ss, func(i, j int) bool { return ss[i].Name < ss[j].Name })
	return ss
}

// Run executes seeders by name. Empty or "all" runs every registered seeder.
func Run(ctx context.Context, dbs map[string]*sql.DB, name string) error {
	if name == "" || name == "all" {
		for _, s := range All() {
			if err := runOne(ctx, dbs, s); err != nil {
				return err
			}
		}
		return nil
	}

	for _, s := range All() {
		if s.Name == name {
			return runOne(ctx, dbs, s)
		}
	}
	return fmt.Errorf("unknown seeder %q (registered: %v)", name, names())
}

func names() []string {
	out := make([]string, 0, len(registry))
	for _, s := range All() {
		out = append(out, s.Name)
	}
	return out
}

func runOne(ctx context.Context, dbs map[string]*sql.DB, s Seeder) error {
	db, ok := dbs[s.DBKey]
	if !ok || db == nil {
		return fmt.Errorf("no database connection for key %q", s.DBKey)
	}
	if err := s.Run(ctx, db); err != nil {
		return fmt.Errorf("seed %s [%s]: %w", s.Name, s.DBKey, err)
	}
	fmt.Printf("seeded %s [%s]\n", s.Name, s.DBKey)
	return nil
}
