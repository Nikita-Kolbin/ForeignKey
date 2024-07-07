package notification_client

import "errors"

const contentJSON = "application/json"

const (
	statusOk    = "Ok"
	statusError = "Error"
)

var ErrUnknownStatus = errors.New("unknown response status")

type EmailNotification struct {
	ReceiverEmail string `json:"receiver_email"`
	Message       string `json:"message"`
}

type TelegramNotification struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
