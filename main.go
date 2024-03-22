package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/selfphonics/api/internal/handler"
	"github.com/selfphonics/api/internal/middleware"
	"github.com/selfphonics/api/internal/server"
	"github.com/selfphonics/api/internal/storage/memory"
)

//go:embed all:build
var fe embed.FS

func init() {
	build, _ := debug.ReadBuildInfo()

	enableDebug := flag.Bool("debug", false, "whether to enable debug mode")
	flag.Parse()

	opts := &slog.HandlerOptions{}
	if *enableDebug {
		opts.Level = slog.LevelDebug
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, opts)).With(
		slog.Group("program_info",
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", build.GoVersion),
		),
	))
}

func main() {
	k := koanf.New(".")

	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := k.Load(env.Provider("SELFPHONICS_", ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, "SELFPHONICS_")), "_", ".")
	}), nil); err != nil {
		log.Fatalf("error merging environment variables: %v", err)
	}

	serverPort := k.Int("server.port")
	databaseType := k.String("database.type")

	slog.Info("Config",
		slog.Group("server", slog.Int("port", serverPort)),
		slog.Group("database", slog.String("type", databaseType)),
	)

	var srw server.StorageReaderWriter
	switch databaseType {
	case "memory":
		srw = memory.New()
	default:
		log.Fatalf("unsupported databaseType: %s", databaseType)
	}

	s := server.New(srw)
	h := handler.New(s)

	distFS, err := fs.Sub(fe, "build")
	if err != nil {
		log.Fatalf("error getting static files: %+v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(distFS)))
	mux.HandleFunc("GET /api/words", h.ListWords)
	mux.HandleFunc("GET /api/word/{id}", h.GetWordByID)
	mux.HandleFunc("GET /api/word/random", h.GetRandomWord)
	mux.HandleFunc("POST /api/word", h.PostWord)
	wmux := middleware.NewRequestID(middleware.NewLogger(mux))

	slog.Info("Server running", "addr", fmt.Sprintf(":%d", serverPort))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serverPort), wmux); err != nil {
		panic(err)
	}
}
