package main

import (
	"go_bedu/config"
	"go_bedu/routes"

	"github.com/labstack/echo"
)

func main() {
	db := config.InitDB()
	e := echo.New()

	routes.NewRoute(e, db)

	e.Logger.Fatal(e.Start(":8080"))
}
