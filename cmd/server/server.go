package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/sebastianvera/ghreviews/pkg/database"
	"github.com/sebastianvera/ghreviews/pkg/graph"
)

var (
	listenAddr string
	env        string
)

func main() {
	parseFlags()

	logger := logrus.New()

	store := database.NewStore()

	isProduction := env == "production"
	if !isProduction {
		logger.SetLevel(logrus.DebugLevel)
	}

	r := graph.NewResolver(logger, store)
	graphqlServer := graph.NewServer(r, isProduction)

	router := http.NewServeMux()
	router.Handle("/query", graphqlServer)
	if env != "production" {
		router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	}

	logger.Printf("connect to http://%s/ for GraphQL playground", listenAddr)
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	w := logger.Writer()
	defer w.Close()
	server := &http.Server{
		Addr:         listenAddr,
		Handler:      cors.AllowAll().Handler(router),
		ErrorLog:     log.New(w, "", 0),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		<-quit
		logger.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	logger.Info("Server is ready to handle requests at ", listenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("Could not listen on %s: %v\n", listenAddr, err)
	}

	<-done
	logger.Info("Server stopped")
}

func parseFlags() {
	flag.StringVar(&listenAddr, "listen-addr", "localhost:8080", "server listen address")
	flag.StringVar(&env, "env", "development", "application environment (development, production)")
	flag.Parse()
}
