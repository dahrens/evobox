package evolution

import (
	"log"
	"sync"
	"time"
)

type Requester interface {
	Obj() Evolver
	X() int
	Y() int
}

type World struct {
	W         int
	H         int
	Evolvers  Evolvers `json:"-"`
	Flowers   Flowers
	Creatures Creatures
	Requests  chan Requester `json:"-"`
	Clock     *time.Ticker   `json:"-"`
	Tick      int
	Speed     time.Duration
	Plan      *Plan
	Client    *Client `json:"-"`
	running   bool
}

func NewWorld(w, h int, client *Client) *World {
	world := new(World)
	world.W = w
	world.H = h
	world.Evolvers = make(Evolvers, 0)
	world.Creatures = make(Creatures, 0)
	world.Flowers = make(Flowers, 0)
	world.Requests = make(chan Requester)
	world.Speed = 400 * time.Millisecond
	world.Clock = time.NewTicker(world.Speed)
	world.Tick = 0
	world.Plan = NewPlan(w, h, world)
	world.running = false
	world.Client = client
	go world.serve()
	return world
}

func (world *World) Init(initialCreatures, initialFlowers int) {
	world.SpawnManyCreatures(initialCreatures/2, GENDER_MALE)
	world.SpawnManyCreatures(initialCreatures/2, GENDER_FEMALE)
	world.SpawnManyFlowers(initialFlowers/3, "flower-orange.png")
	world.SpawnManyFlowers(initialFlowers/3, "flower-red.png")
	world.SpawnManyFlowers(initialFlowers/3, "flower-yellow.png")
}

func (world *World) SpawnManyFlowers(n int, color string) {
	for i := 0; i < n; i++ {
		e := NewFlower(random(300, 1800), random(300, 1800), color, world)
		world.Spawn(e)
	}
}

func (world *World) SpawnManyCreatures(n int, gender Gender) {
	for i := 0; i < n; i++ {
		e := NewCreature(float32(i+10), gender, world.Client)
		world.Spawn(e)
	}
}

func (world *World) Spawn(e Evolver) {
	switch obj := e.(type) {
	case Evolver:
		r := NewPutRequest(obj)
		world.Requests <- r
	}
}

func (world *World) Run() {
	world.running = true
	go world.pulse()
}

func (world *World) Pause() {
	world.running = false
	world.Clock.Stop()
	world.Clock = time.NewTicker(world.Speed)
}

func (world *World) Reset(tick_interval, map_width, map_height int) {
	world.Speed = time.Duration(tick_interval) * time.Millisecond
	world.W = map_width
	world.H = map_height
	if world.running {
		world.Pause()
	} else {
		world.Clock = time.NewTicker(world.Speed)
	}
	world.Evolvers = make(Evolvers, 0)
	world.Creatures = make(Creatures, 0)
	world.Flowers = make(Flowers, 0)
	world.Plan = NewPlan(map_width, map_height, world)
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
		var w sync.WaitGroup
		w.Add(len(t))
		for _, e := range t {
			if e.Alive() {
				e.Pulse() <- &Tick{Count: world.Tick, Wait: &w}
			}
		}
		w.Wait()
		//world.Client.Write(NewMessage("update-creatures", world.Creatures))
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
		go o.Evolve(world)
	case *Flower:
		world.Flowers = append(world.Flowers, o)
		go o.Evolve(world)
	}
}

func (world *World) handlePost(req *PostRequest) {
	if world.Evolvers.Contains(req.Obj()) == false {
		world.handlePut(NewPutRequest(req.Obj()))
	}
	switch o := req.Obj().(type) {
	case *Creature:
		o.SetX(req.X())
		o.SetY(req.Y())
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
		world.removeCreature(o)
	case *Flower:
		world.removeFlower(o)
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

func (world *World) removeFlower(flower *Flower) {
	var i int
	var t *Flower
	for i, t = range world.Flowers {
		if t == flower {
			break
		}
	}
	world.Flowers[i] = nil
	world.Flowers = append(world.Flowers[:i], world.Flowers[i+1:]...)
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
func (p PutRequest) X() int       { return p.obj.GetX() }
func (p PutRequest) Y() int       { return p.obj.GetY() }

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
func (p DeleteRequest) X() int       { return p.obj.GetX() }
func (p DeleteRequest) Y() int       { return p.obj.GetY() }
