package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/condemo/nes-cards-backend/api"
	"github.com/condemo/nes-cards-backend/store"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	addr := flag.String("p", ":3000", "addr")
	flag.Parse()

	sqlStorage := store.NewPostgresqlStore()
	db, err := sqlStorage.Init()
	if err != nil {
		log.Fatal("database error: ", err)
	}

	store := store.NewStorage(db)
	apiServer := api.NewApiServer(*addr, store)
	fmt.Println("Server starting at port ->", *addr)
	apiServer.Run()
}
