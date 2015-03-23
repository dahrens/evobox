package evolution

import "log"

const (
	WOODS = 8
	TREES = 10
)

type Coordinate struct {
	Point
	Passable  bool
	Fragments []Fragmenter
}

func NewCoordinate(x, y int, passable bool) *Coordinate {
	coordinate := new(Coordinate)
	coordinate.X = x
	coordinate.Y = y
	coordinate.Passable = passable
	coordinate.Fragments = make([]Fragmenter, 0)
	return coordinate
}

func (coordinate *Coordinate) addFragment(fragment Fragmenter) {
	coordinate.Fragments = append(coordinate.Fragments, fragment)
}

func (coordinate *Coordinate) deleteFragment(fragment Fragmenter) {
	if len(coordinate.Fragments) == 0 {
		return
	}
	var i int
	var f Fragmenter
	for i, f = range coordinate.Fragments {
		if f == fragment {
			break
		}
	}
	coordinate.Fragments[i] = nil
	coordinate.Fragments = append(coordinate.Fragments[:i], coordinate.Fragments[i+1:]...)
}

type Line []*Coordinate

type Plan struct {
	Map       []Line `json:"-"`
	Fragments []Fragmenter
	world     *World
}

func NewPlan(w, h int, world *World) *Plan {
	plan := new(Plan)
	plan.Map = make([]Line, w)
	plan.Fragments = make([]Fragmenter, 0)
	plan.world = world
	for x := 0; x < w; x++ {
		plan.Map[x] = make(Line, h)
		for y := range plan.Map[x] {
			plan.Map[x][y] = NewCoordinate(x, y, true)
		}
	}
	plan.generateStaticFragments(w, h)
	return plan
}

func (plan *Plan) generateStaticFragments(w, h int) {
	for i := 0; i < WOODS; i++ {
		init_x := random(120, 1800)
		start_x := init_x
		start_y := random(300, 1500)
		for j := 0; j < TREES; j++ {
			t := NewTree(start_x, start_y)
			plan.addFragment(t)
			if j%3 == 0 {
				start_y += t.GetH()
				start_x = init_x
			} else {
				start_x += t.GetW() + random(-5, 25)
				start_y += random(-25, 25)
			}
		}
	}
}

func (plan *Plan) IsInside(x, y int) bool {
	if x <= -1 || x >= plan.world.W {
		return false
	}
	if y <= -1 || y >= plan.world.H {
		return false
	}
	return true
}

func (plan *Plan) addFragment(fragment Fragmenter) {
	plan.Fragments = append(plan.Fragments, fragment)
	plan.Map[fragment.GetX()][fragment.GetY()].addFragment(fragment)
	log.Printf("Added %#v\n", fragment)
}

func (plan *Plan) updateFragment(fragment Fragmenter, newX, newY int) {
	plan.Map[fragment.GetX()][fragment.GetY()].deleteFragment(fragment)
	fragment.SetX(newX)
	fragment.SetY(newY)
	plan.Map[fragment.GetX()][fragment.GetY()].addFragment(fragment)
}

func (plan *Plan) deleteFragment(fragment Fragmenter) {
	var i int
	var f Fragmenter
	for i, f = range plan.Fragments {
		if f == fragment {
			break
		}
	}
	plan.Fragments[i] = nil
	plan.Fragments = append(plan.Fragments[:i], plan.Fragments[i+1:]...)
	plan.Map[fragment.GetX()][fragment.GetY()].deleteFragment(fragment)
}

func (plan *Plan) getFragmentsInCircle(startx, starty, radius int) []Fragmenter {
	objects := make([]Fragmenter, 0)
	// http://stackoverflow.com/a/15856549
	for x := startx - radius; x < startx; x++ {
		for y := starty - radius; y < starty; y++ {
			// we don't have to take the square root, it's slow
			if (x-startx)*(x-startx)+(y-starty)*(y-starty) <= radius*radius {
				xSym := startx - (x - startx)
				ySym := starty - (y - starty)
				// (x, y), (x, ySym), (xSym , y), (xSym, ySym) are in the circle
				if plan.IsInside(x, y) {
					objects = append(objects, plan.Map[x][y].Fragments...)
				}
				if plan.IsInside(x, ySym) {
					objects = append(objects, plan.Map[x][ySym].Fragments...)
				}
				if plan.IsInside(xSym, y) {
					objects = append(objects, plan.Map[xSym][y].Fragments...)
				}
				if plan.IsInside(xSym, ySym) {
					objects = append(objects, plan.Map[xSym][ySym].Fragments...)
				}
			}
		}
	}
	return objects
}
