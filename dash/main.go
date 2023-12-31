package main

import (
	"fmt"
	"github.com/caarlos0/env/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	AppNAme    string `env:"APPNAME,required"`
	Version    string `env:"VERSION,required"`
	ApiVersion string `env:"APIVERSION,required"`
	Listen     string `env:"LISTEN,required"`
}

func main() {
	err := godotenv.Load()
	cfg := Config{} // 👈 new instance of `Config`

	err = env.Parse(&cfg) // 👈 Parse environment variables into `Config`
	if err != nil {
		fmt.Printf("unable to parse ennvironment variables: %e\n", err)
	}
	if err != nil {
		fmt.Println(err)
	}

	app := fiber.New()
	app.Get("/metrics", monitor.New(monitor.Config{Title: "KREDS Dash Report Monitor"}))
	defer func(app *fiber.App) {
		err := app.Shutdown()
		if err != nil {
			fmt.Println("Impossible to shutdown the application")
		}
	}(app)
	file, err := os.OpenFile("/logs/dash.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:        "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat:    "02/Jan/2006:15:04:05",
		TimeZone:      "Europe/Rome",
		TimeInterval:  500 * time.Millisecond,
		Output:        file,
		DisableColors: false,
	}))

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	// Static file server
	app.Static("/", "./public/")

	// => http://localhost:3000/hello.txt
	// => http://localhost:3000/gopher.gif

	// Start the server

	if err := app.Listen(cfg.Listen); err != nil {
		panic(err)
	}
	fmt.Printf("%v started\n", cfg.AppNAme)
}
