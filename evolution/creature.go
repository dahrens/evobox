package evolution

import (
	"github.com/Pallinder/go-randomdata"
	"log"
)

type Gender string

const (
	GENDER_MALE     Gender = "M"
	GENDER_FEMALE   Gender = "F"
	HUNGER_PER_TICK        = 0.1
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
	target     Evolver
	client     *Client `json:"-"`
	world      *World  `json:"-"`
}

func NewCreature(health float32, gender Gender, client *Client) *Creature {
	creature := new(Creature)
	// Fragment values
	creature.X = client.Rand.Intn(client.World.W)
	creature.Y = client.Rand.Intn(client.World.H)
	creature.W = 32
	creature.H = 32
	creature.Birth = client.World.Tick
	creature.Age = 0
	creature.pulse = make(chan *Tick)
	// Creature values
	creature.Name = randomdata.SillyName()
	creature.Health = health
	creature.Hunger = 5.0
	creature.Libido = 0.0
	creature.Sanity = 100.0
	creature.Speed = 10
	creature.Gender = gender
	creature.hunger_max = 10.0
	creature.libido_max = 10.0
	creature.alive = false
	creature.Id = client.AutoInc.Id()
	creature.client = client
	creature.target = nil
	return creature
}

func (creature *Creature) Evolve(world *World) {
	creature.alive = true
	creature.world = world
	creature.world.Client.Write(NewMessage("add-creature", creature))
	for tick := range creature.pulse {
		creature.Age = tick.Count - creature.Birth
		creature.calculateLibido()
		creature.calculateHunger()
		creature.calculateHealth()
		creature.act()
		tick.Wait.Done()
		creature.client.Write(NewMessage("update-creature", creature))
		if creature.Health == 0 {
			break
		}
	}
	creature.alive = false
	world.Requests <- &DeleteRequest{Request{obj: creature}}
	creature.world = nil
}

func (creature *Creature) act() {
	log.Println("act...")
	picked := creature.pickFlower()
	if picked {
		creature.target = nil
		return
	}
	if creature.target == nil {
		flower := creature.searchFlower()
		if flower != nil {
			creature.target = flower
		} else {
			creature.moveRandom()
		}
	}
	if creature.target != nil {
		creature.moveTo(creature.target)
	}

	log.Println("act...done")
}

func (creature *Creature) pickFlower() bool {
	picked := false
	reach := creature.world.Plan.getFragmentsInCircle(creature.X, creature.Y, 50)
PickFlower:
	for i := 0; i < len(reach); i++ {
		switch t := reach[i].(type) {
		case *Flower:
			creature.world.Requests <- &DeleteRequest{Request{obj: t}}
			creature.Hunger = creature.Hunger - t.NutritionalValue
			picked = true
			break PickFlower
		}
	}
	return picked
}

func (creature *Creature) searchFlower() *Flower {
	var flower *Flower
	vision := creature.world.Plan.getFragmentsInCircle(creature.X, creature.Y, 200)
SearchForFlower:
	for j := 0; j < len(vision); j++ {
		switch f := vision[j].(type) {
		case *Flower:
			flower = f
			break SearchForFlower
		}
	}
	return flower
}

func (creature *Creature) moveRandom() {
	coin := creature.client.Rand.Intn(4)
	newX := creature.X
	newY := creature.Y
	switch coin {
	case 0:
		newX -= creature.Speed
	case 1:
		newX += creature.Speed
	case 2:
		newY -= creature.Speed
	case 3:
		newY += creature.Speed
	default:
		return
	}
	req := new(PostRequest)
	req.obj = creature
	req.x = newX
	req.y = newY
	creature.client.World.Requests <- req
}

func (creature *Creature) moveTo(fragment Fragmenter) {
	log.Println("move from %d %d to flower at (%d, %d)", creature.X, creature.Y, fragment.GetX(), fragment.GetY())
	newX := creature.X
	newY := creature.Y

	if fragment.GetX() > creature.X+creature.Speed {
		newX = creature.X + creature.Speed
	}
	if fragment.GetY() > creature.Y+creature.Speed {
		newY = creature.Y + creature.Speed
	}
	if fragment.GetX() < creature.X-creature.Speed {
		newX = creature.X - creature.Speed
	}
	if fragment.GetY() < creature.Y-creature.Speed {
		newY = creature.Y - creature.Speed
	}
	log.Println("newx newy %d %d", newX, newY)
	req := new(PostRequest)
	req.obj = creature
	req.x = newX
	req.y = newY
	creature.world.Requests <- req
}

func (creature *Creature) calculateHunger() {
	if creature.Hunger < creature.hunger_max {
		creature.Hunger = creature.Hunger + HUNGER_PER_TICK
	}
	if creature.Hunger > creature.hunger_max {
		creature.Hunger = creature.hunger_max
	}
}

func (creature *Creature) calculateLibido() {
	if creature.Libido < creature.libido_max {
		creature.Libido = creature.Libido + LIBIDO_PER_TICK
	}
	if creature.Libido > creature.libido_max {
		creature.Libido = creature.libido_max
	}
}

func (creature *Creature) calculateHealth() {
	if creature.Hunger == creature.hunger_max {
		creature.Health--
	}
}

func (creature *Creature) GetX() int {
	return creature.X
}

func (creature *Creature) SetX(x int) {
	creature.X = x
}

func (creature *Creature) GetY() int {
	return creature.Y
}

func (creature *Creature) SetY(y int) {
	creature.Y = y
}

func (creature *Creature) GetW() int {
	return creature.W
}

func (creature *Creature) SetW(w int) {
	creature.W = w
}

func (creature *Creature) GetH() int {
	return creature.H
}

func (creature *Creature) SetH(h int) {
	creature.H = h
}

func (creature *Creature) Pulse() chan *Tick {
	return creature.pulse
}

func (creature *Creature) Alive() bool {
	return creature.alive
}

func (creature *Creature) Collides() bool { return true }

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
