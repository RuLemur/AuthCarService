package main

import (
	_ "AuthCarService/docs"
	"AuthCarService/internal/config"
	"AuthCarService/internal/pkg"
)

func main() {
	app := pkg.NewApp()
	cfg := config.ReadConfig()
	app.RunApp(cfg)
}
