package signup

import (
	"errors"
	resp "film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"film_library/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	resp.Response
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserSaver
type UserSaver interface {
	SaveUser(username string, password string) error
}

// @Summary		Create a new user
// @Description	Create a new user by username and password
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			username	body		string	true	"Username"
// @Param			password	body		string	true	"Password"
// @Success		200			{object}	response.Response
// @Failure		400			{object}	response.Response
// @Router			/signup [post]
func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.signup.New"

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

		err = userSaver.SaveUser(req.Username, req.Password)
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("user already exists", slog.String("username", req.Username))

			render.JSON(w, r, resp.Error("user already exists"))

			return
		}

		if err != nil {
			log.Error("failed to save user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to save user"))

			return
		}

		render.JSON(w, r, Response{
			Response: resp.OK(),
		})
	}
}

func validateRequest(req Request) (bool, slog.Attr, string) {
	if req.Username == "" {
		return false, slog.String("username", req.Username), "field username is required"
	}

	if req.Password == "" {
		return false, slog.String("password", req.Password), "field password is required"
	}

	return true, slog.Attr{}, ""
}
