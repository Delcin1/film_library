package delete

import (
	"film_library/internal/lib/api/response"
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
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorMovieDeleter
type ActorMovieDeleter interface {
	DeleteActorMovie(movieId int, actorsIds []int) error
}

// @Summary		Delete actors from movie
// @Description	Delete actors from movie by movie_id and actors_ids
// @Tags			Actor-Movie
// @Accept			json
// @Produce		json
// @Param			movie_id	path		int		true	"Movie ID"
// @Param			actors_ids	path		[]int	true	"Actors IDs"
// @Success		200			{object}	Response
// @Failure		400			{object}	response.Response
// @Failure		401			{object}	response.Response
// @Failure		403			{object}	response.Response
// @Router			/actor-movie/delete [delete]
func New(log *slog.Logger, actorMovieDeleter ActorMovieDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor-movie.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", sl.Err(err))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if ok, field, msg := validateRequest(req); !ok {
			log.Error("invalid request", field)

			render.JSON(w, r, response.Error(msg))

			return
		}

		err = actorMovieDeleter.DeleteActorMovie(req.MovieId, req.ActorsIds)
		if err != nil {
			log.Error("failed to delete actor-movie", sl.Err(err))

			render.JSON(w, r, response.Error("failed to delete actor-movie"))

			return
		}

		log.Info("actors deleted from movie", slog.Int("movie_id", req.MovieId))

		render.JSON(w, r, Response{response.OK()})
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
