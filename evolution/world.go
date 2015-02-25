package evolution

import (
	"log"
	"math/rand"
	"time"
)

type Requester interface {
	Obj() Evolver
	X() int
	Y() int
}

type World struct {
	W, H      int
	Evolvers  Evolvers
	Creatures Creatures
	Map       []Creatures
	Requests  chan Requester
	Clock     *time.Ticker
	Tick      int
	Speed     time.Duration
	Rand      *rand.Rand
}

func NewWorld(w, h int) *World {
	world := new(World)
	world.W = w
	world.H = h
	world.Evolvers = make(Evolvers, 0)
	world.Creatures = make(Creatures, 0)
	world.Map = make([]Creatures, w)
	for x := 0; x < w; x++ {
		world.Map[x] = make(Creatures, h)
	}
	world.Requests = make(chan Requester)
	world.Clock = time.NewTicker(time.Second)
	world.Tick = 0
	world.Speed = time.Second
	world.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	go world.serve()
	return world
}

func (world *World) Init() {
	world.SpawnMany(15, GENDER_MALE)
	world.SpawnMany(15, GENDER_FEMALE)
}

func (world *World) Run() {
	go world.pulse()
}

func (world *World) Pause() {
	world.Clock.Stop()
	world.Clock = time.NewTicker(world.Speed)
}

func (world *World) Reset() {

}

func (world *World) serve() {
	for {
		request := <-world.Requests
		world.handle(request) // maybe as goroutine?
	}
}

func (world *World) pulse() {
	for range world.Clock.C {
		world.Tick++
		t := make(Evolvers, len(world.Evolvers), cap(world.Evolvers))
		copy(t, world.Evolvers)
		for _, e := range t {
			if e.Alive() {
				e.Pulse() <- world.Tick
			}
		}
	}
}

func (world *World) SpawnMany(n int, gender Gender) {
	for i := 0; i < n; i++ {
		e := NewCreature(world.Rand.Intn(32), world.Rand.Intn(32), float32(i+10), gender, world.Rand)
		world.Spawn(e)
	}
}

func (world *World) Spawn(e Evolver) {
	go e.Evolve(world.Requests)
	switch obj := e.(type) {
	case *Creature:
		r := NewPutRequest(obj)
		world.Requests <- r
	}
}

func (world *World) handle(req Requester) {
	// this checks should move somewhere else....
	if req.X() <= -1 || req.X() >= world.W {
		log.Printf("Out of bounds x %d reject request!", req.X())
		return
	}
	if req.Y() <= -1 || req.Y() >= world.H {
		log.Printf("Out of bounds y %d reject request!", req.Y())
		return
	}
	// dispatch requests base on their type to the corresponding handler.
	switch r := req.(type) {
	case *PutRequest:
		world.handlePut(r)
	case *PostRequest:
		world.handlePost(r)
	case *DeleteRequest:
		world.handleDelete(r)
	}
}

func (world *World) handlePut(req *PutRequest) {
	if world.Evolvers.Contains(req.Obj()) {
		log.Println("Cannot Put this evolvers because it is already inside... Use Post to update...")
		return
	}
	world.Evolvers = append(world.Evolvers, req.Obj())
	switch o := req.Obj().(type) {
	case *Creature:
		world.Creatures = append(world.Creatures, o)
		world.Map[o.X()][o.Y()] = o
		go o.Evolve(world.Requests)
	}
}

func (world *World) handlePost(req *PostRequest) {
	if world.Evolvers.Contains(req.Obj()) == false {
		world.handlePut(NewPutRequest(req.Obj()))
	}
	switch o := req.Obj().(type) {
	case *Creature:
		if world.Map[req.X()][req.Y()] != nil {
			return
		}
		o.SetX(req.X())
		o.SetY(req.Y())
		world.Map[o.X()][o.Y()] = o
		world.Map[o.X()][o.Y()] = nil
	}
}

func (world *World) handleDelete(req *DeleteRequest) {
	var i int
	var t Evolver
	for i, t = range world.Evolvers {
		if t == req.Obj() {
			break
		}
	}
	world.Evolvers[i] = nil
	world.Evolvers = append(world.Evolvers[:i], world.Evolvers[i+1:]...)
	switch o := req.Obj().(type) {
	case *Creature:
		world.Map[o.X()][o.Y()] = nil
		world.removeCreature(o)
	}
}

func (world *World) removeCreature(creature *Creature) {
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

type Request struct {
	obj Evolver
}

type PutRequest struct {
	Request
}

func NewPutRequest(obj Evolver) *PutRequest {
	r := new(PutRequest)
	r.obj = obj
	return r
}

func (p PutRequest) Obj() Evolver { return p.obj }
func (p PutRequest) X() int       { return p.obj.X() }
func (p PutRequest) Y() int       { return p.obj.Y() }

type PostRequest struct {
	Request
	x, y int
}

func NewPostRequest(obj Evolver, newx, newy int) *PostRequest {
	r := new(PostRequest)
	r.obj = obj
	r.x = newx
	r.y = newy
	return r
}

func (p PostRequest) Obj() Evolver { return p.obj }
func (p PostRequest) X() int       { return p.x }
func (p PostRequest) Y() int       { return p.y }

type DeleteRequest struct {
	Request
}

func NewDeleteRequest(obj Evolver, x, y int) *DeleteRequest {
	r := new(DeleteRequest)
	r.obj = obj
	return r
}

func (p DeleteRequest) Obj() Evolver { return p.obj }
func (p DeleteRequest) X() int       { return p.obj.X() }
func (p DeleteRequest) Y() int       { return p.obj.Y() }
