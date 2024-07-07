package order

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type StatusChanger interface {
	SetOrderStatus(orderId int, status int) error
	GetOrderById(id int) (*storage.Order, error)
	GetCustomer(id int) (*storage.Customer, error)
	GetWebsiteById(id int) (adminId int, alias string, err error)
}

type UpdateStatusRequest struct {
	OrderId int `json:"order_id"`
	Status  int `json:"status"`
}

// NewSetStatus godoc
// @Summary Change order status
// @Security ApiKeyAuth
// @Tags order
// @Accept json
// @Produce  json
// @Param input body UpdateStatusRequest true "style to website"
// @Success 200 {object} response.Response
// @Router /order/set-status [patch]
func NewSetStatus(sc StatusChanger, ns NotificationService, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.order.NewSetStatus"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req UpdateStatusRequest

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

		order, err := sc.GetOrderById(req.OrderId)
		if err != nil {
			log.Error("failed to get order", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find order"))
			return
		}

		customer, err := sc.GetCustomer(order.CustomerId)
		if err != nil {
			log.Error("failed to get customer", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find customer"))
			return
		}

		adminId, alias, err := sc.GetWebsiteById(customer.WebsiteId)
		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", alias))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("admin is not owner"))
			return
		}

		if req.Status < 0 || req.Status > 7 {
			log.Error("wrong status")
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("wrong status"))
			return
		}
		err = sc.SetOrderStatus(req.OrderId, req.Status)
		if err != nil {
			log.Error("failed to change status", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to change status"))
			return
		}

		log.Info("status changed", slog.Int("order id", req.OrderId))

		render.JSON(w, r, response.OK())

		go sendNotification(customer, &req, ns, log)
	}
}

func sendNotification(
	customer *storage.Customer,
	req *UpdateStatusRequest,
	ns NotificationService,
	log *slog.Logger,
) {
	const op = "handlers.order.SetStatus.sendNotification"

	log = log.With(
		slog.String("op", op),
	)

	var status string
	switch req.Status {
	case storage.StatusAwaitingConfirm:
		status = "Ожидает подтверждения"
	case storage.StatusAcceptedForProcessing:
		status = "Принят в работу"
	case storage.StatusInProgress:
		status = "В работе"
	case storage.StatusMade:
		status = "Заказ готов"
	case storage.StatusSent:
		status = "Заказ отправлен"
	case storage.StatusDelivered:
		status = "Заказ доставлен"
	case storage.StatusCompleted:
		status = "Заказ завершен"
	case storage.StatusUnusualSituation:
		status = "Нестандартная ситуация"
	default:
		status = "Неизвестный статус"
	}

	msg := fmt.Sprintf("Статус заказа %d обновлен: '%s'", req.OrderId, status)

	if customer.EmailNotification == 1 {
		err := ns.SendEmail(customer.Email, msg)
		if err != nil {
			log.Error("can't send email", slog.String("err", err.Error()))
		}
	}

	if customer.TelegramNotification == 1 {
		err := ns.SendTelegram(customer.Telegram, msg)
		if err != nil {
			log.Error("can't send telegram", slog.String("err", err.Error()))
		}
	}
}
