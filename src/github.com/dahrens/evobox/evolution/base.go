package evolution

type Evolver interface {
	X() int
	Y() int
	Evolve()
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
