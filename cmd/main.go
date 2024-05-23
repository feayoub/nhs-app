package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/adlio/trello"
	"github.com/feayoub/nhs-app/internal/config"
	"github.com/feayoub/nhs-app/internal/handlers"
	"github.com/feayoub/nhs-app/internal/middleware"
	"github.com/feayoub/nhs-app/internal/usecase"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.LoadConfig()

	trelloClient := trello.NewClient(cfg.TrelloAPIKey, cfg.TrelloAPIToken)
	createTrelloCardUseCase := usecase.NewCreateTrelloCardUseCase(trelloClient)

	homeHandler := handlers.NewHomeHandler()
	uploadHandler := handlers.NewUploadHandler(createTrelloCardUseCase)

	router := http.NewServeMux()
	stack := middleware.CreateStack(
		middleware.LoggerMiddleware,
		middleware.CSPMiddleware,
		middleware.TextHTMLMiddleware,
	)

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	router.Handle("/", homeHandler)
	router.Handle("POST /upload", uploadHandler)

	killSig := make(chan os.Signal, 1)

	signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

	server := http.Server{
		Addr:    ":8080",
		Handler: stack(router),
	}

	go func() {
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			logger.Info("Server shutdown complete")
		} else if err != nil {
			logger.Error("Server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	logger.Info("Server started", slog.String("port", "8080"))
	<-killSig

	logger.Info("Shutting down server")

	// Create a context with a timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", slog.Any("err", err))
		os.Exit(1)
	}

	logger.Info("Server shutdown complete")
}
