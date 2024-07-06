package storage

type UserStorage interface {
	CreateUser(userId string, chatId int) error
	GetChatIdByUsername(username string) (int, error)
}
