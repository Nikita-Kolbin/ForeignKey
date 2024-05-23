package product

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

type ProductsCreator interface {
	CreateProduct(name, description string, websiteId, price, imageId int) error
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type CreateRequest struct {
	Alias       string `json:"alias"`
	ProductInfo Info   `json:"product_info"`
}

type Info struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImageId     int    `json:"image_id"`
}

// NewCreate godoc
// @Summary Create product
// @Security ApiKeyAuth
// @Tags product
// @Accept json
// @Produce  json
// @Param input body CreateRequest true "alias new website"
// @Success 200 {object} response.Response
// @Router /product/create [post]
func NewCreate(pc ProductsCreator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.products.NewCreate"

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

		websiteId, adminId, err := pc.GetWebsite(req.Alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to find website"))

			return
		}

		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", req.Alias))

			render.JSON(w, r, response.Error("admin is not owner"))

			return
		}

		pi := req.ProductInfo
		err = pc.CreateProduct(pi.Name, pi.Description, websiteId, pi.Price, pi.ImageId)
		if err != nil {
			log.Error("failed to create product", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to create product"))

			return
		}

		log.Info("product created", slog.String("website alias", req.Alias))

		render.JSON(w, r, response.OK())
	}
}
