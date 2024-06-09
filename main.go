package main

import (
	"fmt"
	"log"
	"piotrek813/wayback-photos/mimetype"
	"piotrek813/wayback-photos/wayback"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Form struct {
	Website   string   `form:"website"`
	Limit     int      `form:"limit"`
	Mimetype  []string `form:"mimetype"`
	ResumeKey string   `form:"resumeKey"`
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":    "Wayback Photos",
			"Mimetype": []string{mimetype.JPEG, mimetype.GIF, mimetype.PNG},
		})
	})

	app.Post("/getPhotos", func(c *fiber.Ctx) error {
		f := new(Form)

		start := time.Now()
		fmt.Printf("start: %v\n", start)
		if err := c.BodyParser(f); err != nil {
			fmt.Printf("err: %v\n", err)
			return c.SendString("Nie działa :(")
		}

		res, ok := wayback.MockGetUrls(f.Website, f.Mimetype, f.Limit, f.ResumeKey)

		if !ok {
			end := time.Now()
			fmt.Printf("end: %v\n", end)
			fmt.Printf("duration: %v\n", end.Sub(start))

			return c.SendString("Coś nie poszło")
		}

		end := time.Now()
		fmt.Printf("end: %v\n", end)

		f.ResumeKey = res.ResumeKey

		return c.Render("results", fiber.Map{
			"Urls":     res.Urls,
			"FormNext": f,
			"Duration": end.Sub(start).String(),
		})
	})
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/public", "./public")

	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
