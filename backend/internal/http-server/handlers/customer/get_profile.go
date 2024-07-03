package customer

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ProfileGetter interface {
	GetCustomer(id int) (*storage.Customer, error)
}

type GetProfileResponse struct {
	response.Response
	Profile *storage.Customer `json:"profile"`
}

// NewGetProfile godoc
// @Summary Get customer profile
// @Security ApiKeyAuth
// @Tags customer
// @Produce  json
// @Success 200 {object} GetProfileResponse
// @Router /customer/get-profile [get]
func NewGetProfile(pg ProfileGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.customer.NewGetProfile"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		auth := r.Header.Get("Authorization")
		token, err := jwt_token.GetTokenFromRequest(auth)
		if err != nil {
			log.Error("failed to get token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("invalid token format"))
			return
		}

		id, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("invalid token"))
			return
		}
		if role != jwt_token.RoleCustomer {
			log.Info("permission denied", slog.String("role", role))
			render.Status(r, http.StatusForbidden)
			render.JSON(w, r, responseError("permission denied"))
			return
		}

		customer, err := pg.GetCustomer(id)
		if err != nil {
			log.Error("failed to get customer", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get customer"))
			return
		}

		log.Info("customer profile given", slog.Int("customer id", id))

		render.JSON(w, r, responseOk(customer))
	}
}

func responseError(msg string) GetProfileResponse {
	return GetProfileResponse{
		response.Error(msg),
		nil,
	}
}

func responseOk(customer *storage.Customer) GetProfileResponse {
	return GetProfileResponse{
		response.OK(),
		customer,
	}
}
