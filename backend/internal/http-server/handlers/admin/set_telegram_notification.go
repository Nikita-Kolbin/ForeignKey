package admin

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type TelegramNotificationChanger interface {
	SetAdminTelegramNotification(adminId, notification int) error
}

type UpdateTelegramNotificationRequest struct {
	Notification int `json:"notification"`
}

// NewTelegramNotification godoc
// @Summary Change telegram notification status
// @Security ApiKeyAuth
// @Tags admin
// @Accept json
// @Produce  json
// @Param input body UpdateTelegramNotificationRequest true "notification status"
// @Success 200 {object} response.Response
// @Router /admin/set-telegram-notification [patch]
func NewTelegramNotification(tnc TelegramNotificationChanger, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.admin.NewTelegramNotification"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateEmailNotificationRequest

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

		err = tnc.SetAdminTelegramNotification(id, req.Notification)
		if errors.Is(err, storage.ErrInvalidNotification) {
			log.Error("failed to set notification status", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid notification status"))
			return
		}
		if err != nil {
			log.Error("failed to set notification status", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to set notification status"))
			return
		}

		log.Info("notification status changed", slog.Int("admin id", id))

		render.JSON(w, r, response.OK())
	}
}
