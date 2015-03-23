package evolution

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

func (tree *Tree) Collides() bool { return true }
