package admin

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	response.Response
	Token string `json:"token"`
}

// NewSignIn godoc
// @Summary      SingIn admin
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param input body SignInRequest true "sign in"
// @Success      200  {object}   SignInResponse
// @Router       /admin/sign-in [post]
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

		id, err := ac.GetAdminId(req.Email, req.Password)
		if err != nil {
			log.Error("failed to get admin", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("wrong email or password"))
			return
		}

		t, err := jwt_token.GenerateToken(id, jwt_token.RoleAdmin, "")
		if err != nil {
			log.Error("failed to generate token", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate token"))
			return
		}

		log.Info("admin sing in", slog.String("email", req.Email))

		render.JSON(w, r, responseOK(t))
	}
}

func responseOK(token string) SignInResponse {
	return SignInResponse{
		Response: response.OK(),
		Token:    token,
	}
}
