package storage

import "errors"

var (
	ErrEmailRegistered  = errors.New("email is already registered")
	ErrAliasTaken       = errors.New("alias is already taken")
	ErrEmptyOrder       = errors.New("order is empty")
	ErrInvalidEmail     = errors.New("email is invalid")
	ErrInvalidImagesIs  = errors.New("invalid images id")
	ErrAdminHaveWebsite = errors.New("admin already have website")
	ErrInvalidActive    = errors.New("invalid active status")
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
	Telegram   string `json:"telegram"`
	ImageId    int    `json:"image_id"`
}

type Customer struct {
	Id           int    `json:"id"`
	WebsiteId    int    `json:"website_id"`
	Email        string `json:"email"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	FatherName   string `json:"father_name"`
	Phone        string `json:"phone"`
	Telegram     string `json:"telegram"`
	DeliveryType string `json:"delivery_type"`
	PaymentType  string `json:"payment_type"`
}

type ProductInfo struct {
	Id          int    `json:"id"`
	WebsiteId   int    `json:"website_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ImagesId    string `json:"images_id"`
	Active      int    `json:"active"`
	Tags        string `json:"tags"`
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

const (
	StatusAwaitingConfirm int = iota
	StatusAcceptedForProcessing
	StatusInProgress
	StatusMade
	StatusSent
	StatusDelivered
	StatusCompleted
	StatusUnusualSituation
)
