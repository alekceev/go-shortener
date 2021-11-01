package handler

import (
	"context"
	"fmt"

	"github.com/alekceev/go-shortener/app/repos/target"
	"github.com/google/uuid"
)

type Handlers struct {
	target *target.Urls
}

func NewHandlers(target *target.Urls) *Handlers {
	h := &Handlers{
		target: target,
	}
	return h
}

type URL struct {
	ID           uuid.UUID `json:"id"`
	URL          string    `json:"url"`
	ShortURL     string    `json:"short_url"`
	NumRedirects int       `json:"num_redirects"`
}

func (rt *Handlers) Create(ctx context.Context, ur URL) (URL, error) {
	u := target.URL{
		URL: ur.URL,
	}
	nu, err := rt.target.Create(ctx, u)
	if err != nil {
		return URL{}, fmt.Errorf("error when creating: %w", err)
	}

	return URL{
		ID:           nu.ID,
		URL:          nu.URL,
		ShortURL:     nu.ShortURL,
		NumRedirects: nu.NumRedirects,
	}, nil
}

func (rt *Handlers) GetURL(ctx context.Context, shortURL string) (URL, error) {
	u, err := rt.target.GetURL(ctx, shortURL)
	if err != nil {
		return URL{}, fmt.Errorf("error when get shor url: %w", err)
	}

	// increment NumRedirects
	if err = rt.target.Update(ctx, *u); err != nil {
		return URL{}, fmt.Errorf("error when get shor url: %w", err)
	}

	return URL{
		ID:           u.ID,
		URL:          u.URL,
		ShortURL:     u.ShortURL,
		NumRedirects: u.NumRedirects,
	}, nil
}

func (rt *Handlers) GetStats(ctx context.Context, shortURL string) (URL, error) {
	u, err := rt.target.GetURL(ctx, shortURL)
	if err != nil {
		return URL{}, fmt.Errorf("error when get stats: %w", err)
	}

	return URL{
		ID:           u.ID,
		URL:          u.URL,
		ShortURL:     u.ShortURL,
		NumRedirects: u.NumRedirects,
	}, nil
}
