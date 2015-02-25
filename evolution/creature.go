package evolution

import (
	"github.com/Pallinder/go-randomdata"
	"math/rand"
	"sort"
	"strings"
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
	Name          string
	Health        float32
	Hunger        float32
	Libido        float32
	Sanity        float32
	Gender        Gender
	hunger_max    float32
	libido_max    float32
	alive         bool
	rand          *rand.Rand
	worldRequests chan Requester
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

func (self *Creature) Evolve(worldRequests chan Requester) {
	self.alive = true
	self.worldRequests = worldRequests
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
	worldRequests <- &DeleteRequest{Request{obj: self}}
	self.worldRequests = nil
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

func (self *Creature) SetX(x int) {
	self.x = x
}

func (self *Creature) Y() int {
	return self.y
}

func (self *Creature) SetY(y int) {
	self.y = y
}

func (self *Creature) Pulse() chan int {
	return self.pulse
}

func (self *Creature) Alive() bool {
	return self.alive
}

func (self *Creature) move() {
	coin := self.rand.Intn(4)
	newX := self.x
	newY := self.y
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
	self.worldRequests <- req
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
			record["X"] = e.X()
			record["Y"] = e.Y()
			data[p] = record
			p++
		}
	}
	return data
}

func (c Creatures) Sort(column, direction string) {
	s := strings.Join(append(make([]string, 2), column, direction), "")
	switch s {
	case "Healthasc":
		sort.Sort(ByHealthAsc(c))
	case "Healthdesc":
		sort.Sort(ByHealthDesc(c))
	case "Nameasc":
		sort.Sort(ByNameAsc(c))
	case "Namedesc":
		sort.Sort(ByNameDesc(c))
	case "Genderasc":
		sort.Sort(ByGenderAsc(c))
	case "Genderdesc":
		sort.Sort(ByGenderDesc(c))
	case "Ageasc":
		sort.Sort(ByAgeAsc(c))
	case "Agedesc":
		sort.Sort(ByAgeDesc(c))
	case "Hungerasc":
		sort.Sort(ByHungerAsc(c))
	case "Hungerdesc":
		sort.Sort(ByHungerDesc(c))
	case "Libidoasc":
		sort.Sort(ByLibidoAsc(c))
	case "Libidodesc":
		sort.Sort(ByLibidoDesc(c))
	}
}

type ByNameAsc Creatures

func (a ByNameAsc) Len() int           { return len(a) }
func (a ByNameAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNameAsc) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ByNameDesc Creatures

func (a ByNameDesc) Len() int           { return len(a) }
func (a ByNameDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNameDesc) Less(i, j int) bool { return a[i].Name > a[j].Name }

type ByGenderAsc Creatures

func (a ByGenderAsc) Len() int           { return len(a) }
func (a ByGenderAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenderAsc) Less(i, j int) bool { return a[i].Gender < a[j].Gender }

type ByGenderDesc Creatures

func (a ByGenderDesc) Len() int           { return len(a) }
func (a ByGenderDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenderDesc) Less(i, j int) bool { return a[i].Gender > a[j].Gender }

type ByAgeAsc Creatures

func (a ByAgeAsc) Len() int           { return len(a) }
func (a ByAgeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAgeAsc) Less(i, j int) bool { return a[i].Age < a[j].Age }

type ByAgeDesc Creatures

func (a ByAgeDesc) Len() int           { return len(a) }
func (a ByAgeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAgeDesc) Less(i, j int) bool { return a[i].Age > a[j].Age }

type ByHungerAsc Creatures

func (a ByHungerAsc) Len() int           { return len(a) }
func (a ByHungerAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHungerAsc) Less(i, j int) bool { return a[i].Hunger < a[j].Hunger }

type ByHungerDesc Creatures

func (a ByHungerDesc) Len() int           { return len(a) }
func (a ByHungerDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHungerDesc) Less(i, j int) bool { return a[i].Hunger > a[j].Hunger }

type ByLibidoAsc Creatures

func (a ByLibidoAsc) Len() int           { return len(a) }
func (a ByLibidoAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLibidoAsc) Less(i, j int) bool { return a[i].Libido < a[j].Libido }

type ByLibidoDesc Creatures

func (a ByLibidoDesc) Len() int           { return len(a) }
func (a ByLibidoDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLibidoDesc) Less(i, j int) bool { return a[i].Libido > a[j].Libido }

type ByHealthAsc Creatures

func (a ByHealthAsc) Len() int           { return len(a) }
func (a ByHealthAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHealthAsc) Less(i, j int) bool { return a[i].Health < a[j].Health }

type ByHealthDesc Creatures

func (a ByHealthDesc) Len() int           { return len(a) }
func (a ByHealthDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHealthDesc) Less(i, j int) bool { return a[i].Health > a[j].Health }
