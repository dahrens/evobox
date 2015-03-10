package evolution

// import (
// 	"github.com/Pallinder/go-randomdata"
// )

type Flower struct {
	EvolverFragment
	Placeable
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
	for tick := range flower.pulse {
		flower.Age = tick.Count - flower.Birth
		tick.Wait.Done()
	}
}

type Flowers []*Flower
