package search_by_id

import (
	"film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"film_library/internal/storage/postgres"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	MovieId int `json:"movie_id"`
}

type Response struct {
	response.Response
	Movie postgres.Movie `json:"movie"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=MovieSearcherById
type MovieSearcherById interface {
	GetMovie(movieId int) (postgres.Movie, error)
}

// @Summary		Search a movie by movie_id
// @Description	Search a movie by movie_id
// @Tags			Movie
// @Accept			json
// @Produce		json
// @Param			movie_id	path		int	true	"Movie ID"
// @Success		200			{object}	Response
// @Failure		400			{object}	response.Response
// @Failure		401			{object}	response.Response
// @Router			/movie/search_by_id [get]
func New(log *slog.Logger, movieSearcher MovieSearcherById) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movie.search_by_id.New"

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

		if req.MovieId < 1 {
			log.Error("invalid movie_id", slog.Int("movie_id", req.MovieId))

			render.JSON(w, r, response.Error("field movie_id is not valid"))

			return
		}

		movie, err := movieSearcher.GetMovie(req.MovieId)
		if err != nil {
			log.Error("movie search failed", sl.Err(err))

			render.JSON(w, r, response.Error("movie search failed"))

			return
		}

		log.Info("movie found", slog.Int("movie_id", req.MovieId))

		render.JSON(w, r, Response{
			response.OK(),
			movie,
		})
	}
}
