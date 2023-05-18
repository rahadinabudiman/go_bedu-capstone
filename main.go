package main

import (
	"go_bedu/config"
	"go_bedu/routes"
)

func main() {
	config.InitDB()
	e := routes.New()

	e.Logger.Fatal(e.Start(":8000"))
}
