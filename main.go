package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"stitchdata.com/dbt/runner/dbt"
)

var e *echo.Echo
var d *dbt.Dbt

func test(c echo.Context) error {
	return c.String(http.StatusOK, "Wow this actually worked... new version!")
}

func GetProfiles(c echo.Context) error {
	result := "{\"hello\":\"world\"}"
	fmt.Println(c.Param("account_id"))

	return c.JSON(http.StatusOK, result)
}

func PostCompile(c echo.Context) error {
	accountid := c.Param("account_id")
	d = dbt.New(accountid)
	//response := d.Compile("--project-dir", "/"+d.Path+"/dbt_test")
	response := d.Compile("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)

	return c.String(http.StatusOK, response.Message)
}

func main() {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", test)
	e.GET("/profiles/:account_id", GetProfiles)
	e.POST("/compile/:account_id", PostCompile)
	fmt.Println("Listening for requests on port :8080")

	e.Logger.Fatal(e.Start(":8080"))
}
