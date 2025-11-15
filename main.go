package main

import (
	"database/sql"
	"examples/SimpleBankProject/api"
	db "examples/SimpleBankProject/db/sqlc"
	"fmt"
	"log"
	"os"

	"examples/SimpleBankProject/config"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var (
	dbDriver = "postgres"
	//dbSource      = viper.GetString("dbSource")
)

func main() {
	// Example usage of the function
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Hello World!")

	_ = godotenv.Load(".env")

	viper.AutomaticEnv()
	dbSource := viper.GetString("DBSOURCE")
	if dbSource == "" {
		dbSource = "postgres://postgres:kanikaindro@db.vybiuledkxidxrtdxday.supabase.co:5432/postgres"
	}
	//a
	port := os.Getenv("PORT")
	serverAddress := ":9090" // fallback for local dev
	if port != "" {
		serverAddress = ":" + port
	}
	// This is config connection
	configCache := config.SetupRedisCache()

	Client := configCache.Client

	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store, Client)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	server.MountHandlers()

	if err := server.Start(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Println("Server started on", serverAddress)
}
