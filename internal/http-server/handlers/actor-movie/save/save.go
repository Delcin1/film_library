package save

import (
	resp "film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	MovieId   int   `json:"movie_id"`
	ActorsIds []int `json:"actors_ids"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorMovieSaver
type ActorMovieSaver interface {
	SaveActorMovie(movieId int, actorsIds []int) error
}

func New(log *slog.Logger, actorMovieSaver ActorMovieSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor-movie.save.New"

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

		err = actorMovieSaver.SaveActorMovie(req.MovieId, req.ActorsIds)
		if err != nil {
			log.Error("failed to save actor-movie", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save actor-movie"))

			return
		}

		log.Info("actors added to movie", slog.Int("movie_id", req.MovieId))

		render.JSON(w, r, Response{resp.OK()})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.MovieId < 1 {
		return false, slog.String("field", "movie_id"), "field movie_id is not valid"
	}
	for _, id := range req.ActorsIds {
		if id < 1 {
			return false, slog.String("field", "actors_ids"), "field actors_ids is not valid"
		}
	}
	return true, slog.Attr{}, ""
}
