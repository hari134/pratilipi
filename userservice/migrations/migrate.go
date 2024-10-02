package migrations

import (
	"context"
	"embed"
	"log"
	"os"

	"github.com/hari134/pratilipi/pkg/db" // Import the common db package
	"github.com/uptrace/bun/migrate"
)

//go:embed scripts/*.sql
var sqlMigrations embed.FS

var Migrations = migrate.NewMigrations()

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}
func RunMigrations(dbInstance *db.DB){
	// Initialize the database using the common db package
	defer db.CloseDB(dbInstance) // Ensure the connection is closed when done
	migrator := migrate.NewMigrator(dbInstance, Migrations)
	// Create the migration runner with the migrations folder
	group, err := migrator.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if group.IsZero() {
		log.Printf("there are no new migrations to run (database is up to date)\n")
	}

	// Get the context for running migrations
	ctx := context.Background()

	// Parse the command-line arguments for migration commands
	switch os.Args[1] {
	case "init":
		// Initialize the migrations table
		if err := migrator.Init(ctx); err != nil {
			log.Fatalf("Failed to initialize migrations: %v", err)
		}

	case "up":
		// Apply all the up migrations
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Printf("Applied migrations: %s", group)

	case "down":
		// Roll back the last migration
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Printf("Rolled back migrations: %s", group)

	case "reset":
		// Roll back all migrations and then apply them again
		if err := migrator.Reset(ctx); err != nil {
			log.Fatalf("Failed to reset migrations: %v", err)
		}
		log.Println("Reset migrations")

	default:
		log.Println("Usage: go run main.go [init|up|down|reset]")
	}
}
