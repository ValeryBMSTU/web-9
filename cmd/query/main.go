package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/api/user", getHandler)

	e.Logger.Fatal(e.Start(":8000"))
}

func getHandler(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		name = "Stranger"
	}
	return c.String(http.StatusOK, "Hello, "+name+"!")
}
