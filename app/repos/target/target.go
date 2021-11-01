package target

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/alekceev/go-shortener/app/decoder"
	"github.com/google/uuid"
)

type URL struct {
	ID           uuid.UUID
	URL          string
	ShortURL     string
	NumRedirects int
}

// URLStore interface for storing and getting url data.
type URLStore interface {
	Create(ctx context.Context, url URL) (*uuid.UUID, error)
	Update(ctx context.Context, url URL) error
	GetURL(ctx context.Context, shortURL string) (*URL, error)
}

type Urls struct {
	store URLStore
}

func NewUrls(store URLStore) *Urls {
	return &Urls{
		store: store,
	}
}

func (u *Urls) Create(ctx context.Context, url URL) (*URL, error) {
	url.ID = uuid.New()
	url.ShortURL, _ = decoder.Decode(url.ID)
	id, err := u.store.Create(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error when creating: %w", err)
	}
	url.ID = *id
	return &url, nil
}

func (u *Urls) Update(ctx context.Context, url URL) error {
	// increment stat
	url.NumRedirects = url.NumRedirects + 1

	err := u.store.Update(ctx, url)
	if err != nil {
		return fmt.Errorf("error when updating: %w", err)
	}
	return nil
}

func (u *Urls) GetURL(ctx context.Context, shortURL string) (*URL, error) {
	url, err := u.store.GetURL(ctx, shortURL)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("URL not found")
		default:
			return nil, fmt.Errorf("error when getting url: %w", err)
		}
	}

	return url, nil
}
