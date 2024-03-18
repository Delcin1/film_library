package delete

import (
	resp "film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	ActorId int `json:"actor_id"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorDeleter
type ActorDeleter interface {
	DeleteActor(actorId int) error
}

func New(log *slog.Logger, actorDeleter ActorDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor.delete.New"

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

		if req.ActorId < 1 {
			log.Error("invalid actor_id", slog.Int("actor_id", req.ActorId))

			render.JSON(w, r, resp.Error("field actor_id is not valid"))

			return
		}

		err = actorDeleter.DeleteActor(req.ActorId)
		if err != nil {
			log.Error("failed to delete actor", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to delete actor"))

			return
		}

		log.Info("actor deleted", slog.Int("actor_id", req.ActorId))

		render.JSON(w, r, Response{resp.OK()})
	}
}
