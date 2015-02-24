package evolution

import "log"

type Changer interface {
	Obj() Evolver
	X() int
	Y() int
}

type World struct {
	W, H      int
	Creatures Creatures
	_m        []Creatures
	requests  chan Changer
}

func NewWorld(w, h int) *World {
	world := new(World)
	world.W = w
	world.H = h
	world.Creatures = make(Creatures, 0)
	world._m = make([]Creatures, w)
	for x := 0; x < w; x++ {
		world._m[x] = make(Creatures, h)
	}
	world.requests = make(chan Changer)
	return world
}

func (world *World) Run() {
	for {
		request := <-world.requests
		go world.handle(request) // Don't wait for handle to finish.
	}
}

func (world *World) handle(req Changer) {
	switch o := req.Obj().(type) {
	case *Creature:
		if world.Creatures.Contains(o) {
			log.Println("Not supported atm...")
		} else {
			world.Creatures = append(world.Creatures, o)
			world._m[o.X()][o.Y()] = o
		}
	}

}

func (world *World) PutCreature(c *Creature) {
	req := new(Put)
	req.obj = c
	req.x = c.X()
	req.y = c.Y()
	world.requests <- req
}

type Request struct {
	obj Evolver
}

type Put struct {
	Request
	x, y int
}

func (p Put) Obj() Evolver { return p.obj }
func (p Put) X() int       { return p.x }
func (p Put) Y() int       { return p.y }
