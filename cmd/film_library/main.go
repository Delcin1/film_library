package main

import (
	"film_library/internal/config"
	deleteActorMovie "film_library/internal/http-server/handlers/actor-movie/delete"
	saveActorMovie "film_library/internal/http-server/handlers/actor-movie/save"
	allActors "film_library/internal/http-server/handlers/actor/all"
	deleteActor "film_library/internal/http-server/handlers/actor/delete"
	saveActor "film_library/internal/http-server/handlers/actor/save"
	searchActor "film_library/internal/http-server/handlers/actor/search"
	updateActor "film_library/internal/http-server/handlers/actor/update"
	allMovies "film_library/internal/http-server/handlers/movie/all"
	deleteMovie "film_library/internal/http-server/handlers/movie/delete"
	saveMovie "film_library/internal/http-server/handlers/movie/save"
	searchMovieById "film_library/internal/http-server/handlers/movie/search_by_id"
	searchMovieByPart "film_library/internal/http-server/handlers/movie/search_by_part"
	updateMovie "film_library/internal/http-server/handlers/movie/update"
	mwLogger "film_library/internal/http-server/middleware/logger"
	"film_library/internal/lib/logger/sl"
	"film_library/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting film_library api", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/actor/save", saveActor.New(log, storage))
	router.Post("/movie/save", saveMovie.New(log, storage))
	router.Post("/actor/update", updateActor.New(log, storage))
	router.Post("/movie/update", updateMovie.New(log, storage))
	router.Post("/actor/delete", deleteActor.New(log, storage))
	router.Post("/movie/delete", deleteMovie.New(log, storage))
	router.Post("/actor-movie/save", saveActorMovie.New(log, storage))
	router.Post("/actor-movie/delete", deleteActorMovie.New(log, storage))
	router.Get("/actor/search", searchActor.New(log, storage))
	router.Get("/movie/search_by_id", searchMovieById.New(log, storage))
	router.Get("/movie/all", allMovies.New(log, storage))
	router.Get("/actor/all", allActors.New(log, storage))
	router.Get("/movie/search_by_part", searchMovieByPart.New(log, storage))

	log.Info("starting http server", slog.String("address", cfg.HTTPServer.Address))
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start http server", sl.Err(err))
	}

	log.Error("http server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
