package go_server_core

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type WebServer struct {
	ListenPort string
	SendAddr   string
	App        *fiber.App
}

func NewWebServer(listenPort string, sendAddr string) *WebServer {
	wb := &WebServer{
		ListenPort: listenPort,
		SendAddr:   sendAddr,
	}
	wb.App = fiber.New(fiber.Config{})
	wb.App.Use(logger.New())

	go wb.App.Listen(listenPort)
	return wb
}

func (w *WebServer) SendGetPacket(api string, json []byte) {
	go func() {
		url := "http://" + w.SendAddr + "/" + api
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(json))
		if err != nil {
			log.Println(err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("WEBSEND RESP :", string(body))
		defer resp.Body.Close()
	}()
}
