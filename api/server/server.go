package server

import (
	"context"
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
	s.srv.Shutdown(ctx)
	cancel()
}

func (s *Server) Start() {
	go s.srv.ListenAndServe()
}
