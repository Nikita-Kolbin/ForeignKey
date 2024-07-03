package client

type EmailClient interface {
	Send(receiverEmail, msg string) error
}

type TelegramClient interface {
	Send(username, msg string) error
}
