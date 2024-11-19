package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var count int = 0

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/hello", getHandler)

	e.Logger.Fatal(e.Start(":2000"))
}

func getHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, Web!")
}
