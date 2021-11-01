package starter

import (
	"context"
	"sync"

	"github.com/alekceev/go-shortener/app/repos/target"
)

type App struct {
	urls *target.Urls
}

func NewApp(urls *target.Urls) *App {
	a := &App{
		urls: urls,
	}
	return a
}

type APIServer interface {
	Start()
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start()
	<-ctx.Done()
	hs.Stop()
}
