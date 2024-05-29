package customer

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

type CustomersGetter interface {
	GetCustomerId(websiteId int, login, password string) (int, error)
	GetWebsite(alias string) (websiteId, adminId int, err error)
}

type SignInRequest struct {
	Alias    string `json:"alias"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	response.Response
	Token string `json:"token"`
}

// NewSignIn godoc
// @Summary      SingIn customer
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param input body SignInRequest true "sign in"
// @Success      200  {object}   SignInResponse
// @Router       /customer/sign-in [post]
func NewSignIn(cg CustomersGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.customer.NewSignIn"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignInRequest

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

		websiteId, _, err := cg.GetWebsite(req.Alias)
		if err != nil {
			log.Error("failed to get website", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to find website"))
			return
		}

		id, err := cg.GetCustomerId(websiteId, req.Email, req.Password)
		if err != nil {
			log.Error("failed to get customer", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("wrong email or password"))
			return
		}

		t, err := jwt_token.GenerateToken(id, jwt_token.RoleCustomer, req.Alias)
		if err != nil {
			log.Error("failed to generate token", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to generate token"))
			return
		}

		log.Info("customer sing in", slog.String("email", req.Email))

		render.JSON(w, r, responseOK(t))
	}
}

func responseOK(token string) SignInResponse {
	return SignInResponse{
		Response: response.OK(),
		Token:    token,
	}
}
