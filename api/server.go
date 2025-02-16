package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/condemo/nes-cards-backend/api/handlers"
	"github.com/condemo/nes-cards-backend/api/middlewares"
	"github.com/condemo/nes-cards-backend/store"
)

type ApiServer struct {
	addr  string
	store store.Store
}

func NewApiServer(addr string, s store.Store) *ApiServer {
	return &ApiServer{
		addr:  addr,
		store: s,
	}
}

func (s *ApiServer) Run() {
	router := http.NewServeMux()
	api := http.NewServeMux()
	game := http.NewServeMux()
	player := http.NewServeMux()

	basicMiddlewares := middlewares.MiddlewareStack(
		middlewares.AddCors,
		middlewares.Recover,
		middlewares.Logger,
	)

	router.Handle("/api/v1/", http.StripPrefix("/api/v1", basicMiddlewares(api)))
	api.Handle("/game/", http.StripPrefix("/game", game))
	api.Handle("/player/", http.StripPrefix("/player", player))

	gameHandler := handlers.NewGameHandler(s.store)
	gameHandler.RegisterRoutes(game)

	playerHandler := handlers.NewPlayerHandler(s.store)
	playerHandler.RegisterRoutes(player)

	server := http.Server{
		Addr:         s.addr,
		Handler:      router,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
	}

	go func() {
		log.Fatal(server.ListenAndServe())
	}()

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigC

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	server.Shutdown(ctx)
}
