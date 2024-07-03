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

type Updater interface {
	UpdateProduct(name, description, imagesId, tags string, id, price int) error
	GetWebsiteById(id int) (adminId int, alias string, err error)
	GetProduct(productId int) (*storage.ProductInfo, error)
}

type UpdateProductRequest struct {
	ProductId    int          `json:"product_id"`
	ProductsInfo ProductsInfo `json:"product_info"`
}

type ProductsInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImagesId    string `json:"images_id"`
	Tags        string `json:"tags"`
}

// NewUpdate godoc
// @Summary Update product info
// @Security ApiKeyAuth
// @Tags product
// @Accept json
// @Produce  json
// @Param input body UpdateProductRequest true "new profile data"
// @Success 200 {object} response.Response
// @Router /product/update [put]
func NewUpdate(u Updater, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.product.NewUpdate"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateProductRequest

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

		product, err := u.GetProduct(req.ProductId)
		if err != nil {
			log.Error("failed to get product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find product"))
			return
		}

		adminId, alias, err := u.GetWebsiteById(product.WebsiteId)
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

		pi := req.ProductsInfo
		err = u.UpdateProduct(pi.Name, pi.Description, pi.ImagesId, pi.Tags, req.ProductId, pi.Price)
		if errors.Is(err, storage.ErrInvalidImagesIs) {
			log.Error("failed to update product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid images id"))
			return
		}
		if err != nil {
			log.Error("failed to update product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to update product"))
			return
		}

		log.Info("product updates", slog.Int("product id", req.ProductId))

		render.JSON(w, r, response.OK())
	}
}
