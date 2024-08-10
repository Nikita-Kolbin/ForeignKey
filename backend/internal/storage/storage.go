package storage

import (
	"errors"
)

var (
	ErrEmailRegistered     = errors.New("email is already registered")
	ErrAliasTaken          = errors.New("alias is already taken")
	ErrEmptyOrder          = errors.New("order is empty")
	ErrInvalidEmail        = errors.New("email is invalid")
	ErrInvalidImagesIs     = errors.New("invalid images id")
	ErrAdminHaveWebsite    = errors.New("admin already have website")
	ErrInvalidActive       = errors.New("invalid active status")
	ErrInvalidNotification = errors.New("invalid notification status")
)

type Admin struct {
	Id                   int    `json:"id"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	FatherName           string `json:"father_name"`
	City                 string `json:"city"`
	Telegram             string `json:"telegram"`
	ImageId              int    `json:"image_id"`
	TelegramNotification int    `json:"telegram_notification"`
	EmailNotification    int    `json:"email_notification"`
}

type Customer struct {
	Id                   int    `json:"id"`
	WebsiteId            int    `json:"website_id"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	FatherName           string `json:"father_name"`
	Phone                string `json:"phone"`
	Telegram             string `json:"telegram"`
	DeliveryType         string `json:"delivery_type"`
	PaymentType          string `json:"payment_type"`
	TelegramNotification int    `json:"telegram_notification"`
	EmailNotification    int    `json:"email_notification"`
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
	Comment    string      `json:"comment"`
}

type WebsiteStyle struct {
	BackgroundColor string `json:"background_color"`
	TextColor       string `json:"text_color"`
	Font            string `json:"font"`

	MainOne string `json:"main_one"`
	MainTwo string `json:"main_two"`

	AboutOne        string `json:"about_one"`
	AboutTwo        string `json:"about_two"`
	AboutThree      string `json:"about_three"`
	AboutFour       string `json:"about_four"`
	AboutFive       string `json:"about_five"`
	AboutSix        string `json:"about_six"`
	AboutImageOne   int    `json:"about_image_one"`
	AboutImageTwo   int    `json:"about_image_two"`
	AboutImageThree int    `json:"about_image_three"`
	AboutImageFour  int    `json:"about_image_four"`

	NewProductOne string `json:"new_product_one"`
	ProductOne    string `json:"product_one"`

	ContactOne   string `json:"contact_one"`
	ContactTwo   string `json:"contact_two"`
	ContactThree string `json:"contact_three"`
	ContactFour  string `json:"contact_four"`
	ContactFive  string `json:"contact_five"`
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
