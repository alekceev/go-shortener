package mem

import (
	"context"
	"database/sql"
	"sync"

	"github.com/alekceev/go-shortener/app/repos/target"
	"github.com/google/uuid"
)

var _ target.URLStore = &MemStore{}

type MemStore struct {
	sync.Mutex
	shortUrlMap map[string]target.URL
	urlMap      map[string]target.URL
}

func NewMemStore() *MemStore {
	return &MemStore{
		shortUrlMap: make(map[string]target.URL),
		urlMap:      make(map[string]target.URL),
	}
}

func (us *MemStore) Create(ctx context.Context, u target.URL) (*uuid.UUID, error) {
	us.Lock()
	defer us.Unlock()

	us.urlMap[u.URL] = u
	us.shortUrlMap[u.ShortURL] = u

	return &u.ID, nil
}

func (us *MemStore) Update(ctx context.Context, u target.URL) error {
	us.Lock()
	defer us.Unlock()

	us.urlMap[u.URL] = u
	us.shortUrlMap[u.ShortURL] = u

	return nil
}

func (us *MemStore) GetURL(ctx context.Context, shortURL string) (*target.URL, error) {
	us.Lock()
	defer us.Unlock()

	if url, found := us.shortUrlMap[shortURL]; found {
		return &url, nil
	}

	return nil, sql.ErrNoRows
}
