package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alekceev/go-shortener/app/config"
)

type Server struct {
	srv http.Server
}

func NewServer(conf config.Config, h http.Handler) *Server {
	return &Server{
		srv: http.Server{
			Addr:              conf.Addr(),
			Handler:           h,
			ReadTimeout:       time.Duration(conf.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(conf.WriteTimeout) * time.Second,
			ReadHeaderTimeout: time.Duration(conf.ReadHeaderTimeout) * time.Second,
		},
	}
}

func (s *Server) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	cancel()
}

func (s *Server) Start() {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil {
			log.Panicln(err)
		}
	}()
}
