package router

import (
	"ForeignKey/internal/http-server/handlers/admin"
	"ForeignKey/internal/http-server/handlers/cart"
	"ForeignKey/internal/http-server/handlers/customer"
	img "ForeignKey/internal/http-server/handlers/image"
	"ForeignKey/internal/http-server/handlers/order"
	"ForeignKey/internal/http-server/handlers/product"
	"ForeignKey/internal/http-server/handlers/website"
	mwLogger "ForeignKey/internal/http-server/middleware/logger"
	"ForeignKey/internal/image"
	"ForeignKey/internal/notification_client"
	"ForeignKey/internal/storage/sqlite"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
)

func New(
	storage *sqlite.Storage,
	imageSaver *image.Image,
	nc *notification_client.NotificationClient,
	log *slog.Logger,
) *chi.Mux {
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*", "https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	// TODO: вход по вк
	// handlers
	router.Post("/api/admin/sign-up", admin.NewSignUp(storage, log))
	router.Post("/api/admin/sign-in", admin.NewSignIn(storage, log))
	router.Get("/api/admin/get-profile", admin.NewGetProfile(storage, log))
	router.Put("/api/admin/update-profile", admin.NewUpdateProfile(storage, log))
	router.Patch("/api/admin/set-email-notification", admin.NewEmailNotification(storage, log))
	router.Patch("/api/admin/set-telegram-notification", admin.NewTelegramNotification(storage, log))

	router.Post("/api/image/upload", img.NewUpload(imageSaver, storage, log))
	router.Get("/api/image/download/{id}", img.NewDownload(imageSaver, storage, log))

	router.Post("/api/website/create", website.NewCreate(storage, log))
	router.Get("/api/website/aliases", website.NewGetAliases(storage, log))
	router.Delete("/api/website/delete/{alias}", website.NewDelete(storage, log))
	router.Patch("/api/website/set-style", website.NewSetStyle(storage, log))
	router.Get("/api/website/get-style/{alias}", website.NewGetStyle(storage, log))

	router.Post("/api/product/create", product.NewCreate(storage, log))
	router.Get("/api/product/get-by-alias/{alias}", product.NewGetByAlias(storage, log))
	router.Get("/api/product/get-all-by-alias/{alias}", product.NewGetAllByAlias(storage, log))
	router.Patch("/api/product/set-active", product.NewSetActive(storage, log))
	router.Put("/api/product/update", product.NewUpdate(storage, log))
	router.Delete("/api/product/delete/{id}", product.NewDelete(storage, log))

	router.Post("/api/customer/sign-up", customer.NewSignUp(storage, log))
	router.Post("/api/customer/sign-in", customer.NewSignIn(storage, log))
	router.Get("/api/customer/get-by-alias/{alias}", customer.NewGetByAlias(storage, log))
	router.Get("/api/customer/get-profile", customer.NewGetProfile(storage, log))
	router.Put("/api/customer/update-profile", customer.NewUpdateProfile(storage, log))
	router.Patch("/api/customer/set-email-notification", customer.NewEmailNotification(storage, log))
	router.Patch("/api/customer/set-telegram-notification", customer.NewTelegramNotification(storage, log))

	router.Post("/api/cart/add", cart.NewAdd(storage, log))
	router.Patch("/api/cart/change-count", cart.NewChangeCount(storage, log))
	router.Get("/api/cart/get", cart.NewGet(storage, log))

	router.Post("/api/order/make", order.NewMakeOrder(storage, nc, log))
	router.Get("/api/order/get", order.NewGet(storage, log))
	router.Get("/api/order/get-by-alias/{alias}", order.NewGetByAlias(storage, log))
	router.Get("/api/order/get-completed/{alias}", order.NewGetCompletedOrders(storage, log))
	router.Patch("/api/order/set-status", order.NewSetStatus(storage, nc, log))

	return router
}
