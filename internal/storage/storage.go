package storage

import "errors"

var (
	ErrLoginTaken = errors.New("login is already taken")
	ErrAliasTaken = errors.New("alias is already taken")
)

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
	OrderItems []OrderItem `json:"order_items"`
}
