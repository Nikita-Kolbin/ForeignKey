package api

import (
	"errors"
	"github.com/go-chi/render"
	"io"
	"log"
	"net/http"
)

type TgSender interface {
	Send(username, msg string) error
}

type SendTelegramRequest struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

// NewSendTelegram godoc
// @Summary Telegram email
// @Tags notification
// @Accept json
// @Produce  json
// @Param input body SendTelegramRequest true "notification info"
// @Success 200 {object} Response
// @Router /send-telegram [post]
func NewSendTelegram(tgs TgSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendTelegramRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Println("request body is empty")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("empty request"))
			return
		}
		if err != nil {
			log.Printf("failed to decode request body: %s", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("failed to decode request"))
			return
		}

		err = tgs.Send(req.Username, req.Message)
		if err != nil {
			log.Printf("failed to send telegram: %s", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("failed to send telegram"))
			return
		}

		log.Printf("telegram sent: %s", req.Username)
		render.JSON(w, r, responseOk())
	}
}
