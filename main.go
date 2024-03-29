package main

import (
	"brandonplank.org/neptune/database"
	neptuneembed "brandonplank.org/neptune/embed"
	"brandonplank.org/neptune/global"
	"brandonplank.org/neptune/routes"
	"fmt"
	. "github.com/crypticplank/israilway"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
	"github.com/mileusna/crontab"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	_ "time/tzdata"
)

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
		favicon.New(favicon.Config{
			FileSystem: http.FS(neptuneembed.Content),
			File:       "public/favicon.ico",
		}),
		func(ctx *fiber.Ctx) error {
			ctx.Append("Access-Control-Allow-Origin", "*")
			ctx.Append("Developer", "Brandon Plank")
			ctx.Append("License", "GNU Affero General Public License v3.0")
			return ctx.Next()
		},
	)
	app.Use("/assets", filesystem.New(filesystem.Config{
		Root:       http.FS(neptuneembed.Content),
		PathPrefix: "public/assets",
		Browse:     false,
	}))

	app.Get("/", routes.Home)
	app.Get("/admin", routes.Admin)
	app.Get("/passwordReset", routes.PasswordChangeSite)

	// API

	v1 := app.Group("/v1")

	v1.Post("/register", routes.Register)
	v1.Post("/adminRegister", routes.AdminRegister)
	v1.Post("/login", routes.Login)
	v1.Get("/logout", routes.Logout)

	v1.Post("/id/:name", routes.Id)
	v1.Post("/isOut/:name", routes.IsOut)
	v1.Post("/addSchool", routes.AddSchool)
	v1.Post("/changePassword", routes.PasswordChange)

	v1.Get("/GetCSV", routes.GetCSV)
	v1.Get("/GetAdminCSV", routes.GetAdminCSV)
	v1.Get("/admin.csv", routes.GetAdminCSVFile)
	v1.Get("/classroom.csv", routes.CSVFile)
	v1.Get("/user", routes.User)
	v1.Post("/search", routes.AdminSearchStudent)
	v1.Post("/search/:name", routes.AdminSearchStudent)
	v1.Get("/getSchool", routes.GetSchool)
	v1.Get("/getSchools", routes.GetSchools)
	v1.Get("/getUserPermissionLevel", routes.GetUserPermissionLevel)
}

func init() {
	if IsRailway() {
		log.Println("Running on railway")
		global.Port = 443
	} else {
		log.Println("Reading from .env")
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}
}

// main SignX program entry point
func main() {
	// Pre Start
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Println("Could not set time to New York")
	}
	time.Local = loc

	global.EmailPassword, _ = os.LookupEnv("EMAIL_PASSWORD")
	log.Println("Starting Neptune")

	// Sentry
	err = sentry.Init(sentry.ClientOptions{
		Dsn:   "https://********************.ingest.sentry.io/6340501",
		Debug: false,
	})
	if err != nil {
		log.Fatalln(fmt.Sprintf("Sentry failed:%s", err.Error()))
	}

	defer sentry.Flush(2 * time.Second)
	defer sentry.Recover()

	log.Println("Started Sentry")

	// MySQL
	database.Connect()

	// Page rendering
	engine := html.NewFileSystem(http.FS(neptuneembed.ViewContent), ".html")
	router := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		Views:                 engine,
	})

	// Cronjobs
	ctab := crontab.New()

	ctab.MustAddJob("5 15 * * 1-5", func() { // 03:05 PM every weekday
		global.EmailUsers()
	})

	ctab.MustAddJob("0 0 * * 1-5", func() { // 12:00 AM every weekday
		global.CleanStudents()
	})

	// Setup all the server routes
	setupRoutes(router)
	log.Println("Finished setting up routes")

	log.Println("Started Neptune on port", strconv.Itoa(global.Port))
	log.Fatalln(router.Listen(":" + strconv.Itoa(global.Port)))
}
