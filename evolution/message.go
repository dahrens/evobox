package evolution

import "strconv"
import "log"

type Message struct {
	Action string
	Data   interface{}
}

func NewMessage(action string, data interface{}) *Message {
	m := new(Message)
	m.Action = action
	m.Data = data
	return m
}

func DispatchIncomingMessage(msg *Message, c *Client) {
	switch msg.Action {
	case "Reset":
		fallthrough
	case "Connect":
		handleConnect(msg.Data, c)
	case "Start":
		c.World.Run()
	case "Pause":
		c.World.Pause()
	case "Spawn":
		handleSpawn(msg.Data, c)
	}
}

func handleConnect(data interface{}, c *Client) {
	initial_creatures, _ := strconv.Atoi(data.(map[string]interface{})["initial_creatures"].(string))
	initial_flowers, _ := strconv.Atoi(data.(map[string]interface{})["initial_flowers"].(string))
	log.Println(initial_flowers)
	tick_interval, _ := strconv.Atoi(data.(map[string]interface{})["tick_interval"].(string))
	c.World.Reset(tick_interval, 2048, 2048)
	c.World.Init(initial_creatures, initial_flowers)
	c.Write(NewMessage("load-world", c.World))
}

func handleSpawn(data interface{}, c *Client) {
	coin := c.Rand.Intn(2)
	var gender Gender
	switch coin {
	case 0:
		gender = GENDER_FEMALE
	default:
		gender = GENDER_MALE
	}
	e := NewCreature(100.0, gender, c)
	c.World.Spawn(e)
	c.Write(NewMessage("add-creature", e))
}
