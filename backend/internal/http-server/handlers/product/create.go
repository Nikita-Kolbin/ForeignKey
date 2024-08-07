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

type ProductsCreator interface {
	CreateProduct(name, description, imagesId, tags string, websiteId, price int) error
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
	ImagesId    string `json:"images_id"`
	Tags        string `json:"tags"`
}

// NewCreate godoc
// @Summary Create product
// @Security ApiKeyAuth
// @Tags product
// @Accept json
// @Produce  json
// @Param input body CreateRequest true "product info"
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

		websiteId, adminId, err := pc.GetWebsite(req.Alias)
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

		pi := req.ProductInfo
		err = pc.CreateProduct(pi.Name, pi.Description, pi.ImagesId, pi.Tags, websiteId, pi.Price)
		if errors.Is(err, storage.ErrInvalidImagesIs) {
			log.Error("failed to create product", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid images id"))
			return
		}
		if err != nil {
			log.Error("failed to create product", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create product"))
			return
		}

		log.Info("product created", slog.String("website alias", req.Alias))

		render.JSON(w, r, response.OK())
	}
}
