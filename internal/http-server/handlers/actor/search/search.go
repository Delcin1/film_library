package search

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
	ActorId int `json:"actor_id"`
}

type Response struct {
	response.Response
	Actor postgres.Actor `json:"actor"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorSearcher
type ActorSearcher interface {
	GetActor(actorId int) (postgres.Actor, error)
}

//	@Summary		Get an actor
//	@Description	Get an actor by actor_id
//	@Tags			Actor
//	@Accept			json
//	@Produce		json
//	@Param			actor_id	path		int	true	"Actor ID"
//	@Success		200			{object}	Response
//	@Failure		400			{object}	response.Response
//	@Failure		401			{object}	response.Response
//	@Router			/actor/search [get]
func New(log *slog.Logger, actorSearcher ActorSearcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor.search.New"

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

		if req.ActorId < 1 {
			log.Error("invalid actor_id", slog.Int("actor_id", req.ActorId))

			render.JSON(w, r, response.Error("field actor_id is not valid"))

			return
		}

		actor, err := actorSearcher.GetActor(req.ActorId)
		if err != nil {
			log.Error("actor search failed", sl.Err(err))

			render.JSON(w, r, response.Error("actor search failed"))

			return
		}

		log.Info("actor found", slog.Int("actor_id", req.ActorId))

		render.JSON(w, r, Response{
			response.OK(),
			actor,
		})
	}
}
