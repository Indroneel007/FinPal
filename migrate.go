package main

import (
    "log"
    "os"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    dbURL := "postgresql://finpal_postgres_user:wMJTxyATm6dtr2NGq29Vm7Eala082iEZ@dpg-d27efo6uk2gs73e30sh0-a/finpal_postgres?sslmode=disable"
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
