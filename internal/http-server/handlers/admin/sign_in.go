package admin

import (
	"ForeignKey/internal/http-server/jwt_token"
	resp "ForeignKey/internal/http-server/response"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type AdminsGetter interface {
	GetAdminId(username, password string) (int, error)
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignInResponse struct {
	resp.Response
	Token string `json:"token"`
}

func NewSignIn(ac AdminsGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.NewSignIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignInRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		id, err := ac.GetAdminId(req.Username, req.Password)
		if err != nil {
			log.Error("failed to get admin", slog.String("err", err.Error()))

			render.JSON(w, r, resp.Error("wrong login or password"))

			return
		}

		t, err := jwt_token.GenerateToken(id)
		if err != nil {
			log.Error("failed to generate token", slog.String("err", err.Error()))

			render.JSON(w, r, resp.Error("failed to generate token"))

			return
		}

		log.Info("admin sing in", slog.String("username", req.Username))

		render.JSON(w, r, responseOK(t))
	}
}

func responseOK(token string) SignInResponse {
	return SignInResponse{
		Response: resp.OK(),
		Token:    token,
	}
}
