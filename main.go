package main

import (
	"fmt"
	"net/http"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func test(c echo.Context) error {
	return c.String(http.StatusOK, "Wow this actually worked")
}

func main() {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", test)
	fmt.Println("Listening for requests on port :8080")

	e.Logger.Fatal(e.Start(":8080"))
}