package admin

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ProfileGetter interface {
	GetAdminById(id int) (*storage.Admin, error)
}

type GetProfileResponse struct {
	response.Response
	Profile *storage.Admin `json:"profile"`
}

// NewGetProfile godoc
// @Summary Get admin profile
// @Security ApiKeyAuth
// @Tags admin
// @Produce  json
// @Success 200 {object} GetProfileResponse
// @Router /admin/get-profile [get]
func NewGetProfile(pg ProfileGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.NewGetProfile"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("invalid token format"))
			return
		}

		id, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("invalid token"))
			return
		}
		if role != jwt_token.RoleAdmin {
			log.Info("permission denied", slog.String("role", role))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, responseError("permission denied"))
			return
		}

		admin, err := pg.GetAdminById(id)
		if err != nil {
			log.Error("failed to get admin", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get admin"))
			return
		}

		log.Info("admin profile given", slog.Int("admin id", id))

		render.JSON(w, r, responseOk(admin))
	}
}

func responseError(msg string) GetProfileResponse {
	return GetProfileResponse{
		response.Error(msg),
		nil,
	}
}

func responseOk(admin *storage.Admin) GetProfileResponse {
	return GetProfileResponse{
		response.OK(),
		admin,
	}
}
