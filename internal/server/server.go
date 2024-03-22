package server

import (
	"context"
	"log/slog"

	"github.com/selfphonics/api/internal/middleware"
	"github.com/selfphonics/api/internal/storage"
)

type StorageReader interface {
	GetWordByID(ctx context.Context, id string) (*storage.Word, error)
	GetRandomWord(ctx context.Context) (*storage.Word, error)
	ListWords(ctx context.Context) ([]storage.Word, error)
}

type StorageWriter interface {
	AddWord(ctx context.Context, word storage.Word) (*storage.Word, error)
}

type StorageReaderWriter interface {
	StorageReader
	StorageWriter
}

type Server struct {
	storageReader StorageReader
	storageWriter StorageWriter
}

type Word struct {
	ID       string                   `json:"id,omitempty"`
	Word     string                   `json:"word,omitempty"`
	Sections []map[string]interface{} `json:"sections,omitempty"`
}

func New(srw StorageReaderWriter) *Server {
	return &Server{
		storageReader: srw,
		storageWriter: srw,
	}
}

func (s *Server) ListWords(ctx context.Context) ([]Word, error) {
	res, err := s.storageReader.ListWords(ctx)
	if err != nil {
		return nil, err
	}

	words := make([]Word, len(res))

	for i, w := range res {
		words[i] = Word{ID: w.ID, Word: w.Word, Sections: w.Sections}
	}

	return words, nil
}

func (s *Server) GetWordByID(ctx context.Context, id string) (*Word, error) {
	slog.Info("GetWordByID", "requestID", ctx.Value(middleware.ContextKeyRequestID))

	res, err := s.storageReader.GetWordByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Word{Word: res.Word, Sections: res.Sections}, nil
}

func (s *Server) GetRandomWord(ctx context.Context) (*Word, error) {
	res, err := s.storageReader.GetRandomWord(ctx)
	if err != nil {
		return nil, err
	}

	return &Word{Word: res.Word, Sections: res.Sections}, nil
}

func (s *Server) AddWord(ctx context.Context, w Word) (*Word, error) {
	slog.Info("GetWordByID", "requestID", ctx.Value(middleware.ContextKeyRequestID))

	in := storage.Word{Word: w.Word, Sections: w.Sections}

	res, err := s.storageWriter.AddWord(ctx, in)
	if err != nil {
		return nil, err
	}

	return &Word{ID: res.ID, Word: res.Word, Sections: res.Sections}, nil
}
