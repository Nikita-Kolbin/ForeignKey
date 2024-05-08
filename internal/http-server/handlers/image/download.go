package image

import (
	"ForeignKey/internal/http-server/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type ImagesGetter interface {
	Get(imagePath string) ([]byte, error)
}

type ImagesPathGetter interface {
	GetImagePath(id int) (string, error)
}

// NewDownload godoc
// @Summary      DownloadImage
// @Description  При удачном запросе вернет картинку в body со статусом 200, при неудачном json с ошибкой
// @Tags         files
// @Produce      jpeg
// @Produce      png
// @Produce      json
// @Param id path int true "image id"
// @Success      200  {object}   []byte
// @Router       /image/download/{id} [get]
func NewDownload(ig ImagesGetter, ipg ImagesPathGetter, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.image.NewDownload"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("can't parse id", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		path, err := ipg.GetImagePath(id)
		if err != nil {
			log.Error("can't get image path", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid id"))
			return
		}

		img, err := ig.Get(path)
		if err != nil {
			log.Error("can't get image", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("can't find image"))
			return
		}

		_, err = w.Write(img)
		if err != nil {
			log.Error("can't write response", slog.String("err", err.Error()))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("can't find image"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
