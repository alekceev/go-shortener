package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/alekceev/go-shortener/api/handler"
	"github.com/alekceev/go-shortener/api/router"
	"github.com/alekceev/go-shortener/api/server"
	"github.com/alekceev/go-shortener/app/config"
	"github.com/alekceev/go-shortener/app/repos/target"
	"github.com/alekceev/go-shortener/app/starter"
	"github.com/alekceev/go-shortener/db/mem"
)

func main() {
	conf, err := config.Get()
	if err != nil {
		log.Fatalf("error parsing config: %v\n", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	// var store url.URLStore
	store := mem.NewMemStore()

	urls := target.NewUrls(store)
	app := starter.NewApp(urls)
	h := handler.NewHandlers(urls)
	rt := router.NewRouter(h)
	srv := server.NewServer(conf, rt)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
