package signin

import (
	"film_library/internal/lib/api/response"
	"film_library/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	response.Response
	Token string `json:"token"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=UserAuthenticator
type UserAuthenticator interface {
	GetUser(username string, password string) (int, error)
}

// @Summary		Sign in a user
// @Description	Sign in a user by username and password
// @Tags			User
// @Accept			json
// @Produce		json
// @Param			username	body		string	true	"Username"
// @Param			password	body		string	true	"Password"
// @Success		200			{object}	Response
// @Failure		400			{object}	response.Response
// @Router			/signin [post]
func New(log *slog.Logger, userAuthenticator UserAuthenticator, jwtKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.signin.New"

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

		userId, err := userAuthenticator.GetUser(req.Username, req.Password)
		if err != nil {
			log.Error("failed to authenticate user", sl.Err(err))

			render.JSON(w, r, response.Error("failed to authenticate user"))

			return
		}

		token, err := generateToken(userId, jwtKey)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))

			render.JSON(w, r, response.Error("failed to generate token"))

			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Token:    token,
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

func generateToken(userId int, jwtKey string) (string, error) {
	tokenAuth := jwtauth.New("HS256", []byte(jwtKey), nil)

	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{"user_id": userId})
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
