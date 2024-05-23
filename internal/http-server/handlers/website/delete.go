package website

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type WebsitesDeleter interface {
	GetWebsite(alias string) (websiteId, adminId int, err error)
	DeleteWebsite(alias string) error
}

// NewDelete godoc
// @Summary Delete website
// @Security ApiKeyAuth
// @Description Удаляет сайт по алиасу
// @Tags website
// @Produce  json
// @Param alias path string true "website alias"
// @Success 200 {object} response.Response
// @Router /website/delete/{alias} [delete]
func NewDelete(wd WebsitesDeleter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.website.NewDelete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token format"))

			return
		}

		id, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token"))

			return
		}
		if role != jwt_token.RoleAdmin {
			log.Info("permission denied", slog.String("role", role))

			render.JSON(w, r, response.Error("permission denied"))

			return
		}

		_, adminId, err := wd.GetWebsite(alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to find website"))

			return
		}
		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", alias))

			render.JSON(w, r, response.Error("admin is not owner"))

			return
		}

		err = wd.DeleteWebsite(alias)
		if err != nil {
			log.Error("failed to delete website", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to delete website"))

			return
		}

		log.Info("website deleted", slog.String("alias", alias))

		render.JSON(w, r, response.OK())
	}
}
