package all

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
	SortBy string `json:"sort_by"`
}

type Response struct {
	response.Response
	Movies []postgres.Movie `json:"movies"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=MoviesAllGetter
type MoviesAllGetter interface {
	GetMovies(sortBy string) ([]postgres.Movie, error)
}

//	@Summary		Get all movies
//	@Description	Get all movies
//	@Tags			Movie
//	@Accept			json
//	@Produce		json
//	@Param			sort_by	query		string	true	"Sort by"
//	@Success		200		{object}	Response
//	@Failure		400		{object}	response.Response
//	@Failure		401		{object}	response.Response
//	@Router			/movie/all [get]
func New(log *slog.Logger, moviesAllGetter MoviesAllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movie.all.New"

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

		if !validateSortBy(req.SortBy) {
			log.Error("invalid sort_by", slog.String("sort_by", req.SortBy))

			render.JSON(w, r, response.Error("field sort_by is not valid"))

			return
		}

		movies, err := moviesAllGetter.GetMovies(req.SortBy)
		if err != nil {
			log.Error("movies search failed", sl.Err(err))

			render.JSON(w, r, response.Error("movies search failed"))

			return
		}

		log.Info("movies found", slog.Int("movies_count", len(movies)))

		render.JSON(w, r, Response{
			response.OK(),
			movies,
		})
	}
}

func validateSortBy(sortBy string) bool {
	return sortBy == postgres.OrderByTitleAsc || sortBy == postgres.OrderByTitleDesc ||
		sortBy == postgres.OrderByReleaseDateAsc || sortBy == postgres.OrderByReleaseDateDesc ||
		sortBy == postgres.OrderByRatingAsc || sortBy == postgres.OrderByRatingDesc
}
