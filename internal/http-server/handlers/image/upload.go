package image

import (
	"ForeignKey/internal/http-server/response"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type ImagesSaver interface {
	Save(img []byte, extension string) (string, error)
}

type ImagesCreator interface {
	CreateImage(path string) (int, error)
}

type UploadResponse struct {
	response.Response
	Id int `json:"id"`
}

// NewUpload godoc
// @Summary      UploadImage
// @Description  В боди должна быть картинка в виде массива байт
// @Tags         files
// @Accept       jpeg
// @Accept       png
// @Produce      json
// @Param input body []byte true "byte image"
// @Success      200  {object}   UploadResponse
// @Router       /image/upload [post]
func NewUpload(imagesSever ImagesSaver, imagesCreator ImagesCreator, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.image.NewUpload"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Error("failed to read request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseErr("failed to read request"))
			return
		}

		ext, err := parseExtension(r)
		if err != nil {
			log.Error("failed to read request body", slog.String("err", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseErr("can't parse file extension"))
			return
		}

		path, err := imagesSever.Save(b, ext)
		if err != nil {
			log.Error("failed to save image", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseErr("failed to save image"))
			return
		}

		id, err := imagesCreator.CreateImage(path)
		if err != nil {
			log.Error("failed to create image", slog.String("err", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, responseErr("failed to save image"))
			return
		}

		log.Info("image created", slog.Int("id", id))

		render.JSON(w, r, responseOK(id))
	}
}

func responseOK(id int) UploadResponse {
	return UploadResponse{
		Response: response.OK(),
		Id:       id,
	}
}

func responseErr(msg string) UploadResponse {
	return UploadResponse{
		Response: response.Error(msg),
		Id:       -1,
	}
}

func parseExtension(r *http.Request) (string, error) {
	ct := r.Header.Get("Content-Type")
	sp := strings.Split(ct, "/")
	if len(sp) != 2 {
		return "", fmt.Errorf("can't parse file extension")
	}
	return sp[1], nil
}
