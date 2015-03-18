package evolution

import (
	"log"
	"math/rand"
)

type Flower struct {
	EvolverFragment
	Placeable
	Cycle            int
	ThrowingRangeMin int
	ThrowingRangeMax int
	NutritionalValue float32
	world            *World `json:"-"`
	alive            bool
	Id               int
}

func NewFlower(x, y int, sprite string, world *World) *Flower {
	flower := new(Flower)
	flower.X = x
	flower.Y = y
	flower.W = 39
	flower.H = 69
	flower.Anchor.X = 0.5
	flower.Anchor.Y = 0.5
	flower.Sprite = sprite
	flower.Sheet = "default"
	flower.Birth = world.Tick
	flower.world = world
	flower.Id = world.Client.AutoInc.Id()
	flower.alive = false
	flower.Age = 0
	flower.Cycle = random(20, 30)
	flower.ThrowingRangeMin = 20
	flower.ThrowingRangeMax = 50
	flower.pulse = make(chan *Tick)
	return flower
}

func (flower *Flower) GetX() int { return flower.X }
func (flower *Flower) GetY() int { return flower.Y }
func (flower *Flower) GetW() int { return flower.W }
func (flower *Flower) GetH() int { return flower.H }

func (flower *Flower) SetX(x int) { flower.X = x }
func (flower *Flower) SetY(y int) { flower.Y = y }
func (flower *Flower) SetW(w int) { flower.W = w }
func (flower *Flower) SetH(h int) { flower.H = h }

func (flower *Flower) Alive() bool {
	return flower.alive
}

func (flower *Flower) Pulse() chan *Tick {
	return flower.pulse
}

func (flower *Flower) Evolve(world *World) {
	flower.alive = true
	flower.world = world
	flower.world.Client.Write(NewMessage("add-flower", flower))
	for tick := range flower.pulse {
		flower.Age = tick.Count - flower.Birth
		flower.NutritionalValue += 0.2
		if flower.Age%flower.Cycle == 0 {
			x := random(flower.ThrowingRangeMin, flower.ThrowingRangeMax)
			switch rand.Intn(2) {
			case 1:
				x = x * -1
			}
			y := random(flower.ThrowingRangeMin, flower.ThrowingRangeMax)
			switch rand.Intn(2) {
			case 1:
				y = y * -1
			}
			f := NewFlower(flower.X+x, flower.Y+y, flower.Sprite, flower.world)
			flower.world.Requests <- NewPutRequest(f)
		}
		tick.Wait.Done()
		flower.world.Client.Write(NewMessage("update-flower", flower))
		// in 90 percent of the cases the flower dies after a cycle
		if flower.Age%flower.Cycle == 0 && rand.Intn(100) < flower.Age {
			break
		}
	}
	flower.alive = false
	world.Requests <- &DeleteRequest{Request{obj: flower}}
	flower.world.Client.Write(NewMessage("delete-flower", flower))
	flower.world = nil
}

type Flowers []*Flower
