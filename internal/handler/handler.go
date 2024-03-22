package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/selfphonics/api/internal/server"
)

type WordReader interface {
	GetWordByID(ctx context.Context, id string) (*server.Word, error)
	GetRandomWord(ctx context.Context) (*server.Word, error)
	ListWords(ctx context.Context) ([]server.Word, error)
}

type WordWriter interface {
	AddWord(ctx context.Context, data server.Word) (*server.Word, error)
}

type WordReaderWriter interface {
	WordReader
	WordWriter
}

type Handler struct {
	wordReader WordReader
	wordWriter WordWriter
}

func New(wrw WordReaderWriter) *Handler {
	return &Handler{
		wordReader: wrw,
		wordWriter: wrw,
	}
}

func (h *Handler) ListWords(w http.ResponseWriter, r *http.Request) {
	out, err := h.wordReader.ListWords(r.Context())
	if err != nil {
		JSONError(w, map[string]string{"error": "unable to list words", "details": err.Error()}, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		JSONError(w, map[string]string{"error": "failed to encode response", "details": err.Error()}, http.StatusInternalServerError)
	}
}

func (h *Handler) GetWordByID(w http.ResponseWriter, r *http.Request) {
	out, err := h.wordReader.GetWordByID(r.Context(), r.PathValue("id"))
	if err != nil {
		JSONError(w, map[string]string{"error": "unable to get word", "details": err.Error()}, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		JSONError(w, map[string]string{"error": "failed to encode response", "details": err.Error()}, http.StatusInternalServerError)
	}
}

func (h *Handler) GetRandomWord(w http.ResponseWriter, r *http.Request) {
	out, err := h.wordReader.GetRandomWord(r.Context())
	if err != nil {
		JSONError(w, map[string]string{"error": "unable to get random word", "details": err.Error()}, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		JSONError(w, map[string]string{"error": "failed to encode response", "details": err.Error()}, http.StatusInternalServerError)
	}
}

func (h *Handler) PostWord(w http.ResponseWriter, r *http.Request) {
	var in server.Word

	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		JSONError(w, map[string]string{"error": "invalid json", "details": err.Error()}, http.StatusBadRequest)
		return
	}

	out, err := h.wordWriter.AddWord(r.Context(), in)
	if err != nil {
		JSONError(w, map[string]string{"error": "unable to add word", "details": err.Error()}, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(out); err != nil {
		JSONError(w, map[string]string{"error": "failed to encode response", "details": err.Error()}, http.StatusInternalServerError)
	}
}

func JSONError(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
}
