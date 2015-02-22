package evolution

import "sort"
import "strings"

type Creatures []*Creature

func (c Creatures) Sort(column, direction string) {
	s := strings.Join(append(make([]string, 2), column, direction), "")
	switch s {
	case "Healthasc":
		sort.Sort(ByHealthAsc(c))
	case "Healthdesc":
		sort.Sort(ByHealthDesc(c))
	case "Nameasc":
		sort.Sort(ByNameAsc(c))
	case "Namedesc":
		sort.Sort(ByNameDesc(c))
	case "Genderasc":
		sort.Sort(ByGenderAsc(c))
	case "Genderdesc":
		sort.Sort(ByGenderDesc(c))
	case "Ageasc":
		sort.Sort(ByAgeAsc(c))
	case "Agedesc":
		sort.Sort(ByAgeDesc(c))
	case "Hungerasc":
		sort.Sort(ByHungerAsc(c))
	case "Hungerdesc":
		sort.Sort(ByHungerDesc(c))
	case "Libidoasc":
		sort.Sort(ByLibidoAsc(c))
	case "Libidodesc":
		sort.Sort(ByLibidoDesc(c))
	}
}

func (c Creatures) ToMap(start, end int) []interface{} {
	data := make([]interface{}, end-start)
	if start < len(c) {
		p := 0
		for i := start; i < end; i++ {
			e := c[i]
			record := make(map[string]interface{})
			record["Name"] = e.Name
			record["Age"] = e.Age
			record["Health"] = e.Health
			record["Libido"] = e.Libido
			record["Hunger"] = e.Hunger
			record["Gender"] = e.Gender
			record["X"] = e.X()
			record["Y"] = e.Y()
			data[p] = record
			p++
		}
	}
	return data
}

type ByNameAsc Creatures

func (a ByNameAsc) Len() int           { return len(a) }
func (a ByNameAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNameAsc) Less(i, j int) bool { return a[i].Name < a[j].Name }

type ByNameDesc Creatures

func (a ByNameDesc) Len() int           { return len(a) }
func (a ByNameDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByNameDesc) Less(i, j int) bool { return a[i].Name > a[j].Name }

type ByGenderAsc Creatures

func (a ByGenderAsc) Len() int           { return len(a) }
func (a ByGenderAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenderAsc) Less(i, j int) bool { return a[i].Gender < a[j].Gender }

type ByGenderDesc Creatures

func (a ByGenderDesc) Len() int           { return len(a) }
func (a ByGenderDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByGenderDesc) Less(i, j int) bool { return a[i].Gender > a[j].Gender }

type ByAgeAsc Creatures

func (a ByAgeAsc) Len() int           { return len(a) }
func (a ByAgeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAgeAsc) Less(i, j int) bool { return a[i].Age < a[j].Age }

type ByAgeDesc Creatures

func (a ByAgeDesc) Len() int           { return len(a) }
func (a ByAgeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAgeDesc) Less(i, j int) bool { return a[i].Age > a[j].Age }

type ByHungerAsc Creatures

func (a ByHungerAsc) Len() int           { return len(a) }
func (a ByHungerAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHungerAsc) Less(i, j int) bool { return a[i].Hunger < a[j].Hunger }

type ByHungerDesc Creatures

func (a ByHungerDesc) Len() int           { return len(a) }
func (a ByHungerDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHungerDesc) Less(i, j int) bool { return a[i].Hunger > a[j].Hunger }

type ByLibidoAsc Creatures

func (a ByLibidoAsc) Len() int           { return len(a) }
func (a ByLibidoAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLibidoAsc) Less(i, j int) bool { return a[i].Libido < a[j].Libido }

type ByLibidoDesc Creatures

func (a ByLibidoDesc) Len() int           { return len(a) }
func (a ByLibidoDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLibidoDesc) Less(i, j int) bool { return a[i].Libido > a[j].Libido }

type ByHealthAsc Creatures

func (a ByHealthAsc) Len() int           { return len(a) }
func (a ByHealthAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHealthAsc) Less(i, j int) bool { return a[i].Health < a[j].Health }

type ByHealthDesc Creatures

func (a ByHealthDesc) Len() int           { return len(a) }
func (a ByHealthDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHealthDesc) Less(i, j int) bool { return a[i].Health > a[j].Health }
