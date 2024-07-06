package client

import "tg-bot/internal/client/tgClient"

type EmailClient interface {
	Send(receiverEmail, msg string) error
}

type TelegramClient interface {
	Send(chatId int, msg string) error
	Updates(offset, limit int) ([]tgClient.Update, error)
}
