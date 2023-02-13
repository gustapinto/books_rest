package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gustapinto/books_rest/go_gin_sqlx/docs"
	"github.com/gustapinto/books_rest/go_gin_sqlx/internal/config"
	"github.com/gustapinto/books_rest/go_gin_sqlx/pkg/controller"
	"github.com/gustapinto/books_rest/go_gin_sqlx/pkg/middleware"
	"github.com/gustapinto/books_rest/go_gin_sqlx/pkg/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	swaggoFiles "github.com/swaggo/files"
	swaggoGin "github.com/swaggo/gin-swagger"
)

// @title Books REST - GO + Gin + SQLX
// @version dev
// @description Just a simple book management API written using Go, Gin and SQLX

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the authorized JWT token
func main() {
	docs.SwaggerInfo.BasePath = "/api"

	db, err := sqlx.Connect("pgx", config.DB_DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	pingController := controller.NewPingController()
	userController := controller.NewUserController(userRepository)
	authController := controller.NewAuthController(userRepository)

	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/auth", authController.Login)
		api.GET("/ping", pingController.Pong)
		api.GET("/swagger/*any", swaggoGin.WrapHandler(swaggoFiles.Handler))

		user := api.Group("/user")
		{
			user.GET("", userController.All).Use(middleware.Auth)
			user.GET(":userId", userController.Find).Use(middleware.Auth)
			user.POST("", userController.Create)
			user.PUT(":userId", userController.Update).Use(middleware.Auth)
			user.DELETE(":userId", userController.Delete).Use(middleware.Auth)
		}
	}

	if err := router.Run(config.APP_ADDR); err != nil {
		log.Fatal(err)
	}
}
