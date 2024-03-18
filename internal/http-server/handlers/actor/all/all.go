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

type Response struct {
	response.Response
	Movies []postgres.Actor `json:"actors"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorsAllGetter
type ActorsAllGetter interface {
	GetActors() ([]postgres.Actor, error)
}

//	@Summary		Get all actors
//	@Description	Get all actors
//	@Tags			Actor
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Response
//	@Failure		400	{object}	response.Response
//	@Failure		401	{object}	response.Response
//	@Router			/actor/all [get]
func New(log *slog.Logger, actorsAllGetter ActorsAllGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor.all.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		actors, err := actorsAllGetter.GetActors()
		if err != nil {
			log.Error("actors search failed", sl.Err(err))

			render.JSON(w, r, response.Error("actors search failed"))

			return
		}

		log.Info("actors found", slog.Int("actors_count", len(actors)))

		render.JSON(w, r, Response{
			response.OK(),
			actors,
		})
	}
}
