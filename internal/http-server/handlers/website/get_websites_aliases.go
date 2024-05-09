package website

import (
	"ForeignKey/internal/http-server/jwt_token"
	"ForeignKey/internal/http-server/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type AliasesGetter interface {
	GetWebsitesAliases(adminId int) ([]string, error)
}

type AliasesResponse struct {
	response.Response
	Aliases []string `json:"aliases"`
}

// NewGetAliases godoc
// @Summary Get all users aliases
// @Security ApiKeyAuth
// @Tags website
// @Produce  json
// @Success 200 {object} AliasesResponse
// @Router /website/aliases [get]
func NewGetAliases(ag AliasesGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.website.NewCreate"

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

		id, role, _, err := jwt_token.ParseToken(token)
		if err != nil {
			log.Error("failed to parse token", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("invalid token"))

			return
		}
		if role != jwt_token.RoleAdmin {
			log.Info("permission denied", slog.String("role", role))

			render.JSON(w, r, response.Error("permission denied"))

			return
		}

		aliases, err := ag.GetWebsitesAliases(id)
		if err != nil {
			log.Error("failed to get aliases", slog.String("err", err.Error()))

			render.JSON(w, r, response.Error("failed to get aliases"))

			return
		}

		log.Info("give aliases", slog.Int("admin id", id))

		render.JSON(w, r, responseOK(aliases))
	}
}

func responseOK(aliases []string) AliasesResponse {
	return AliasesResponse{
		Response: response.OK(),
		Aliases:  aliases,
	}
}
