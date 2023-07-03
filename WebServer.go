package go_server_core

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type WebServer struct {
	Port string
	App  *fiber.App
}

func NewWebServer(port string) *WebServer {
	wb := &WebServer{
		Port: port,
	}
	wb.App = fiber.New(fiber.Config{})
	wb.App.Use(logger.New())

	go wb.App.Listen(port)

	resp, err := http.Get("http://127.0.0.1:3000/matchcomplete")
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	return wb
}
