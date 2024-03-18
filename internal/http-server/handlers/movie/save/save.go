package save

import (
	resp "film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate string `json:"release_date"`
	Rating      int    `json:"rating"`
	ActorsIds   []int  `json:"actors_ids"`
}

type Response struct {
	resp.Response
	MovieId int `json:"movie_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=MovieSaver
type MovieSaver interface {
	SaveMovie(title string, description string, releaseDate string, rating int, actorsIds []int) (int, error)
}

func New(log *slog.Logger, movieSaver MovieSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movie.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if ok, field, msg := validateRequest(req); !ok {
			log.Error("invalid request", field)

			render.JSON(w, r, resp.Error(msg))

			return
		}

		movieId, err := movieSaver.SaveMovie(req.Title, req.Description, req.ReleaseDate, req.Rating, req.ActorsIds)
		if err != nil {
			log.Error("failed to save movie", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save movie"))

			return
		}

		log.Info("movie saved", slog.Int("movie_id", movieId))

		render.JSON(w, r, Response{
			resp.OK(),
			movieId,
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if len(req.Title) < 1 || len(req.Title) > 150 {
		return false, slog.String("field", "title"), "field title is not valid"
	}
	if len(req.Description) > 1000 {
		return false, slog.String("field", "description"), "field description is not valid"
	}
	if _, err := time.Parse("2006-01-02", req.ReleaseDate); err != nil {
		return false, slog.String("field", "release_date"), "field release_date is not valid"
	}
	if req.Rating < 0 || req.Rating > 10 {
		return false, slog.String("field", "rating"), "field rating is not valid"
	}
	for _, id := range req.ActorsIds {
		if id < 1 {
			return false, slog.String("field", "actors_ids"), "field actors_ids is not valid"
		}
	}
	return true, slog.Attr{}, ""
}
