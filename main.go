package main

import (
	"e-commerce-api/config"
	_productData "e-commerce-api/feature/product/data"
	_productHandler "e-commerce-api/feature/product/handler"
	_productService "e-commerce-api/feature/product/service"
	"e-commerce-api/feature/users/data"
	"e-commerce-api/feature/users/handler"
	"e-commerce-api/feature/users/services"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := config.InitDB(*cfg)

	// panggil fungsi Migrate untuk buat table baru di database
	config.Migrate(db)

	v := validator.New()

	productData := _productData.New(db)
	productService := _productService.New(productData, v)
	productHandler := _productHandler.New(productService)

	userData := data.New(db)
	userSrv := services.New(userData)
	userHdl := handler.New(userSrv)

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, error=${error}\n",
	}))

	e.POST("/register", userHdl.Register())
	e.POST("/login", userHdl.Login())
	e.GET("/users", userHdl.Profile(), middleware.JWT([]byte(config.JWT_KEY)))
	e.PUT("/users", userHdl.Update(), middleware.JWT([]byte(config.JWT_KEY)))
	e.DELETE("/users", userHdl.Delete(), middleware.JWT([]byte(config.JWT_KEY)))
	e.POST("/products", productHandler.Add(), middleware.JWT(config.JWT_KEY))
	e.GET("/products", productHandler.GetAll())
	e.GET("/products/:product_id", productHandler.GetByID())

	e.POST("/products", productHandler.Add(), middleware.JWT(config.JWT_KEY))
	if err := e.Start(":8000"); err != nil {
		log.Println(err.Error())
	}

}
