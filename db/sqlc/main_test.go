package db

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var testQueries *Queries
var testDB *sql.DB

// This is the test connection
func waitForDB(dbSource string) error {
	var db *sql.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sql.Open("postgres", dbSource)
		if err == nil {
			err = db.Ping()
			if err == nil {
				return nil // success
			}
		}
		time.Sleep(2 * time.Second)
	}
	return err
}

const (
	dbDriver = "postgres"
)

func TestMain(m *testing.M) {

	/*err = godotenv.Load()
	if err != nil {
		log.Fatal("Error locacaacading .env file:", err)
	}*/

	paths := []string{".env", "../.env", "../../.env"}
	for _, path := range paths {
		_ = godotenv.Load(path)
	}

	var err error
	viper.AutomaticEnv()

	dbSource := viper.GetString("DBSOURCE")

	if dbSource == "" {
		log.Fatal("Unable to read DBSOURCE environment variable in test")
	}

	if err = waitForDB(dbSource); err != nil {
		log.Fatalf("Cannot connect to DB after waiting: %v", err)
	}

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	defer testDB.Close()
	os.Exit(m.Run())
}
