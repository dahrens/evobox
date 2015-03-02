package evolution

import "math/rand"

const (
	FIELDS            = 20
	FIELDS_WIDTH_MIN  = 5
	FIELDS_HEIGHT_MIN = 5
	FIELDS_WIDTH_MAX  = 10
	FIELDS_HEIGHT_MAX = 10
)

type Terrains []string

type Cell struct {
	Passable bool
	Terrains []string
}

func NewCell(terrains []string) *Cell {
	cell := new(Cell)
	cell.Terrains = terrains
	cell.Passable = true
	return cell
}

type Cells []*Cell

type Plan []Cells

func NewPlan(w, h int) Plan {
	plan := make(Plan, w)
	for x := 0; x < w; x++ {
		plan[x] = make(Cells, h)
	}
	plan.generate(w, h)
	return plan
}

func (plan Plan) generate(w, h int) {
	for x, row := range plan {
		for y, _ := range row {
			terrains := Terrains{"Tiles/grass1.png"}
			plan[x][y] = NewCell(terrains)
		}
	}
	for i := 0; i < FIELDS; i++ {
		width := rand.Intn(FIELDS_WIDTH_MAX-FIELDS_WIDTH_MIN) + FIELDS_WIDTH_MIN
		height := rand.Intn(FIELDS_HEIGHT_MAX-FIELDS_HEIGHT_MIN) + FIELDS_HEIGHT_MIN
		start_x := rand.Intn(w)
		start_y := rand.Intn(h)
		for x := start_x; x < start_x+width; x++ {
			for y := start_y; y < start_y+height; y++ {
				if x < w && y < h {
					plan[x][y] = NewCell(Terrains{"Tiles/dirt1.png", "Planter og Farm stuff/wheat4.png"})
				}
			}
		}
	}
}
