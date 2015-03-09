package evolution

import (
	"math/rand"
	"sync"
	"time"
)

type Tick struct {
	Count int
	Wait  *sync.WaitGroup
}

type Point struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

type Anchor struct {
	X float32
	Y float32
}

type Fragment struct {
	Point
	Size
	Anchor Anchor
}

type EvolverFragment struct {
	Fragment
	Age   int
	Birth int
	pulse chan *Tick
}

type Fragmenter interface {
	GetX() int
	GetY() int
	SetX(int)
	SetY(int)
	GetW() int
	GetH() int
	SetW(int)
	SetH(int)
}

type Evolver interface {
	Fragmenter
	Evolve(world *World)
	Pulse() chan *Tick
	Alive() bool
}

type Evolvers []Evolver

func (evo Evolvers) Contains(e Evolver) bool {
	for _, b := range evo {
		if b == e {
			return true
		}
	}
	return false
}

// https://bitbucket.org/mikespook/golib/src/46a4f2a8abcb/autoinc/

type AutoInc struct {
	start, step int
	queue       chan int
	running     bool
}

func NewAutoInc(start, step int) (ai *AutoInc) {
	ai = &AutoInc{
		start:   start,
		step:    step,
		running: true,
		queue:   make(chan int, 4),
	}
	go ai.process()
	return
}

func (ai *AutoInc) process() {
	defer func() { recover() }()
	for i := ai.start; ai.running; i = i + ai.step {
		ai.queue <- i
	}
}

func (ai *AutoInc) Id() int {
	return <-ai.queue
}

func (ai *AutoInc) Close() {
	ai.running = false
	close(ai.queue)
}

// http://golangcookbook.blogspot.de/2012/11/generate-random-number-in-given-range.html

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
