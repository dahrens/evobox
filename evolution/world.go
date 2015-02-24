package evolution

import "log"

type WorldChanger interface {
	Obj() Evolver
	X() int
	Y() int
}

type World struct {
	W, H      int
	Creatures Creatures
	_m        []Creatures
	Requests  chan WorldChanger
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
	world.Requests = make(chan WorldChanger)
	return world
}

func (world *World) Run() {
	for {
		request := <-world.Requests
		world.handle(request) // maybe as goroutine?
	}
}

func (world *World) handle(req WorldChanger) {
	if req.X() <= -1 || req.X() >= world.W {
		log.Printf("Out of bounds x %d reject put request!", req.X())
		return
	}
	if req.Y() <= -1 || req.Y() >= world.H {
		log.Printf("Out of bounds y %d reject put request!", req.Y())
		return
	}
	switch r := req.(type) {
	case *PutRequest:
		world.handlePut(r)
	case *DeleteRequest:
		world.handleDelete(r)
	}
}

func (world *World) handlePut(req *PutRequest) {
	switch o := req.Obj().(type) {
	case *Creature:
		if world.Creatures.Contains(o) {
			log.Printf("x %d, y %d", req.X(), req.Y())
			if world._m[req.X()][req.Y()] != nil {
				log.Println("Can not place here...")
				return
			}
			o.SetX(req.X())
			o.SetY(req.Y())
			world._m[o.X()][o.Y()] = o
			world._m[o.X()][o.Y()] = nil
		} else {
			world.Creatures = append(world.Creatures, o)
			world._m[o.X()][o.Y()] = o
			go o.Evolve(world.Requests)
		}
	}
}

func (world *World) handleDelete(req *DeleteRequest) {
	log.Println("delete me! handler")
	switch o := req.Obj().(type) {
	case *Creature:
		world._m[o.X()][o.Y()] = nil
		world.RemoveCreature(o)
	}
}

func (world *World) RemoveCreature(creature *Creature) {
	var i int
	var t *Creature
	for i, t = range world.Creatures {
		if t == creature {
			break
		}
	}
	world.Creatures[i] = nil
	world.Creatures = append(world.Creatures[:i], world.Creatures[i+1:]...)
}

func (world *World) PutCreature(c *Creature) {
	req := new(PutRequest)
	req.obj = c
	req.x = c.X()
	req.y = c.Y()
	world.Requests <- req
}

type Request struct {
	obj Evolver
}

type PutRequest struct {
	Request
	x, y int
}

func (p PutRequest) Obj() Evolver { return p.obj }
func (p PutRequest) X() int       { return p.x }
func (p PutRequest) Y() int       { return p.y }

type DeleteRequest struct {
	Request
}

func (p DeleteRequest) Obj() Evolver { return p.obj }
func (p DeleteRequest) X() int       { return p.obj.X() }
func (p DeleteRequest) Y() int       { return p.obj.Y() }
