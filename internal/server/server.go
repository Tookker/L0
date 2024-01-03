package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"L0/internal/config"
	"L0/internal/router"
	"L0/internal/store"
)

type Server struct {
	router *router.ChiRouter
	store  store.Store
	ip     string
}

func NewServer(router *router.ChiRouter, store store.Store, config *config.Config) *Server {
	return &Server{
		router: router,
		ip:     config.ServerIp,
		store:  store,
	}
}

func (s *Server) StartServer() error {
	s.router.SetHandlers()
	srv := &http.Server{
		Addr:    s.ip,
		Handler: s.router.GetRouter(),
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}

		close(idleConnsClosed)
	}()

	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	<-idleConnsClosed

	return nil
}
