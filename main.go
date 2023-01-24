package main

import (
	"e-commerce-api/config"
	_productData "e-commerce-api/feature/product/data"
	_productHandler "e-commerce-api/feature/product/handler"
	_productService "e-commerce-api/feature/product/service"
	"log"

	"github.com/go-playground/validator/v10"
	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	c := config.InitConfig()
	db := config.InitDB(*c)
	config.Migrate(db)

	v := validator.New()

	productData := _productData.NewProductData(db)
	productService := _productService.NewProductService(productData, v)
	productHandler := _productHandler.NewProductHandler(productService)

	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom}, method=${method}, uri=${uri}, status=${status}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	e.POST("/products", productHandler.Add(), echojwt.JWT(config.JWT_KEY))

	if err := e.Start(":8000"); err != nil {
		log.Fatal("failed to run server")
	}
}
