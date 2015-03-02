package evolution

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
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
	c.World = NewWorld(64, 64)
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
						settings := msg["settings"].(map[string]interface{})
						count, _ := strconv.Atoi(settings["initial_creatures"].(string))
						tick_interval, _ := strconv.Atoi(settings["tick_interval"].(string))
						map_width, _ := strconv.Atoi(settings["map_width"].(string))
						map_height, _ := strconv.Atoi(settings["map_height"].(string))
						c.World.Reset(tick_interval, map_width, map_height)
						c.Init(count)
						websocket.JSON.Send(c.Conn, NewMessage("load-world", c.World))
					case "Start":
						c.World.Run()
					case "Pause":
						c.World.Pause()
					case "Reset":
						settings := msg["settings"].(map[string]interface{})
						count, _ := strconv.Atoi(settings["initial_creatures"].(string))
						tick_interval, _ := strconv.Atoi(settings["tick_interval"].(string))
						map_width, _ := strconv.Atoi(settings["map_width"].(string))
						map_height, _ := strconv.Atoi(settings["map_height"].(string))
						c.World.Reset(tick_interval, map_width, map_height)
						c.Init(count)
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

func (client *Client) Init(initialCreatures int) {
	client.SpawnMany(initialCreatures/2, GENDER_MALE)
	client.SpawnMany(initialCreatures/2, GENDER_FEMALE)
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
