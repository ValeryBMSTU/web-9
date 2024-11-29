package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "vagab"
	password = "Vagab78901"
	dbname   = "LAB8"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// обработчики http-запросов
func (h *Handlers) GetCounter(c echo.Context) error {
	msg, err := h.dbProvider.SelectCounter()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "counter: "+strconv.Itoa(msg))
}

func (h *Handlers) PostCounter(c echo.Context) error {
	input := struct {
		Msg int `json:"msg"`
	}{}

	err := c.Bind(&input)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	err = h.dbProvider.UpdateCounter(input.Msg)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "counter changed")
}

// методы для работы с базой данных
func (dp *DatabaseProvider) SelectCounter() (int, error) {
	var msg int

	row := dp.db.QueryRow("SELECT number FROM count WHERE id_number = 1")
	err := row.Scan(&msg)
	if err != nil {
		return -1, err
	}

	return msg, nil
}

func (dp *DatabaseProvider) UpdateCounter(msg int) error {
	_, err := dp.db.Exec("UPDATE count SET number = number + $1 WHERE id_number = 1", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Считываем аргументы командной строки
	address := flag.String("address", "127.0.0.1:8082", "адрес для запуска сервера")
	flag.Parse()

	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Создание соединения с сервером postgres
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создаем провайдер для БД с набором методов
	dp := DatabaseProvider{db: db}
	// Создаем экземпляр структуры с набором обработчиков
	h := Handlers{dbProvider: dp}

	e := echo.New()

	e.Use(middleware.Logger())

	e.GET("/counter", h.GetCounter)
	e.POST("/counter", h.PostCounter)

	e.Logger.Fatal(e.Start(*address))
}
