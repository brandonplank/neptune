package main

import (
	"brandonplank.org/neptune/database"
	"brandonplank.org/neptune/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"log"
	"strconv"
)

const Port = 8080

func setupRoutes(app *fiber.App) {
	app.Use(
		cors.New(cors.Config{
			AllowHeaders:     "Origin, Content-Type, Accept, Access-Control-Allow-Origin",
			AllowCredentials: true,
		}),
		logger.New(logger.Config{
			Format:     "${time} [${method}]->${status} Latency->${latency} - ${path} | ${error}\n",
			TimeFormat: "2006/01/02 15:04:05",
		}),
		cors.New(cors.Config{
			AllowCredentials: true,
		}),
		func(ctx *fiber.Ctx) error {
			ctx.Append("Access-Control-Allow-Origin", "*")
			ctx.Append("Developer", "Brandon Plank")
			ctx.Append("License", "GNU Affero General Public License v3.0")
			return ctx.Next()
		},
	)

	app.Static("/", "./Public")

	app.Get("/", routes.Home)

	// API

	v1 := app.Group("/v1")

	v1.Post("/register", routes.Register)
	v1.Post("/adminRegister", routes.AdminRegister)
	v1.Post("/login", routes.Login)
	v1.Get("/logout", routes.Logout)

	v1.Post("/id/:name", routes.Id)
	v1.Post("/isOut/:name", routes.IsOut)
	v1.Post("/addSchool", routes.AddSchool)

	v1.Get("/GetCSV", routes.GetCSV)
	v1.Get("/GetAdminCSV", routes.GetAdminCSV)
	v1.Get("/classroom.csv", routes.CSVFile)
	v1.Get("/user", routes.User)
	v1.Post("/search", routes.AdminSearchStudent)
	v1.Post("/search/:name", routes.AdminSearchStudent)
	v1.Get("/getSchool", routes.GetSchool)
	v1.Get("/getSchools", routes.GetSchools)
}

// main SignX program entry point
func main() {
	log.Println("Starting Neptune")

	// TODO: Sentry

	// MySQL
	database.Connect()

	engine := html.New("./resources/views", ".html")
	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	// Setup all the server routes
	setupRoutes(router)
	log.Println("Finished setting up routes")

	log.Println("Started Neptune on port", strconv.Itoa(Port))
	log.Fatalln(router.Listen(":" + strconv.Itoa(Port)))
}
