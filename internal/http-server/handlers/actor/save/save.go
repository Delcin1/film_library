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
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Birthdate string `json:"birthdate"`
}

type Response struct {
	resp.Response
	ActorId int `json:"actor_id"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=ActorSaver
type ActorSaver interface {
	SaveActor(name string, gender string, birthdate string) (int, error)
}

func New(log *slog.Logger, actorSaver ActorSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.actor.save.New"

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

		actorId, err := actorSaver.SaveActor(req.Name, req.Gender, req.Birthdate)
		if err != nil {
			log.Error("failed to save actor", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save actor"))

			return
		}

		log.Info("actor saved", slog.Int("actor_id", actorId))

		render.JSON(w, r, Response{
			resp.OK(),
			actorId,
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if len(req.Name) < 1 || len(req.Name) > 255 {
		return false, slog.String("field", "name"), "field name is not valid"
	}
	if req.Gender != "male" && req.Gender != "female" {
		return false, slog.String("field", "gender"), "field gender is not valid"
	}
	if _, err := time.Parse("2006-01-02", req.Birthdate); err != nil {
		return false, slog.String("field", "birthdate"), "field birthdate is not valid"
	}
	return true, slog.Attr{}, ""
}
