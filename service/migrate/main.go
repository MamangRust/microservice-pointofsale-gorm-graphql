package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/MamangRust/microservice-point-of-sale-pkg/dotenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

const (
	dialect = "pgx"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "", "directory with migration files (default: collect all service/*/database/migration)")
)

func main() {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	err := dotenv.Viper()
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
	)

	db, err := goose.OpenDBWithDriver(dialect, connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()

	if *dir != "" {
		// Explicit dir (Docker image legacy path): delegate to goose directly.
		if err := goose.RunContext(context.Background(), command, db, *dir, args[1:]...); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		return
	}

	// F1: per-service migrations. The POS database is a single shared
	// instance, so every service/*/database/migration/*.sql file applies to
	// the same DB. File timestamps are globally unique and increasing, so
	// collecting all files and running them in timestamp order preserves the
	// original chronology (per-dir runs would break cross-service FK deps,
	// e.g. auth.refresh_tokens -> users).
	matches, err := filepath.Glob("service/*/database/migration/*.sql")
	if err != nil {
		log.Fatalf("Failed to glob migration files: %v", err)
	}
	sort.Strings(matches)
	if len(matches) == 0 {
		log.Fatalf("No migration files found under service/*/database/migration (use -dir to set one)")
	}

	// Stage the files into a temp dir (goose requires a plain directory) and
	// run them in one pass.
	tmp, err := os.MkdirTemp("", "pos-migrations-")
	if err != nil {
		log.Fatalf("Failed to create temp migrations dir: %v", err)
	}
	defer os.RemoveAll(tmp)
	for _, m := range matches {
		data, err := os.ReadFile(m)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", m, err)
		}
		if err := os.WriteFile(filepath.Join(tmp, filepath.Base(m)), data, 0o644); err != nil {
			log.Fatalf("Failed to stage %s: %v", filepath.Base(m), err)
		}
	}

	log.Printf("Migrating %s (%d files in %s)", command, len(matches), tmp)
	if err := goose.RunContext(context.Background(), command, db, tmp, args[1:]...); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
