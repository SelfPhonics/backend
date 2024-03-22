package memory

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/google/uuid"

	"github.com/selfphonics/api/internal/storage"
)

type Client struct {
	mux   sync.RWMutex
	words map[string]*storage.Word
}

func New() *Client {
	return &Client{
		words: make(map[string]*storage.Word),
	}
}

func (c *Client) ListWords(ctx context.Context) ([]storage.Word, error) {
	words := make([]storage.Word, 0)

	c.mux.RLock()

	for _, w := range c.words {
		words = append(words, *w)
	}

	c.mux.RUnlock()

	return words, nil
}

func (c *Client) GetWordByID(ctx context.Context, id string) (*storage.Word, error) {
	c.mux.RLock()
	w, ok := c.words[id]
	c.mux.RUnlock()

	if !ok {
		return nil, fmt.Errorf(storage.ErrNotFoundFmt, "id")
	}

	return w, nil
}

func (c *Client) GetRandomWord(ctx context.Context) (*storage.Word, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	if len(c.words) == 0 {
		return nil, storage.ErrNoRecords
	}

	k := rand.Intn(len(c.words))
	for _, w := range c.words {
		if k == 0 {
			return w, nil
		}
		k--
	}

	return &storage.Word{}, nil
}

func (c *Client) AddWord(ctx context.Context, word storage.Word) (*storage.Word, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.words[word.Word]; ok {
		return nil, errors.New("word already exists")
	}

	c.words[word.Word] = &word
	word.ID = uuid.NewString()

	return &word, nil
}
