package evolution

import "github.com/Pallinder/go-randomdata"
import "math/rand"

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
	rand       *rand.Rand
}

func NewCreature(x, y int, health float32, gender Gender, r *rand.Rand) *Creature {
	c := new(Creature)
	// Fragment values
	c.x = x
	c.y = y
	c.Age = 0
	c.pulse = make(chan int)
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
	c.rand = r
	return c
}

func (self *Creature) Evolve() {
	self.alive = true
	for tick := range self.pulse {
		self.Age = tick
		self.calculateLibido()
		self.calculateHunger()
		self.calculateHealth()
		self.move()
		if self.Health == 0 {
			break
		}
	}
	self.alive = false
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

func (self *Creature) X() int {
	return self.x
}

func (self *Creature) Y() int {
	return self.y
}

func (self *Creature) Pulse() chan int {
	return self.pulse
}

func (self *Creature) Alive() bool {
	return self.alive
}

func (self *Creature) move() {
	coin := self.rand.Intn(4)
	switch coin {
	case 0:
		self.x--
	case 1:
		self.x++
	case 2:
		self.y--
	case 3:
		self.y++
	}
}
