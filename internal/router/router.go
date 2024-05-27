package router

import (
	"ForeignKey/internal/email"
	"ForeignKey/internal/http-server/handlers/admin"
	"ForeignKey/internal/http-server/handlers/cart"
	"ForeignKey/internal/http-server/handlers/customer"
	img "ForeignKey/internal/http-server/handlers/image"
	"ForeignKey/internal/http-server/handlers/order"
	"ForeignKey/internal/http-server/handlers/product"
	"ForeignKey/internal/http-server/handlers/website"
	mwLogger "ForeignKey/internal/http-server/middleware/logger"
	"ForeignKey/internal/image"
	"ForeignKey/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
)

func New(storage *sqlite.Storage, imageSaver *image.Image, emailSender *email.Email, log *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*", "https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// TODO: указать статусы для всех ответов
	// TODO: поменять некоторык еррор логги на инфо
	// TODO: поменять некоторые пост запросы на патч и делейт
	// TODO: вход по вк
	// handlers
	router.Post("/api/admin/sign-up", admin.NewSignUp(storage, log))
	router.Post("/api/admin/sign-in", admin.NewSignIn(storage, log))

	router.Post("/api/image/upload", img.NewUpload(imageSaver, storage, log))
	router.Get("/api/image/download/{id}", img.NewDownload(imageSaver, storage, log))

	router.Post("/api/website/create", website.NewCreate(storage, log))
	router.Get("/api/website/aliases", website.NewGetAliases(storage, log))
	router.Delete("/api/website/delete/{alias}", website.NewDelete(storage, log))

	router.Post("/api/product/create", product.NewCreate(storage, log))
	router.Get("/api/product/get-by-alias/{alias}", product.NewGetByAlias(storage, log))

	router.Post("/api/customer/sign-up", customer.NewSignUp(storage, log))
	router.Post("/api/customer/sign-in", customer.NewSignIn(storage, log))

	router.Post("/api/cart/add", cart.NewAdd(storage, log))
	router.Patch("/api/cart/change-count", cart.NewChangeCount(storage, log))
	router.Get("/api/cart/get", cart.NewGet(storage, log))

	// TODO: добавить сохранение цены товара на момент заказа
	// TODO: добавить статус заказа, уведы на почту
	router.Post("/api/order/make", order.NewMakeOrder(storage, emailSender, log))
	router.Get("/api/order/get", order.NewGet(storage, log))

	return router
}
