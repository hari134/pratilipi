

package migrations

import (
	"context"
	"embed"
	"log"

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
func RunMigrations(dbInstance *db.DB) {
	// Initialize the database using the common db package
	migrator := migrate.NewMigrator(dbInstance, Migrations)

	if err := migrator.Init(context.Background()); err != nil {
		log.Fatalf("Failed to initialize migrations: %v", err)
	}
	// Create the migration runner with the migrations folder
	group, err := migrator.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if group.IsZero() {
		log.Printf("there are no new migrations to run (database is up to date)\n")
	}

}
