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

	e.GET("/count", getHandler)
	e.POST("/count", postHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

type Counter struct {
	Count int `json:"count"`
}

func getHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]int{"count": count})
}

func postHandler(c echo.Context) error {
	var counter Counter
	if err := c.Bind(&counter); err != nil {
		return err
	}

	count += counter.Count

	return c.JSON(http.StatusOK, counter)
}
