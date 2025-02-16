package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/condemo/nes-cards-backend/api"
	"github.com/condemo/nes-cards-backend/store"
)

func main() {
	addr := flag.String("p", ":3000", "addr")
	flag.Parse()

	sqliteStorage := store.NewSqliteStore()
	db, err := sqliteStorage.Init()
	if err != nil {
		log.Fatal("database error: ", err)
	}

	store := store.NewStorage(db)
	apiServer := api.NewApiServer(*addr, store)
	fmt.Println("Server starting at port ->", *addr)
	apiServer.Run()
}
