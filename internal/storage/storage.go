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
	OrderItems []OrderItem `json:"order_items"`
}

type Styles struct {
	AboutUs  AboutUsStyle  `json:"about_us"`
	Footer   FooterStyle   `json:"footer"`
	Header   HeaderStyle   `json:"header"`
	Products ProductsStyle `json:"products"`
}

type AboutUsStyle struct {
	BackgroundColor string `json:"background_color"`
	Content         string `json:"content"`
	FontSize        string `json:"font_size"`
	Height          string `json:"height"`
	Width           string `json:"width"`
}

type FooterStyle struct {
	BackgroundColor string `json:"background_color"`
	Content         string `json:"content"`
	FontSize        string `json:"font_size"`
	Height          string `json:"height"`
	Width           string `json:"width"`
}

type HeaderStyle struct {
	BackgroundColor string `json:"background_color"`
	Content         string `json:"content"`
	FontSize        string `json:"font_size"`
	Height          string `json:"height"`
	Width           string `json:"width"`
}

type ProductsStyle struct {
	BackgroundColor string `json:"background_color"`
	Content         string `json:"content"`
	FontSize        string `json:"font_size"`
	Height          string `json:"height"`
	Width           string `json:"width"`
}
