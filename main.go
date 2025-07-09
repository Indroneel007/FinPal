package main

import (
	"database/sql"
	"examples/SimpleBankProject/api"
	db "examples/SimpleBankProject/db/sqlc"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	dbDriver = "postgres"
	//dbSource      = viper.GetString("dbSource")
	//serverAddress = "0.0.0.0:9090" // Change this to your desired address and port
)

func main() {
	// Example usage of the function
	fmt.Println("Hello World!")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	viper.AutomaticEnv()
	dbSource := viper.GetString("dbSource")
	if dbSource == "" {
		log.Fatal("dbSource is not set in the environment variables")
	}

	serverAddress := viper.GetString("serverAddress")
	if serverAddress == "" {
		log.Fatal("serverAddress is not set in the environment variables")
	}

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	server.MountHandlers()

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Println("Server started on", serverAddress)
}
