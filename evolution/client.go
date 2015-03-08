package evolution

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
	"math/rand"
	"time"
)

type Client struct {
	Conn    *websocket.Conn
	World   *World
	AutoInc *AutoInc
	Rand    *rand.Rand
	Id      int
	server  *Server
	ch      chan *Message
	doneCh  chan bool
}

var channelBufSize = 1000

// Create new chat client.
func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("ws cannot be nil")
	}

	if server == nil {
		panic("server cannot be nil")
	}
	c := new(Client)
	c.ch = make(chan *Message, channelBufSize)
	c.doneCh = make(chan bool)
	c.AutoInc = NewAutoInc(0, 1)
	c.Id = server.Next.Id()
	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.World = NewWorld(64, 64, c)
	c.server = server
	c.Conn = ws
	return c
}

func (c *Client) WebsocketHandler(ws *websocket.Conn) {
	c.Conn = ws
	io.Copy(ws, ws)
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	go c.listenWrite()
	c.listenRead()
}

// Listen write request via chanel
func (c *Client) listenWrite() {
	for {
		select {

		// send message to the client
		case msg := <-c.ch:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.Conn, msg)

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenRead method
			return
		}
	}
}

// Listen read request via chanel
func (c *Client) listenRead() {
	for {
		select {

		// receive done request
		case <-c.doneCh:
			c.server.Del(c)
			c.doneCh <- true // for listenWrite method
			return

		// read data from websocket connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.Conn, &msg)
			if err == io.EOF {
				c.doneCh <- true
			} else if err != nil {
				c.server.Err(err)
			} else {
				DispatchIncomingMessage(&msg, c)
			}
		}
	}
}

func (c *Client) Write(msg *Message) {
	select {
	case c.ch <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.Id)
		c.server.Err(err)
	}
}
