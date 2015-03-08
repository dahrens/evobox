package evolution

import "math/rand"

const (
	FIELDS            = 12
	FIELDS_WIDTH_MIN  = 2
	FIELDS_HEIGHT_MIN = 2
	FIELDS_WIDTH_MAX  = 4
	FIELDS_HEIGHT_MAX = 4
)

type Terrains []string

type Field struct {
	Passable bool
	Terrains []string
}

func NewField(terrains []string) *Field {
	field := new(Field)
	field.Terrains = terrains
	field.Passable = true
	return field
}

type Fields []*Field

type Plan []Fields

func NewPlan(w, h int) Plan {
	plan := make(Plan, w)
	for x := 0; x < w; x++ {
		plan[x] = make(Fields, h)
	}
	plan.generate(w, h)
	return plan
}

func (plan Plan) generate(w, h int) {
	for x, row := range plan {
		for y, _ := range row {
			terrains := Terrains{"Tiles/grass1.png"}
			plan[x][y] = NewField(terrains)
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
					plan[x][y] = NewField(Terrains{"Tiles/dirt1.png", "Planter og Farm stuff/wheat4.png"})
				}
			}
		}
	}
}
