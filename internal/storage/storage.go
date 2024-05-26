package storage

import "errors"

var (
	ErrLoginTaken   = errors.New("login is already taken")
	ErrAliasTaken   = errors.New("alias is already taken")
	ErrEmptyOrder   = errors.New("order is empty")
	ErrInvalidEmail = errors.New("email is invalid")
)

type Customer struct {
	WebsiteId int    `json:"website_id"`
	Email     string `json:"email"`
}

type ProductInfo struct {
	Id          int    `json:"id"`
	WebsiteId   int    `json:"website_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImageId     int    `json:"image_id"`
}

type CartItem struct {
	Id      int         `json:"id"`
	CartId  int         `json:"cart_id"`
	Product ProductInfo `json:"product"`
	Count   int         `json:"count"`
}

type OrderItem struct {
	Id      int         `json:"id"`
	OrderId int         `json:"order_id"`
	Product ProductInfo `json:"product"`
	Count   int         `json:"count"`
}

type Order struct {
	DateTime   string      `json:"date_time"`
	OrderItems []OrderItem `json:"order_items"`
}
