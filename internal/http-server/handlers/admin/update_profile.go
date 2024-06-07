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

type ProfileUpdater interface {
	UpdateAdminProfile(fin, ln, fan, city string, id, ii int) error
}

type UpdateProfileRequest struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FatherName string `json:"father_name"`
	City       string `json:"city"`
	ImageId    int    `json:"image_id"`
}

// NewUpdateProfile godoc
// @Summary Update admin profile
// @Security ApiKeyAuth
// @Tags admin
// @Accept json
// @Produce  json
// @Param input body UpdateProfileRequest true "new profile data"
// @Success 200 {object} response.Response
// @Router /admin/update-profile [put]
func NewUpdateProfile(pu ProfileUpdater, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.NewUpdateProfile"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateProfileRequest

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

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid token format"))
			return
		}

		id, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid token"))
			return
		}
		if role != jwt_token.RoleAdmin {
			log.Info("permission denied", slog.String("role", role))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("permission denied"))
			return
		}

		err = pu.UpdateAdminProfile(req.FirstName, req.LastName, req.FatherName, req.City, id, req.ImageId)
		if err != nil {
			log.Error("failed update profile", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed update profile"))
			return
		}

		log.Info("profile updated", slog.Int("admin id", id))

		render.JSON(w, r, response.OK())
	}
}
