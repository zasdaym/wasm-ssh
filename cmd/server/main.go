package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

//go:embed static/*
var staticFS embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := flag.String("addr", ":8000", "listen address")
	flag.Parse()

	rootFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		return err
	}

	handler := http.FileServerFS(rootFS)
	srv := &http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	var g errgroup.Group
	g.Go(func() error {
		<-ctx.Done()
		return srv.Shutdown(context.Background())
	})
	g.Go(func() error {
		log.Printf("Listening on %s", *addr)
		return srv.ListenAndServe()
	})
	err = g.Wait()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}
