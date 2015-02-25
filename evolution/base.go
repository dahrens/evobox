package evolution

type Evolver interface {
	X() int
	Y() int
	SetX(int)
	SetY(int)
	Evolve(worldRequests chan Requester)
	Pulse() chan int
	Alive() bool
}

type Position struct {
	x int
	y int
}

type Fragment struct {
	Position
	Age   int
	pulse chan int
}

type Plant struct {
	Fragment
	NutritionalValue float32
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
