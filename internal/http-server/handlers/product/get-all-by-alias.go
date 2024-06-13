package product

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ProductsAdminGetter interface {
	GetAllProducts(websiteId int) ([]storage.ProductInfo, error)
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

// NewGetAllByAlias godoc
// @Summary GetAllByAlias
// @Security ApiKeyAuth
// @Description Возвращает ВСЕ ТОВАРЫ сайта по алиасу, только для админа
// @Tags product
// @Produce  json
// @Param alias path string true "website alias"
// @Success 200 {object} GetResponse
// @Router /product/get-all-by-alias/{alias} [get]
func NewGetAllByAlias(pg ProductsAdminGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.products.NewGetAllByAlias"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

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

		websiteId, adminId, err := pg.GetWebsite(alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find website"))
			return
		}

		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", alias))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("admin is not owner"))
			return
		}

		products, err := pg.GetAllProducts(websiteId)
		if err != nil {
			log.Error("failed to get products", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseError("failed to get products"))
			return
		}

		log.Info("product given", slog.String("website alias", alias))

		render.JSON(w, r, responseOk(products))
	}
}
