package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := "postgres://Indroneel007:Kanika%402602@db.vybiuledkxidxrtdxday.supabase.co:5432/postgres"
	m, err := migrate.New(
		"file://db/migration",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Migration init failed: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Println("Migration ran successfully!")
	os.Exit(0)
}
