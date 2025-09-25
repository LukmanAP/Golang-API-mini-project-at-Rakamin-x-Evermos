package http

import (
	"project-evermos/internal/config"
	authHandler "project-evermos/internal/todo/handler/auth"
	productHandler "project-evermos/internal/todo/handler/product"
	tokoHandler "project-evermos/internal/todo/handler/toko"
	usersHandler "project-evermos/internal/todo/handler/users"
	transactionHandler "project-evermos/internal/todo/handler/transaction"
	authRepo "project-evermos/internal/todo/repository/auth"
	productRepo "project-evermos/internal/todo/repository/product"
	storeRepo "project-evermos/internal/todo/repository/toko"
	usersRepo "project-evermos/internal/todo/repository/users"
	transactionRepo "project-evermos/internal/todo/repository/transaction"
	authService "project-evermos/internal/todo/service/auth"
	productService "project-evermos/internal/todo/service/product"
	tokoService "project-evermos/internal/todo/service/toko"
	usersService "project-evermos/internal/todo/service/users"
	transactionService "project-evermos/internal/todo/service/transaction"
	// Address service imports
	addressHandler "project-evermos/internal/todo/handler/address"
	addressRepo "project-evermos/internal/todo/repository/address"
	addressService "project-evermos/internal/todo/service/address"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// RegisterRoutes registers HTTP routes for the application.
// This keeps the router setup centralized.
func RegisterRoutes(app *fiber.App, gdb *gorm.DB, cfg *config.Config) {
	// Healthcheck endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Auth module wiring
	repo := authRepo.NewRepository(gdb)
	storeR := storeRepo.NewRepository(gdb)
	service := authService.NewService(repo, storeR, cfg)
	h := authHandler.NewHandler(service)

	grp := app.Group("/auth")
	grp.Post("/login", h.Login)
	grp.Post("/register", h.Register)

	// Toko module wiring
	tSvc := tokoService.NewService(storeR)
	tH := tokoHandler.NewHandler(tSvc)

	// Protected Toko endpoints (require JWT)
	jwtMW := tokoHandler.JWTMiddleware(cfg.JWTSecret)
	// Register static route before parameterized route to avoid capture as :id_toko
	app.Get("/toko/my", jwtMW, tH.GetMy)

	// Public Toko endpoints
	app.Get("/toko", tH.List)
	app.Get("/toko/:id_toko", tH.GetByID)

	// Users module wiring
	uRepo := usersRepo.NewRepository(gdb)
	uService := usersService.NewService(uRepo)
	uHandler := usersHandler.NewHandler(uService)

	// Protected user endpoints using JWTMiddleware that reads token header
	uJWT := usersHandler.JWTMiddleware(cfg.JWTSecret)
	app.Get("/user", uJWT, uHandler.GetProfile)
	app.Put("/user", uJWT, uHandler.UpdateProfile)

	// Alamat Kirim endpoints
	app.Get("/user/alamat", uJWT, uHandler.ListAlamat)
	app.Post("/user/alamat", uJWT, uHandler.CreateAlamat)
	app.Get("/user/alamat/:id", uJWT, uHandler.GetAlamat)
	app.Put("/user/alamat/:id", uJWT, uHandler.UpdateAlamat)
	app.Delete("/user/alamat/:id", uJWT, uHandler.DeleteAlamat)

	// Product module wiring
	pRepo := productRepo.NewRepository(gdb)
	pService := productService.NewService(pRepo, cfg.BaseFileURL)
	pHandler := productHandler.NewHandler(pService, storeR, cfg)

	// JWT for protected product endpoints (supports 'token' header and Authorization: Bearer)
	pJWT := usersHandler.JWTMiddleware(cfg.JWTSecret)

	// Public Product endpoints
	app.Get("/product", pHandler.List)
	app.Get("/product/:id", pHandler.GetByID)

	// Protected Product endpoints
	app.Post("/product", pJWT, pHandler.Create)
	app.Put("/product/:id", pJWT, pHandler.Update)
	app.Delete("/product/:id", pJWT, pHandler.Delete)

	// Address (Province/City) public endpoints using EMSIFA
	addrRepo := addressRepo.NewRepository(cfg.EMSIFABase, cfg.HTTPTimeoutMS, cfg.HTTPRetry)
	addrSvc := addressService.NewService(addrRepo, time.Duration(cfg.CacheTTLSeconds)*time.Second)
	addrH := addressHandler.NewHandler(addrSvc)
	app.Get("/provcity/listprovincies", addrH.ListProvinces)
	app.Get("/provcity/listcities/:prov_id", addrH.ListCities)
	app.Get("/provcity/detailprovince/:prov_id", addrH.DetailProvince)
	app.Get("/provcity/detailcity/:city_id", addrH.DetailCity)

	// Transaction module wiring
	trxRepo := transactionRepo.NewRepository(gdb)
	trxService := transactionService.NewService(trxRepo)
	trxHandler := transactionHandler.NewHandler(trxService)

	trxJWT := usersHandler.JWTMiddleware(cfg.JWTSecret)
	app.Get("/trx", trxJWT, trxHandler.List)
	app.Get("/trx/:id", trxJWT, trxHandler.GetByID)
	app.Post("/trx", trxJWT, trxHandler.Create)
}
