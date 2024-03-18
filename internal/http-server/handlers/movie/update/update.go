package update

import (
	"film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type Request struct {
	MovieId     int     `json:"movie_id"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	ReleaseDate *string `json:"release_date,omitempty"`
	Rating      *int    `json:"rating,omitempty"`
}

type Response struct {
	response.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=MovieUpdater
type MovieUpdater interface {
	UpdateMovieTitle(movieId int, title string) error
	UpdateMovieDescription(movieId int, description string) error
	UpdateMovieReleaseDate(movieId int, releaseDate string) error
	UpdateMovieRating(movieId int, rating int) error
}

// @Summary		Update movie
// @Description	Update movie by movie_id
// @Tags			Movie
// @Accept			json
// @Produce		json
// @Param			movie_id		path		int		true	"Movie ID"
// @Param			title			body		string	false	"Title"
// @Param			description		body		string	false	"Description"
// @Param			release_date	body		string	false	"Release Date"
// @Param			rating			body		int		false	"Rating"
// @Success		200				{object}	Response
// @Failure		400				{object}	response.Response
// @Failure		401				{object}	response.Response
// @Failure		403				{object}	response.Response
// @Router			/movie/update [post]
func New(log *slog.Logger, movieSaver MovieUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.movie.update.New"

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

		if req.Title != nil {
			err := movieSaver.UpdateMovieTitle(req.MovieId, *req.Title)
			if err != nil {
				log.Error("failed to update movie title", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update movie title"))

				return
			}
		}

		if req.Description != nil {
			err := movieSaver.UpdateMovieDescription(req.MovieId, *req.Description)
			if err != nil {
				log.Error("failed to update movie description", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update movie description"))

				return
			}
		}

		if req.ReleaseDate != nil {
			err := movieSaver.UpdateMovieReleaseDate(req.MovieId, *req.ReleaseDate)
			if err != nil {
				log.Error("failed to update movie release date", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update movie release date"))

				return
			}
		}

		if req.Rating != nil {
			err := movieSaver.UpdateMovieRating(req.MovieId, *req.Rating)
			if err != nil {
				log.Error("failed to update movie rating", sl.Err(err))

				render.JSON(w, r, response.Error("failed to update movie rating"))

				return
			}
		}

		if req.Title == nil && req.Description == nil && req.ReleaseDate == nil && req.Rating == nil {
			log.Error("no fields to update")

			render.JSON(w, r, response.Error("no fields to update"))

			return
		}

		log.Info("movie updated", slog.Int("movie_id", req.MovieId))

		render.JSON(w, r, Response{response.OK()})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.MovieId < 1 {
		return false, slog.String("field", "movie_id"), "field movie_id is not valid"
	}
	if req.Title != nil && (len(*req.Title) < 1 || len(*req.Title) > 150) {
		return false, slog.String("field", "title"), "field title is not valid"
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		return false, slog.String("field", "description"), "field description is not valid"
	}
	if req.ReleaseDate != nil {
		if _, err := time.Parse("2006-01-02", *req.ReleaseDate); err != nil {
			return false, slog.String("field", "release_date"), "field release_date is not valid"
		}
	}
	if req.Rating != nil && (*req.Rating < 0 || *req.Rating > 10) {
		return false, slog.String("field", "rating"), "field rating is not valid"
	}
	return true, slog.Attr{}, ""
}
