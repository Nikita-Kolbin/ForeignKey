package website

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

type WebsitesCreator interface {
	CreateWebsite(alias string, adminId int) error
}

type CreateRequest struct {
	Alias string `json:"alias"`
}

// NewCreate godoc
// @Summary Create website
// @Security ApiKeyAuth
// @Tags website
// @Accept json
// @Produce  json
// @Param input body CreateRequest true "alias new website"
// @Success 200 {object} response.Response
// @Router /website/create [post]
func NewCreate(wc WebsitesCreator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.website.NewCreate"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, response.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		token, err := jwt_token.GetTokenFromRequest(r)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token format"))

			return
		}

		id, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token"))

			return
		}

		err = wc.CreateWebsite(req.Alias, id)
		if errors.Is(err, storage.ErrAliasTaken) {
			log.Error("alias already taken", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("alias already taken"))

			return
		}
		if err != nil {
			log.Error("failed to create website", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to create website"))

			return
		}

		log.Info("website created", slog.String("alias", req.Alias))

		render.JSON(w, r, response.OK())
	}
}
