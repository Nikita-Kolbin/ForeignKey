package cart

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

type ItemsCreator interface {
	GetCartId(customerId int) (int, error)
	CreateCartItem(cartId, productId, count int) error
	GetWebsite(alias string) (websiteId, adminId int, err error)
	GetProduct(productId int) (*storage.ProductInfo, error)
}

type AddRequest struct {
	ProductId int `json:"product_id"`
	Count     int `json:"count"`
}

// NewAdd godoc
// @Summary Create curt item
// @Description Добавляет товар в корзину, если товар уже в корзине, увеличивает количество
// @Security ApiKeyAuth
// @Tags cart
// @Accept json
// @Produce  json
// @Param input body AddRequest true "product id and count"
// @Success 200 {object} response.Response
// @Router /cart/add [post]
func NewAdd(ic ItemsCreator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cart.NewAdd"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req AddRequest

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

		customerId, role, alias, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token"))

			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))

			render.JSON(w, r, response.Error("only customers can by products"))

			return
		}

		customerWebsiteId, _, err := ic.GetWebsite(alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("can't find website"))

			return
		}
		product, err := ic.GetProduct(req.ProductId)
		if err != nil {
			log.Error("failed to get product", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("can't find product"))

			return
		}
		if customerWebsiteId != product.WebsiteId {
			log.Info("wrong website", slog.Int(
				"customer", customerWebsiteId),
				slog.Int("product", product.WebsiteId),
			)

			render.JSON(w, r, response.Error("websites id is not equal"))

			return
		}

		cartId, err := ic.GetCartId(customerId)
		if err != nil {
			log.Error("failed to get cart", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("can't find cart"))

			return
		}

		err = ic.CreateCartItem(cartId, req.ProductId, req.Count)
		if err != nil {
			log.Error("failed to create cart item", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to create cart item"))

			return
		}

		log.Info("cart item created", slog.Int("cart id", cartId))

		render.JSON(w, r, response.OK())
	}
}
