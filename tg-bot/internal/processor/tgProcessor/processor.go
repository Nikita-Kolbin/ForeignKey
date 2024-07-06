package tgProcessor

import (
	"fmt"
	"log"
	"tg-bot/internal/client"
	"tg-bot/internal/storage"
	"time"
)

const limit = 100

type Processor struct {
	tgClient client.TelegramClient
	storage  storage.UserStorage
	offset   int
}

func New(tgClient client.TelegramClient, storage storage.UserStorage) *Processor {
	p := &Processor{
		tgClient: tgClient,
		storage:  storage,
		offset:   0,
	}

	go func() {
		for {
			p.processUpdates()
			time.Sleep(time.Millisecond)
		}
	}()

	return p
}

func (p *Processor) processUpdates() {
	updates, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		log.Println("can't, get updates:", err)
	}

	for _, u := range updates {
		err = p.storage.CreateUser(u.Message.From.Username, u.Message.Chat.ID)
		if err != nil {
			log.Println("can't, create user:", err)
		}

		log.Printf("fetch message: %s, from: %s", u.Message.Text, u.Message.From.Username)
	}

	if len(updates) > 0 {
		p.offset = updates[len(updates)-1].ID + 1
	}
}

func (p *Processor) Send(username, message string) error {
	chatId, err := p.storage.GetChatIdByUsername(username)
	if err != nil {
		return fmt.Errorf("can't get chat id: %w", err)
	}

	err = p.tgClient.Send(chatId, message)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	log.Printf("send message: %s, to: %s", message, username)

	return nil
}
