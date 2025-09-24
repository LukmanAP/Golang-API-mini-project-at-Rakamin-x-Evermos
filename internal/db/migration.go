package db

import (
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "time"

    "gorm.io/gorm"
)

// RunMigrations applies all *.up.sql files in the given directory in alphanumeric order.
// It records applied migrations in the schema_migrations table.
func RunMigrations(gdb *gorm.DB, dir string) error {
    // Ensure schema_migrations table exists
    if err := ensureSchemaMigrationsTable(gdb); err != nil {
        return err
    }

    applied, err := getAppliedMigrations(gdb)
    if err != nil {
        return err
    }

    entries, err := os.ReadDir(dir)
    if err != nil {
        // If the directory does not exist or is empty, treat as no migrations to run
        if errorsIs(err, os.ErrNotExist) {
            return nil
        }
        return err
    }

    // Collect .up.sql files
    var ups []string
    for _, e := range entries {
        if e.IsDir() {
            continue
        }
        name := e.Name()
        if filepath.Ext(name) == ".sql" && hasSuffix(name, ".up.sql") {
            ups = append(ups, filepath.Join(dir, name))
        }
    }

    sort.Strings(ups)

    for _, path := range ups {
        base := filepath.Base(path)
        if applied[base] {
            continue // already applied
        }
        // Read SQL content
        sqlBytes, err := os.ReadFile(path)
        if err != nil {
            return fmt.Errorf("read migration %s: %w", base, err)
        }
        // Execute
        if err := gdb.Exec(string(sqlBytes)).Error; err != nil {
            return fmt.Errorf("apply migration %s: %w", base, err)
        }
        // Record
        if err := recordApplied(gdb, base); err != nil {
            return fmt.Errorf("record migration %s: %w", base, err)
        }
    }

    return nil
}

func ensureSchemaMigrationsTable(gdb *gorm.DB) error {
    const ddl = `
CREATE TABLE IF NOT EXISTS schema_migrations (
  name VARCHAR(255) PRIMARY KEY,
  applied_at DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
    return gdb.Exec(ddl).Error
}

func getAppliedMigrations(gdb *gorm.DB) (map[string]bool, error) {
    type row struct{ Name string }
    var rows []row
    if err := gdb.Raw("SELECT name FROM schema_migrations").Scan(&rows).Error; err != nil {
        // If table doesn't exist yet, return empty set
        // but ensureSchemaMigrationsTable should have created it already
        return map[string]bool{}, nil
    }
    m := make(map[string]bool, len(rows))
    for _, r := range rows {
        m[r.Name] = true
    }
    return m, nil
}

func recordApplied(gdb *gorm.DB, name string) error {
    now := time.Now().Format("2006-01-02 15:04:05")
    return gdb.Exec("INSERT INTO schema_migrations (name, applied_at) VALUES (?, ?)", name, now).Error
}

// helpers (avoid importing strings for tiny helpers)
func hasSuffix(s, suf string) bool {
    if len(s) < len(suf) {
        return false
    }
    return s[len(s)-len(suf):] == suf
}

// errorsIs is a tiny wrapper to avoid importing errors for a single call
func errorsIs(err error, target error) bool {
    type is interface{ Is(error) bool }
    if err == nil {
        return target == nil
    }
    if x, ok := err.(is); ok {
        return x.Is(target)
    }
    return err == target
}