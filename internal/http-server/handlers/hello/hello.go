package hello

import (
	resp "ForeignKey/internal/http-server/response"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type Response struct {
	resp.Response
	Message string `json:"message"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.hello.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		name := chi.URLParam(r, "name")
		if name == "" {
			log.Info("name is empty")

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		msg := fmt.Sprintf("Hello, %s!", name)

		response := responseOK(msg)

		log.Info("say hello", slog.String("name", name))

		render.JSON(w, r, response)
	}
}

func responseOK(msg string) Response {
	return Response{
		Response: resp.OK(),
		Message:  msg,
	}
}
