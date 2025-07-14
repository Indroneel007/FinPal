package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var testQueries *Queries
var testDB *sql.DB

const (
	dbDriver = "postgres"
	//dbSource = "postgresql://root:rootpassword@localhost:5433/simple_bank_2?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	err = godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	viper.AutomaticEnv()

	dbSource := viper.GetString("dbSource")
	if dbSource == "" {
		log.Fatal("dbSource is not set in the environment variables")
	}

	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testQueries = New(testDB)
	defer testDB.Close()
	os.Exit(m.Run())
}
