package website

import (
	"ForeignKey/internal/http-server/response"
	"ForeignKey/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type StyleGetter interface {
	GetWebsiteStyle(alias string) (style *storage.WebsiteStyle, err error)
}

type GetStyleResponse struct {
	response.Response
	Style *storage.WebsiteStyle `json:"style"`
}

// NewGetStyle godoc
// @Summary Get website style
// @Tags website
// @Produce  json
// @Param alias path string true "website alias"
// @Success 200 {object} GetStyleResponse
// @Router /website/get-style/{alias} [get]
func NewGetStyle(pg StyleGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.products.NewGetStyle"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		style, err := pg.GetWebsiteStyle(alias)
		if err != nil {
			log.Error("failed to get style", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseError("failed to get style"))
			return
		}

		log.Info("product given", slog.String("website alias", alias))

		render.JSON(w, r, responseOk(style))
	}
}

func responseError(msg string) GetStyleResponse {
	return GetStyleResponse{
		response.Error(msg),
		nil,
	}
}

func responseOk(style *storage.WebsiteStyle) GetStyleResponse {
	return GetStyleResponse{
		response.OK(),
		style,
	}
}
