package evolution

const (
	WOODS   = 3
	TREES   = 10
	FLOWERS = 150
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

type PlaceableFragment struct {
	Fragment
	Sprite string
	Sheet  string
}

type Tree PlaceableFragment

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

type Flower PlaceableFragment

func NewFlower(x, y int, sprite string) *Tree {
	flower := new(Tree)
	flower.X = x
	flower.Y = y
	flower.W = 39
	flower.H = 69
	flower.Anchor.X = 0.5
	flower.Anchor.Y = 0.5
	flower.Sprite = sprite
	flower.Sheet = "default"
	return flower
}

type Line []Coordinate

type Plan struct {
	_m        []Line
	Fragments []Fragmenter
	Evolvers  Evolvers
}

func NewPlan(w, h int) *Plan {
	plan := new(Plan)
	plan._m = make([]Line, w)
	for x := 0; x < w; x++ {
		plan._m[x] = make(Line, h)
	}
	plan.generateStaticFragments(w, h)
	return plan
}

func (plan *Plan) generateStaticFragments(w, h int) {
	// for x, line := range plan._m {
	// 	for y, _ := range line {
	// 		passable := true
	// 		plan._m[x][y] = NewCoordinate(x, y, passable)
	// 	}
	// }

	for i := 0; i < FLOWERS; i++ {
		var colored_name string
		switch random(1, 4) {
		case 1:
			colored_name = "flower-red.png"
		case 2:
			colored_name = "flower-yellow.png"
		case 3:
			colored_name = "flower-orange.png"
		}
		t := NewFlower(random(300, 1800), random(300, 1800), colored_name)
		plan.addFragment(t)
	}

	for i := 0; i < WOODS; i++ {
		init_x := random(300, 1800)
		start_x := init_x
		start_y := random(300, 1800)
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
