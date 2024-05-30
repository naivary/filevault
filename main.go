package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)
const filevault = "filevault"

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Args, os.Getenv, os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func newFlagset(cfg *Config, getenv func(string) string) *flag.FlagSet {
    defaultValue := func(key string) string {
        val := getenv(key)
        if val != "" {
            return val
        }
        switch key {
        case "FILEVAULT_HOST":
            return cfg.Host
        case "FILEVAULT_PORT":
            return cfg.Port
        case "FILEVAULT_DIR":
            return cfg.Dir
        default:
            return ""
        }
    }

	fs := flag.NewFlagSet(filevault, flag.ContinueOnError)
	fs.StringVar(&cfg.Dir, "dir", defaultValue("FILEVAULT_DIR"), "directory which will be used for storing the files")
	fs.StringVar(&cfg.Host, "host", defaultValue("FILEVAULT_HOST"), "host addr on which the http server will run")
	fs.StringVar(&cfg.Port, "port", defaultValue("FILEVAULT_PORT"), "port on which the http server will listen")
	return fs
}


func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	cfg := NewConfig()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	fs := newFlagset(&cfg, getenv)
	fs.SetOutput(stdout)
	fs.Parse(args[1:])

	svc := NewFilevaultService(cfg)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: NewServer(cfg, svc),
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()

	return nil
}
