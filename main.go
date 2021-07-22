package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var e *echo.Echo

func test(c echo.Context) error {
	return c.String(http.StatusOK, "Wow this actually worked... new version!")
}

func GetProfiles(c echo.Context) error {
	result := "{\"hello\":\"world\"}"
	fmt.Println(c.Param("account_id"))
	
	return c.JSON(http.StatusOK, result)
}

func main() {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", test)
	e.GET("/profiles/:account_id", GetProfiles)
	fmt.Println("Listening for requests on port :8080")

	e.Logger.Fatal(e.Start(":8080"))
}
