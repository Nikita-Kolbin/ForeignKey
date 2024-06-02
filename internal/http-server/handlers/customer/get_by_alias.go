package customer

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type WebsiteCustomersGetter interface {
	GetCustomersByWebsite(websiteId int) ([]storage.Customer, error)
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type GetResponse struct {
	response.Response
	Customers []storage.Customer `json:"customers"`
}

// NewGetByAlias godoc
// @Summary Get all customers by alias
// @Security ApiKeyAuth
// @Tags customer
// @Param alias path string true "website alias"
// @Produce json
// @Success 200 {object} GetResponse
// @Router /customer/get-by-alias/{alias} [get]
func NewGetByAlias(wcg WebsiteCustomersGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.order.NewGetByAlias"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseGetError("invalid token format"))
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

		websiteId, adminId, err := wcg.GetWebsite(alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find website"))
			return
		}

		if adminId != id {
			log.Info("admin is not owner", slog.String("alias", alias))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, response.Error("admin is not owner"))
			return
		}

		customers, err := wcg.GetCustomersByWebsite(websiteId)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseGetError("can't find website's customers"))
			return
		}

		log.Info("orders given", slog.Int("customer id", id))

		render.JSON(w, r, responseGetOk(customers))
	}
}

func responseGetOk(orders []storage.Customer) GetResponse {
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
