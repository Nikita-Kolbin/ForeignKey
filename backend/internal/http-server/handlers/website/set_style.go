package website

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

type StyleChanger interface {
	UpdateStyle(alias, backgroundColor, font string) error
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type StyleRequest struct {
	Alias           string `json:"alias"`
	BackgroundColor string `json:"background_color"`
	Font            string `json:"font"`
}

// NewSetStyle godoc
// @Summary Change style
// @Security ApiKeyAuth
// @Tags website
// @Accept json
// @Produce  json
// @Param input body StyleRequest true "style to website"
// @Success 200 {object} response.Response
// @Router /website/set-style [patch]
func NewSetStyle(sc StyleChanger, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.website.NewSetStyle"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req StyleRequest

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

		_, adminId, err := sc.GetWebsite(req.Alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find website"))
			return
		}
		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", req.Alias))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("admin is not owner"))
			return
		}

		err = sc.UpdateStyle(req.Alias, req.BackgroundColor, req.Font)
		if err != nil {
			log.Error("failed to change style", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to change style"))
			return
		}

		log.Info("style changed", slog.String("alias", req.Alias))

		render.JSON(w, r, response.OK())
	}
}
