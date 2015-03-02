package evolution

import (
	"github.com/Pallinder/go-randomdata"
)

type Gender string

const (
	GENDER_MALE     Gender = "M"
	GENDER_FEMALE   Gender = "F"
	HUNGER_PER_TICK        = 0.5
	LIBIDO_PER_TICK        = 0.1
)

type Creature struct {
	Fragment
	Name       string
	Health     float32
	Hunger     float32
	Libido     float32
	Sanity     float32
	Gender     Gender
	hunger_max float32
	libido_max float32
	alive      bool
	Id         int
	client     *Client `json:"-"`
	world      *World  `json:"-"`
}

func NewCreature(health float32, gender Gender, client *Client) *Creature {
	c := new(Creature)
	// Fragment values
	c.X = client.Rand.Intn(client.World.W)
	c.Y = client.Rand.Intn(client.World.H)
	c.Birth = client.World.Tick
	c.Age = 0
	c.pulse = make(chan *Tick)
	// Creature values
	c.Name = randomdata.SillyName()
	c.Health = health
	c.Hunger = 5.0
	c.Libido = 0.0
	c.Sanity = 100.0
	c.Gender = gender
	c.hunger_max = 10.0
	c.libido_max = 10.0
	c.alive = false
	c.Id = client.AutoInc.Id()
	c.client = client
	return c
}

func (self *Creature) Evolve(world *World) {
	self.alive = true
	self.world = world
	for tick := range self.pulse {
		self.Age = tick.Count - self.Birth
		self.calculateLibido()
		self.calculateHunger()
		self.calculateHealth()
		self.move()
		self.client.Write(NewMessage("update", self))
		tick.Wait.Done()
		if self.Health == 0 {
			break
		}
	}
	self.alive = false
	world.Requests <- &DeleteRequest{Request{obj: self}}
	self.client.Write(NewMessage("delete", self))
	self.world = nil
}

func (self *Creature) calculateHunger() {
	if self.Hunger < self.hunger_max {
		self.Hunger = self.Hunger + HUNGER_PER_TICK
	}
	if self.Hunger > self.hunger_max {
		self.Hunger = self.hunger_max
	}
}

func (self *Creature) calculateLibido() {
	if self.Libido < self.libido_max {
		self.Libido = self.Libido + LIBIDO_PER_TICK
	}
	if self.Libido > self.libido_max {
		self.Libido = self.libido_max
	}
}

func (self *Creature) calculateHealth() {
	if self.Hunger == self.hunger_max {
		self.Health--
	}
}

func (self *Creature) GetX() int {
	return self.X
}

func (self *Creature) SetX(x int) {
	self.X = x
}

func (self *Creature) GetY() int {
	return self.Y
}

func (self *Creature) SetY(y int) {
	self.Y = y
}

func (self *Creature) Pulse() chan *Tick {
	return self.pulse
}

func (self *Creature) Alive() bool {
	return self.alive
}

func (self *Creature) move() {
	coin := self.client.Rand.Intn(4)
	newX := self.X
	newY := self.Y
	switch coin {
	case 0:
		newX--
	case 1:
		newX++
	case 2:
		newY--
	case 3:
		newY++
	}
	req := new(PostRequest)
	req.obj = self
	req.x = newX
	req.y = newY
	self.world.Requests <- req
}

type Creatures []*Creature

func (c Creatures) Contains(creature *Creature) bool {
	for _, b := range c {
		if b == creature {
			return true
		}
	}
	return false
}

func (c Creatures) Remove(creature *Creature) {
	var i int
	var t *Creature
	for i, t = range c {
		if t == creature {
			break
		}
	}
	nc := make(Creatures, len(c)-1)
	c[i] = nil
	nc = append(c[:i], c[i+1:]...)
	c = nc
}

func (c Creatures) ToMap(start, end int) []interface{} {
	data := make([]interface{}, end-start)
	if start < len(c) {
		p := 0
		for i := start; i < end; i++ {
			e := c[i]
			record := make(map[string]interface{})
			record["Name"] = e.Name
			record["Age"] = e.Age
			record["Health"] = e.Health
			record["Libido"] = e.Libido
			record["Hunger"] = e.Hunger
			record["Gender"] = e.Gender
			record["X"] = e.X
			record["Y"] = e.Y
			record["Id"] = e.Id
			data[p] = record
			p++
		}
	}
	return data
}
