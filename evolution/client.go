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

var maxId int = 0
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
	maxId++
	c.ch = make(chan *Message, channelBufSize)
	c.doneCh = make(chan bool)
	c.AutoInc = NewAutoInc(0, 1)
	c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	c.World = NewWorld(32, 32)
	c.server = server
	c.Conn = ws
	return c
}

func (client *Client) WebsocketHandler(ws *websocket.Conn) {
	client.Conn = ws
	io.Copy(ws, ws)
}

// Listen Write and Read request via chanel
func (c *Client) Listen() {
	c.Init()
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
				if action, ok := msg["action"]; ok {
					switch action {
					case "connect":
						websocket.JSON.Send(c.Conn, NewMessage("load-world", c.World))
					case "Start":
						c.World.Run()
					case "Pause":
						c.World.Pause()
					case "Reset":
						c.World.Reset()
						c.Init()
						websocket.JSON.Send(c.Conn, NewMessage("load-world", c.World))
					}
				}
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

func (client *Client) Init() {
	client.SpawnMany(10, GENDER_MALE)
	client.SpawnMany(10, GENDER_FEMALE)
}

func (client *Client) SpawnMany(n int, gender Gender) {
	for i := 0; i < n; i++ {
		e := NewCreature(float32(i+10), gender, client)
		client.Spawn(e)
	}
}

func (client *Client) Spawn(e Evolver) {
	switch obj := e.(type) {
	case *Creature:
		r := NewPutRequest(obj)
		client.World.Requests <- r
	}
}
