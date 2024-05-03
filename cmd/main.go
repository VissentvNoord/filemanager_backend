package main

import (
	"example.com/myproject/pkg/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	startServer()

	//storagemanager.Connect()
}

func startServer() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	controllers.Initialize(e)

	e.Logger.Fatal(e.Start(":42069"))
}
