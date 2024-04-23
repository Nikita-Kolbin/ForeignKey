package admin

import (
	resp "ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type AdminsCreator interface {
	CreateAdmin(username, password string) error
}

type SignUpRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSignUp(ac AdminsCreator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.NewSignUp"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignUpRequest

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

		err = ac.CreateAdmin(req.Username, req.Password)
		if errors.Is(err, storage.ErrUsernameTaken) {
			log.Error("username is already taken", slog.String("username", req.Username))

			render.JSON(w, r, resp.Error("username is already taken"))

			return
		}
		if err != nil {
			log.Error("failed to create admin", slog.String("err", err.Error()))

			render.JSON(w, r, resp.Error("failed to create admin"))

			return
		}

		log.Info("admin created", slog.String("username", req.Username))

		render.JSON(w, r, resp.OK())
	}
}
