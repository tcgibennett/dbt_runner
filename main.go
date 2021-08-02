package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"stitchdata.com/dbt/runner/rest"
)

var e *echo.Echo

func main() {
	e = echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/compiled/sql/:account_id", rest.GetCompiledSQL)
	e.GET("/project/:account_id/vars", rest.GetProjectVars)
	e.GET("/profiles/:account_id", rest.GetProfiles)
	e.POST("/compile/:account_id", rest.PostCompile)
	e.POST("/execute/:account_id", rest.PostExecute)
	e.POST("/project/:account_id", rest.CreateProject)
	e.POST("/account/:account_id", rest.CreateAccount)
	fmt.Println("Listening for requests on port :8080")

	e.Logger.Fatal(e.Start(":8080"))
}
