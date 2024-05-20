package order

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OrdersMaker interface {
	CreateOrder(customerId int) error
}

// NewMakeOrder godoc
// @Summary Make order
// @Description Создает заказ исходя из корзины покупателя
// @Security ApiKeyAuth
// @Tags order
// @Produce json
// @Success 200 {object} response.Response
// @Router /order/make [post]
func NewMakeOrder(om OrdersMaker, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.cart.NewMakeOrder"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token format"))

			return
		}

		customerId, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token"))

			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))

			render.JSON(w, r, response.Error("only customers can make order"))

			return
		}

		err = om.CreateOrder(customerId)
		if err != nil {
			log.Error("failed to make order", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to make order"))

			return
		}

		log.Info("order made", slog.Int("customer id", customerId))

		render.JSON(w, r, response.OK())
	}
}
