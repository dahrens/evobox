package evolution

import (
	"github.com/Pallinder/go-randomdata"
	"log"
)

type Gender string

const (
	GENDER_MALE     Gender = "M"
	GENDER_FEMALE   Gender = "F"
	HUNGER_PER_TICK        = 0.5
	LIBIDO_PER_TICK        = 0.1
)

type Creature struct {
	EvolverFragment
	Name       string
	Health     float32
	Hunger     float32
	Libido     float32
	Sanity     float32
	Speed      int // pixels per tick
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
	c.W = 32
	c.H = 32
	c.Birth = client.World.Tick
	c.Age = 0
	c.pulse = make(chan *Tick)
	// Creature values
	c.Name = randomdata.SillyName()
	c.Health = health
	c.Hunger = 5.0
	c.Libido = 0.0
	c.Sanity = 100.0
	c.Speed = 32
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
		tick.Wait.Done()
		if self.Health == 0 {
			break
		}
	}
	self.alive = false
	world.Requests <- &DeleteRequest{Request{obj: self}}
	self.client.Write(NewMessage("delete-creature", self))
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

func (self *Creature) GetW() int {
	return self.W
}

func (self *Creature) SetW(w int) {
	self.W = w
}

func (self *Creature) GetH() int {
	return self.H
}

func (self *Creature) SetH(h int) {
	self.H = h
}

func (self *Creature) Pulse() chan *Tick {
	return self.pulse
}

func (self *Creature) Alive() bool {
	return self.alive
}

func (self *Creature) move() {
	coin := self.client.Rand.Intn(12)
	newX := self.X
	newY := self.Y
	switch coin {
	case 0:
		newX -= 64
	case 1:
		newX += 64
	case 2:
		newY -= 64
	case 3:
		newY += 64
	default:
		log.Println("idle")
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
