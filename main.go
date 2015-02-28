package main

import (
	"code.google.com/p/go.net/websocket"
	evo "github.com/dahrens/evobox/evolution"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

var (
	server *evo.Server = evo.NewServer("/connect")
)

func connect(ws *websocket.Conn) {
	client := evo.NewClient(ws, server)
	server.Add(client)
	client.Listen()
}

func main() {
	r := gin.Default()
	go server.Listen()

	r.GET("/connect", func(c *gin.Context) {
		handler := websocket.Handler(connect)
		handler.ServeHTTP(c.Writer, c.Request)
	})

	r.Use(static.Serve("/", static.LocalFile("assets", false)))

	// Listen and server on 0.0.0.0:8080
	r.Run(":8080")
}
