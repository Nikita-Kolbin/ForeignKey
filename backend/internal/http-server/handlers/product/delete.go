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
	"strconv"
)

type ProductsDeleter interface {
	DeleteProduct(id int) error
	GetWebsiteById(id int) (adminId int, alias string, err error)
	GetProduct(productId int) (*storage.ProductInfo, error)
}

// NewDelete godoc
// @Summary Delete product
// @Security ApiKeyAuth
// @Description Удаляет товар по id
// @Tags product
// @Produce  json
// @Param id path int true "product id"
// @Success 200 {object} response.Response
// @Router /product/delete/{id} [delete]
func NewDelete(pd ProductsDeleter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.product.NewDelete"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		param := chi.URLParam(r, "id")
		productId, err := strconv.Atoi(param)
		if err != nil {
			log.Error("invalid id", slog.String("id", param))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
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

		idFromToken, role, _, err := jwt_token.ParseToken(token)
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

		product, err := pd.GetProduct(productId)
		if err != nil {
			log.Error("failed to get product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find product"))
			return
		}

		adminId, alias, err := pd.GetWebsiteById(product.WebsiteId)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find website"))
			return
		}
		if adminId != idFromToken {
			log.Info("admin is not owner", slog.String("alias", alias))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("admin is not owner"))
			return
		}

		err = pd.DeleteProduct(productId)
		if err != nil {
			log.Error("failed to delete product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to delete product"))
			return
		}

		log.Info("product deleted", slog.Int("product id", productId))

		render.JSON(w, r, response.OK())
	}
}
