package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "r0ckwe11"
	password = "postgres1"
	dbname   = "mydatabase"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func (h *Handlers) PostCounter(c echo.Context) error {
	a, err := strconv.Atoi(c.FormValue("count"))
	if err != nil {
		c.Logger().Error(err)
		return c.String(http.StatusBadRequest, "это не число")
	}
	err = h.dbProvider.AddCount(a)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, "OK!")
}

func (h *Handlers) GetCounter(c echo.Context) error {
	value, err := h.dbProvider.GetCount()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, strconv.Itoa(value))
}

func (dp *DatabaseProvider) GetCount() (int, error) {
	var value int

	row := dp.db.QueryRow("SELECT COALESCE(count, 0) FROM count WHERE name=$1", "key1")
	err := row.Scan(&value)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (dp *DatabaseProvider) AddCount(a int) error {
	_, err := dp.db.Exec("INSERT INTO count (name, count) VALUES ($2, $1) ON CONFLICT (name) DO UPDATE SET count = count.count + $1", a, "key1")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	address := flag.String("address", "127.0.0.1:3333", "адрес для запуска сервера")
	flag.Parse()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dp := DatabaseProvider{db: db}

	h := Handlers{dbProvider: dp}

	e := echo.New()
	e.POST("/count", h.PostCounter)
	e.GET("/count", h.GetCounter)
	e.Logger.Fatal(e.Start(*address))
}
