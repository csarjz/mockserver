package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/csarjz/mockserver/cmd/entity"
	"github.com/gin-gonic/gin"
)

var server *http.Server

func main() {

	var serverConfigFileName = "server.json"
	var serverConfig, err = decodeServerConfigFile(serverConfigFileName)
	if err != nil {
		log.Fatal(err)
	}

	startServer(serverConfig)

	select {}
}

func decodeServerConfigFile(fileName string) (*entity.ServerConfig, error) {
	var serverConfigFile, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer serverConfigFile.Close()

	var serverConfig = entity.ServerConfig{}
	decoder := json.NewDecoder(serverConfigFile)
	err = decoder.Decode(&serverConfig)
	if err != nil {
		return nil, errors.New("malformed " + fileName)
	}

	return &serverConfig, nil
}

func startServer(serverConfig *entity.ServerConfig) {
	gin.SetMode(gin.ReleaseMode)
	var router = gin.Default()

	initializeServerRoutes(serverConfig, router)

	server = &http.Server{
		Addr:    ":" + strconv.Itoa(int(serverConfig.Port)),
		Handler: router,
	}

	go func() {
		var spaceAndGreen, blue, resetAndSpace = "\n\n\n\033[32m", "\033[34m", "\033[0m\n\n\n"

		log.Printf("%sStarting server on %shttp://localhost:%d%s", spaceAndGreen, blue, serverConfig.Port, resetAndSpace)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()
}

func initializeServerRoutes(serverConfig *entity.ServerConfig, router *gin.Engine) {

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to Mock Server By github.com/csarjz")
	})

	var api = router.Group(serverConfig.BaseUrl)

	for i := 0; i < len(serverConfig.Routes); i++ {
		route := serverConfig.Routes[i]
		if route.HttpStatus <= http.StatusContinue {
			route.HttpStatus = http.StatusOK
		}
		switch route.Method {
		case "POST":
			api.POST(route.Path, func(c *gin.Context) {
				delay(route.Delay)
				processResponse(route.HttpStatus, route.ResponseFile, c)
			})
		case "PUT":
			api.PUT(route.Path, func(c *gin.Context) {
				delay(route.Delay)
				processResponse(route.HttpStatus, route.ResponseFile, c)
			})
		case "DELETE":
			api.DELETE(route.Path, func(c *gin.Context) {
				delay(route.Delay)
				processResponse(route.HttpStatus, route.ResponseFile, c)
			})
		default:
			api.GET(route.Path, func(c *gin.Context) {
				delay(route.Delay)
				processResponse(route.HttpStatus, route.ResponseFile, c)
			})
		}
	}
}

func delay(delay uint32) {
	if delay > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}

func processResponse(httpStatus uint32, responseFilePath string, c *gin.Context) {
	var jsonFile, err = os.Open(responseFilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Message: "JSON File Not Found"})
	}
	defer jsonFile.Close()

	c.Header("Content-Type", "application/json")
	c.Status(int(httpStatus))
	_, err = io.Copy(c.Writer, jsonFile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{Message: "JSON File Not Found"})
	}
}
