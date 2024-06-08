package storage

import "errors"

var (
	ErrEmailRegistered = errors.New("email is already registered")
	ErrAliasTaken      = errors.New("alias is already taken")
	ErrEmptyOrder      = errors.New("order is empty")
	ErrInvalidEmail    = errors.New("email is invalid")
)

const (
	DefaultBackgroundColor = "white"
	DefaultFont            = "Arial"
)

type Admin struct {
	Id         int    `json:"id"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	FatherName string `json:"father_name"`
	City       string `json:"city"`
	ImageId    int    `json:"image_id"`
}

type Customer struct {
	Id        int    `json:"id"`
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
	Id         int         `json:"id"`
	CustomerId int         `json:"customer_id"`
	DateTime   string      `json:"date_time"`
	Status     int         `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
}

type OrderStatus int

const (
	StatusAwaitingConfirm OrderStatus = iota
	StatusAcceptedForProcessing
	StatusInProgress
	StatusMade
	StatusSent
	StatusDelivered
	StatusUnusualSituation
)
