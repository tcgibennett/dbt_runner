package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"stitchdata.com/dbt/runner/dbt"
	"stitchdata.com/dbt/runner/messages"
	"stitchdata.com/dbt/runner/rest"
)

var e *echo.Echo
var d *dbt.Dbt

func test(c echo.Context) error {
	return c.String(http.StatusOK, "Wow this actually worked... new version!")
}

func GetProfiles(c echo.Context) error {
	response := messages.Response{
		Status:  "Success",
		Message: "Hello World",
	}
	result, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c.Param("account_id"))
	fmt.Println(string(result))
	return c.JSON(http.StatusOK, string(result))
}

func PostCompile(c echo.Context) error {
	accountid := c.Param("account_id")
	d = dbt.New(accountid)
	//response := d.Compile("--project-dir", "/"+d.Path+"/dbt_test")
	response := d.Compile("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)

	if response.Status == "Success" {
		return c.String(http.StatusOK, response.Message)
	} else {
		return c.String(http.StatusNotAcceptable, response.Message)
	}
}

func PostExecute(c echo.Context) error {
	accountid := c.Param("account_id")
	d = dbt.New(accountid)

	response := d.Execute("--project-dir", "/"+accountid+"/dbt_test", "--profiles-dir", "/"+accountid)

	if response.Status == "Success" {
		return c.String(http.StatusOK, response.Message)
	} else {
		return c.String(http.StatusNotAcceptable, response.Message)
	}
}

func main() {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/test", test)
	e.GET("/compiled/sql/:account_id", rest.GetCompiledSQL)
	e.GET("/profiles/:account_id", GetProfiles)
	e.POST("/compile/:account_id", PostCompile)
	e.POST("/execute/:account_id", PostExecute)
	e.POST("/create/project/:account_id", rest.CreateProject)
	fmt.Println("Listening for requests on port :8080")

	e.Logger.Fatal(e.Start(":8080"))
}
