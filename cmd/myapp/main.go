package main

import (
	"fmt"
	"log"
	"os"
	"payment/internal/db"
	"payment/internal/handler"
	"payment/internal/middleware"
	"payment/internal/repository"
	"payment/internal/service"
	"payment/internal/utils"
	"payment/internal/xenditpay"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if os.Getenv("ENV") != "production" {
	err := godotenv.Load("../../.env") 
		if err != nil {
		log.Fatal("Error loading .env file")
		}
	}


	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	db.Connect(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	fmt.Println(os.Getenv("SECRET_KEY"))

	tokenMaker := utils.NewJWTMaker(os.Getenv("SECRET_KEY"))

	authRepo := repository.NewAuthenticationRepository(db.GetDB())
	authService := service.NewAuthenticationService(authRepo, tokenMaker)
	authHandler := handler.NewAuthenticationHandler(authService)

	subsPayRepo := repository.NewSubscriptionPaymentRepository(db.GetDB())
	subsPayService := service.NewSubcriptionService(subsPayRepo)
	subsPayHandler := handler.NewSubscriptionPaymentHandler(subsPayService)

	packageRepo := repository.NewPackageRepository(db.GetDB())
	packageService := service.NewPackageService(packageRepo)
	packageHandler := handler.NewPackageHandler(packageService)

	xenditHandler := xenditpay.NewPaymentHandler(subsPayRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(tokenMaker))

	r.POST("/api/authentication/register", authHandler.Register)
	r.POST("/api/authentication/login", authHandler.Login)

	r.POST("/api/payment/webhook", xenditHandler.PaymentCallbackHandler)
	r.POST("/api/invoice/payment/webhook", xenditHandler.InvoiceCallbackHandler)
	r.POST("/api/payment/simulate/:xenditId", xenditHandler.SimulatePaymentHandler)

	r.GET("/api/packages", packageHandler.FindAllPackages)

	protected.GET("/package/:id", packageHandler.GetPackageById)

	protected.GET("/current", authHandler.GetDataUser)

	protected.POST("/subscribe", subsPayHandler.CreateSubscriptionPayment)
	protected.POST("/expired", subsPayHandler.CheckPaymentExpired)
	protected.GET("/payment/status/:xenditId", subsPayHandler.CheckPaymentStatusById)
	protected.GET("/user/package/active", subsPayHandler.GetUserPackageAccess)
	protected.GET("/user/subscription/active", subsPayHandler.CheckSubscriptionActive)

	// protected.POST("/customer", xenditpay.CreateCustomerHandler)

	protected.GET("/payment/data/:xenditId", xenditHandler.GetPaymentRequestByIdHandler)
	

	r.Run(":8080")

	r.Run(":" + port)

}
