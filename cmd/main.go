package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/condemo/nes-cards-backend/api"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

func main() {
	var db *bun.DB

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	addr := flag.String("p", ":3000", "addr")
	flag.Parse()

	if os.Getenv("CURRENT_ENV") == "prod" {
		sqlStorage := store.NewPostgresqlStore()
		tempDB, err := sqlStorage.Init()
		if err != nil {
			log.Fatal("database error: ", err)
		}
		db = tempDB
	} else {
		sqlStorage := store.NewSqliteStore()
		tempDB, err := sqlStorage.Init()
		if err != nil {
			log.Fatal("database error: ", err)
		}
		db = tempDB
	}

	store := store.NewStorage(db)
	apiServer := api.NewApiServer(*addr, store)
	fmt.Println("Server starting at port ->", *addr)
	apiServer.Run()
}
