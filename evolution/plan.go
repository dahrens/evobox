package evolution

// import (
// 	"bytes"
// 	"image"
// 	"log"
// )

const (
	WOODS = 8
	TREES = 10
)

type Coordinate struct {
	Point
	Passable bool
}

func NewCoordinate(x, y int, passable bool) Coordinate {
	p := new(Coordinate)
	p.X = x
	p.Y = y
	p.Passable = passable
	return *p
}

type Placeable struct {
	Sprite string
	Sheet  string
}

type Tree struct {
	Placeable
	Fragment
}

func NewTree(x, y int) *Tree {
	tree := new(Tree)
	tree.X = x
	tree.Y = y
	tree.W = 39
	tree.H = 69
	tree.Anchor.X = 0.5
	tree.Anchor.Y = 0.5
	tree.Sprite = "tree1.png"
	tree.Sheet = "default"
	return tree
}

func (tree *Tree) GetX() int { return tree.X }
func (tree *Tree) GetY() int { return tree.Y }
func (tree *Tree) GetW() int { return tree.W }
func (tree *Tree) GetH() int { return tree.H }

func (tree *Tree) SetX(x int) { tree.X = x }
func (tree *Tree) SetY(y int) { tree.Y = y }
func (tree *Tree) SetW(w int) { tree.W = w }
func (tree *Tree) SetH(h int) { tree.H = h }

type Line []Coordinate

type Plan struct {
	_m        []Line
	Fragments []Fragmenter
	Evolvers  Evolvers
	world     *World
}

func NewPlan(w, h int, world *World) *Plan {
	plan := new(Plan)
	plan._m = make([]Line, w)
	plan.world = world
	for x := 0; x < w; x++ {
		plan._m[x] = make(Line, h)
	}
	plan.generateStaticFragments(w, h)
	return plan
}

func (plan *Plan) generateStaticFragments(w, h int) {
	// data, err := Asset("assets/experiment/island-black.jpeg")
	// if err != nil {
	// 	log.Println(err)
	// 	panic("can not load image")
	// }
	// reader := bytes.NewReader(data)
	// img, _, err := image.Decode(reader)
	// if err != nil {
	// 	log.Println(err)
	// 	panic("can not decode image")
	// }
	// bounds := img.Bounds()
	// for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
	// 	for x := bounds.Min.X; x < bounds.Max.X; x++ {
	// 		r, g, b, a := img.At(x, y).RGBA()
	// 		log.Println(r)
	// 		log.Println(g)
	// 		log.Println(b)
	// 		log.Println(a)
	// 		plan._m[x][y] = NewCoordinate(x, y, true)
	// 	}
	// }

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

func (plan *Plan) addFragment(f Fragmenter) (Fragmenter, error) {
	plan.Fragments = append(plan.Fragments, f)
	return f, nil
}
