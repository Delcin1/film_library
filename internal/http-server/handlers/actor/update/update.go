package update

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
	ActorId   int     `json:"actor_id"`
	Name      *string `json:"name,omitempty"`
	Gender    *string `json:"gender,omitempty"`
	Birthdate *string `json:"birthdate,omitempty"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorUpdater
type ActorUpdater interface {
	UpdateActorName(actorId int, name string) error
	UpdateActorGender(actorId int, gender string) error
	UpdateActorBirthdate(actorId int, birthdate string) error
}

func New(log *slog.Logger, actorSaver ActorUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor.update.New"

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

		if req.Name != nil {
			err := actorSaver.UpdateActorName(req.ActorId, *req.Name)
			if err != nil {
				log.Error("failed to update actor name", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to update actor name"))

				return
			}
		}

		if req.Gender != nil {
			err := actorSaver.UpdateActorGender(req.ActorId, *req.Gender)
			if err != nil {
				log.Error("failed to update actor gender", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to update actor gender"))

				return
			}
		}

		if req.Birthdate != nil {
			err := actorSaver.UpdateActorBirthdate(req.ActorId, *req.Birthdate)
			if err != nil {
				log.Error("failed to update actor birthdate", sl.Err(err))

				render.JSON(w, r, resp.Error("failed to update actor birthdate"))

				return
			}
		}

		if req.Name == nil && req.Gender == nil && req.Birthdate == nil {
			log.Error("no fields to update")

			render.JSON(w, r, resp.Error("no fields to update"))

			return
		}

		log.Info("actor updated", slog.Int("actor_id", req.ActorId))

		render.JSON(w, r, Response{resp.OK()})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.ActorId <= 0 {
		return false, slog.String("field", "actor_id"), "field actor_id is not valid"
	}
	if req.Name != nil && (len(*req.Name) < 1 || len(*req.Name) > 255) {
		return false, slog.String("field", "name"), "field name is not valid"
	}
	if req.Gender != nil && *req.Gender != "male" && *req.Gender != "female" {
		return false, slog.String("field", "gender"), "field gender is not valid"
	}
	if req.Birthdate == nil {
		return true, slog.Attr{}, ""
	}
	if _, err := time.Parse("2006-01-02", *req.Birthdate); err != nil {
		return false, slog.String("field", "birthdate"), "field birthdate is not valid"
	}
	return true, slog.Attr{}, ""
}
