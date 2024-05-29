package product

import (
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ProductsGetter interface {
	GetProducts(websiteId int) ([]storage.ProductInfo, error)
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type GetResponse struct {
	response.Response
	Products []storage.ProductInfo `json:"products"`
}

// NewGetByAlias godoc
// @Summary GetByAlias
// @Description Возвращает все товары сайта по алиасу
// @Tags product
// @Produce  json
// @Param alias path string true "website alias"
// @Success 200 {object} GetResponse
// @Router /product/get-by-alias/{alias} [get]
func NewGetByAlias(pg ProductsGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.products.NewCreate"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		websiteId, _, err := pg.GetWebsite(alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("failed to find website"))
			return
		}

		products, err := pg.GetProducts(websiteId)
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

func responseError(msg string) GetResponse {
	return GetResponse{
		response.Error(msg),
		nil,
	}
}

func responseOk(products []storage.ProductInfo) GetResponse {
	return GetResponse{
		response.OK(),
		products,
	}
}
