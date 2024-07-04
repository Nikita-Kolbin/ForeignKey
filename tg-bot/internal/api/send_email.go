package api

import (
	"errors"
	"github.com/go-chi/render"
	"io"
	"log"
	"net/http"
)

type EmailSender interface {
	Send(receiverEmail, msg string) error
}

type SendEmailRequest struct {
	ReceiverEmail string `json:"receiver_email"`
	Message       string `json:"message"`
}

// NewSendEmail godoc
// @Summary Send email
// @Tags notification
// @Accept json
// @Produce  json
// @Param input body SendEmailRequest true "notification info"
// @Success 200 {object} Response
// @Router /send-email [post]
func NewSendEmail(es EmailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req SendEmailRequest

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

		err = es.Send(req.ReceiverEmail, req.Message)
		if err != nil {
			log.Printf("failed to send email: %s", err.Error())
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, responseError("failed to send email"))
			return
		}

		log.Printf("email sent: %s", req.ReceiverEmail)
		render.JSON(w, r, responseOk())
	}
}
