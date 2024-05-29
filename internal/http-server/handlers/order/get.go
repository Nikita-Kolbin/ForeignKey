package order

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type OrdersGetter interface {
	GetOrders(customerId int) ([]storage.Order, error)
}

type GetResponse struct {
	response.Response
	Orders []storage.Order `json:"orders"`
}

// NewGet godoc
// @Summary Get all orders
// @Security ApiKeyAuth
// @Tags order
// @Produce json
// @Success 200 {object} GetResponse
// @Router /order/get [get]
func NewGet(og OrdersGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.order.NewGet"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseGetError("invalid token format"))
			return
		}

		customerId, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseGetError("invalid token"))
			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, responseGetError("only customers can get order"))
			return
		}

		orders, err := og.GetOrders(customerId)
		if err != nil {
			log.Error("failed to get cart", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseGetError("can't find cart"))
			return
		}

		log.Info("orders given", slog.Int("customer id", customerId))

		render.JSON(w, r, responseGetOk(orders))
	}
}

func responseGetOk(orders []storage.Order) GetResponse {
	return GetResponse{
		response.OK(),
		orders,
	}
}

func responseGetError(msg string) GetResponse {
	return GetResponse{
		response.Error(msg),
		nil,
	}
}
