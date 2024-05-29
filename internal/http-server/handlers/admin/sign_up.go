package admin

import (
	"ForeignKey/internal/http-server/response"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewSignUp godoc
// @Summary      SingUp admin
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param input body SignUpRequest true "sign up"
// @Success      200  {object}   response.Response
// @Router       /admin/sign-up [post]
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
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		err = ac.CreateAdmin(req.Email, req.Password)
		if errors.Is(err, storage.ErrInvalidEmail) {
			log.Error("email is invalid", slog.String("email", req.Email))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("email is invalid"))
			return
		}
		if errors.Is(err, storage.ErrEmailRegistered) {
			log.Error("email is already taken", slog.String("email", req.Email))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("email is already taken"))
			return
		}
		if err != nil {
			log.Error("failed to create admin", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create admin"))
			return
		}

		log.Info("admin created", slog.String("email", req.Email))

		render.JSON(w, r, response.OK())
	}
}
