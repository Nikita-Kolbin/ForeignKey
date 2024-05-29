package order

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OrdersMaker interface {
	CreateOrder(customerId int) error
	GetCustomer(id int) (*storage.Customer, error)
}

type EmailSender interface {
	Send(receiverEmail, msg string) error
}

// NewMakeOrder godoc
// @Summary Make order
// @Description Создает заказ исходя из корзины покупателя
// @Security ApiKeyAuth
// @Tags order
// @Produce json
// @Success 200 {object} response.Response
// @Router /order/make [post]
func NewMakeOrder(om OrdersMaker, es EmailSender, log *slog.Logger) http.HandlerFunc {
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
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid token format"))
			return
		}

		customerId, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid token"))
			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("only customers can make order"))
			return
		}

		err = om.CreateOrder(customerId)
		if errors.Is(err, storage.ErrEmptyOrder) {
			log.Error("failed to make order", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("order is empty"))
			return
		}
		if err != nil {
			log.Error("failed to make order", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to make order"))
			return
		}

		log.Info("order made", slog.Int("customer id", customerId))

		render.JSON(w, r, response.OK())

		customer, err := om.GetCustomer(customerId)
		if err != nil {
			log.Error("can't get customer", slog.String("err", err.Error()))
			return
		}

		err = es.Send(customer.Email, "Заказ оформлен")
		if err != nil {
			log.Error("can't send message", slog.String("err", err.Error()))
			return
		}
	}
}
