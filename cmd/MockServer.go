package main

import (
	"encoding/json"
	"github.com/csarjz/mockserver/cmd/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

func main() {
	var app = fiber.New()

	serverConfigFile, err := ioutil.ReadFile("server.json")
	if err != nil {
		log.Fatal("server.json file not found")
	}

	var serverConfig = entity.ServerConfig{}
	err = json.Unmarshal(serverConfigFile, &serverConfig)
	if err != nil {
		log.Fatal("malformed app.json")
	}

	app.Use(logger.New())
	app.Use(recover.New())
	initializeServerRoutes(serverConfig, app)

	log.Fatal(app.Listen(":" + strconv.Itoa(int(serverConfig.Port))))
}

func initializeServerRoutes(serverConfig entity.ServerConfig, app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Mock Server By github.com/csarjz")
	})

	api := app.Group(serverConfig.BaseUrl)

	for i := 0; i < len(serverConfig.Routes); i++ {
		route := serverConfig.Routes[i]
		switch route.Method {
		case "POST":
			api.Post(route.Path, func(c *fiber.Ctx) error {
				delay(route.Delay)
				return processResponse(route.ResponseFile, c)
			})
		case "PUT":
			api.Put(route.Path, func(c *fiber.Ctx) error {
				delay(route.Delay)
				return processResponse(route.ResponseFile, c)
			})
		case "DELETE":
			api.Delete(route.Path, func(c *fiber.Ctx) error {
				delay(route.Delay)
				return processResponse(route.ResponseFile, c)
			})
		default:
			api.Get(route.Path, func(c *fiber.Ctx) error {
				delay(route.Delay)
				return processResponse(route.ResponseFile, c)
			})
		}
	}
}

func processResponse(filePath string, c *fiber.Ctx) error {
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		c.Response().Header.SetStatusCode(fiber.StatusBadRequest)
		return c.JSON(entity.ErrorResponse{
			Message: "JSON File Not Found",
		})
	}
	c.Response().Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
	return c.Send(jsonFile)
}

func delay(delay uint32) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
