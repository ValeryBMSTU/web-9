package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "sandbox9"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

// Обработчики HTTP-запросов
func (h *Handlers) GetQuery(c echo.Context) error {
	name := c.QueryParam("name")

	if name == "" {
		return c.String(http.StatusBadRequest, "Не введен параметр!")
	}

	test, err := h.dbProvider.SelectQuery(name)
	if !test && err == nil {
		return c.String(http.StatusBadRequest, "Запись не добавлена в БД!")
	} else if !test && err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Hello "+name+"!")
}

func (h *Handlers) PostQuery(c echo.Context) error {
	name := c.QueryParam("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Не введен параметр!")
	}

	test, err := h.dbProvider.SelectQuery(name)
	if test && err == nil {
		return c.String(http.StatusBadRequest, "Запись уже добавлена в БД!")
	}

	err = h.dbProvider.InsertQuery(name)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, "Добавили запись!")
}

// Методы для работы с базой данных
func (dp *DatabaseProvider) SelectQuery(msg string) (bool, error) {
	var rec string

	row := dp.db.QueryRow("SELECT record FROM query WHERE record = ($1)", msg)
	err := row.Scan(&rec)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (dp *DatabaseProvider) InsertQuery(msg string) error {
	_, err := dp.db.Exec("INSERT INTO query (record) VALUES ($1)", msg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	// Считываем аргументы командной строки
	address := flag.String("address", "127.0.0.1:8083", "адрес для запуска сервера")
	flag.Parse()

	// Формирование строки подключения для postgres
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"dbname=%s sslmode=disable",
		host, port, user, dbname)

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

	e.GET("/query", h.GetQuery)
	e.POST("/query1", h.PostQuery)

	e.Logger.Fatal(e.Start(*address))
}
