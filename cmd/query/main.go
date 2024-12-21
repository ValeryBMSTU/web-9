package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres1"
	dbname   = "mydatabase"
)

type Handlers struct {
	dbProvider DatabaseProvider
}

type DatabaseProvider struct {
	db *sql.DB
}

func (h *Handlers) Handler(c echo.Context) error {
	name := c.QueryParam("name")
	if name != "" {
		err := h.dbProvider.InsertHello(name)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "Hello, "+name+"!")
	} else {
		msg, err := h.dbProvider.SelectHello()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "Hello, "+msg+"!")
	}
}

func (dp *DatabaseProvider) SelectHello() (string, error) {
	var msg string

	row := dp.db.QueryRow("SELECT message FROM hello ORDER BY RANDOM() LIMIT 1")
	err := row.Scan(&msg)
	if err != nil {
		return "", err
	}

	return msg, nil
}

func (dp *DatabaseProvider) InsertHello(msg string) error {
	query := `SELECT EXISTS(SELECT 1 FROM hello WHERE message = $1);`
	var exists bool
	err := dp.db.QueryRow(query, msg).Scan(&exists)

	if err != nil {
		return fmt.Errorf("ошибка при запросе: %v", err)
	}
	if !exists {
		_, err := dp.db.Exec("INSERT INTO hello (message) VALUES ($1)", msg)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	address := flag.String("address", "127.0.0.1:9000", "адрес для запуска сервера")
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
	e.GET("/api/user", h.Handler)

	e.Logger.Fatal(e.Start(*address))
}
