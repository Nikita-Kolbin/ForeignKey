package cart

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type CartsGetter interface {
	GetCartItems(cartId int) ([]storage.CartItem, error)
	GetCartId(customerId int) (int, error)
}

type GetResponse struct {
	response.Response
	CartItems []storage.CartItem `json:"cart_items"`
}

// NewGet godoc
// @Summary Get all cart items
// @Security ApiKeyAuth
// @Tags cart
// @Produce  json
// @Success 200 {object} GetResponse
// @Router /cart/get [get]
func NewGet(cg CartsGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cart.NewGet"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))

			render.JSON(w, r, responseError("invalid token format"))

			return
		}

		customerId, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, responseError("invalid token"))

			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))

			render.JSON(w, r, responseError("only customers can by products"))

			return
		}

		cartId, err := cg.GetCartId(customerId)
		if err != nil {
			log.Error("failed to get cart", slog.String("err", err.Error()))

			render.JSON(w, r, responseError("can't find cart"))

			return
		}

		cartItems, err := cg.GetCartItems(cartId)
		if err != nil {
			log.Error("failed to get cart item", slog.String("err", err.Error()))

			render.JSON(w, r, responseError("failed to get cart item"))

			return
		}

		log.Info("cart items giver", slog.Int("cart id", cartId))

		render.JSON(w, r, responseOk(cartItems))
	}
}

func responseOk(ci []storage.CartItem) GetResponse {
	return GetResponse{
		response.OK(),
		ci,
	}
}

func responseError(msg string) GetResponse {
	return GetResponse{
		response.Error(msg),
		nil,
	}
}
