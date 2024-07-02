package product

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

type ActiveChanger interface {
	SetProductActive(productId, active int) error
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type UpdateActiveRequest struct {
	Alias     string `json:"alias"`
	ProductId int    `json:"product_id"`
	Active    int    `json:"active"`
}

// NewSetActive godoc
// @Summary Change product active status
// @Security ApiKeyAuth
// @Tags product
// @Accept json
// @Produce  json
// @Param input body UpdateActiveRequest true "active status"
// @Success 200 {object} response.Response
// @Router /product/set-active [patch]
func NewSetActive(ac ActiveChanger, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.product.NewSetActive"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateActiveRequest

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

		_, adminId, err := ac.GetWebsite(req.Alias)
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

		err = ac.SetProductActive(req.ProductId, req.Active)
		if errors.Is(err, storage.ErrInvalidActive) {
			log.Error("failed to set active", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid active"))
			return
		}
		if err != nil {
			log.Error("failed to set active", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to set active"))
			return
		}

		log.Info("active changed", slog.Int("product id", req.ProductId))

		render.JSON(w, r, response.OK())
	}
}
