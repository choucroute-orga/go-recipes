package units

type Unit struct {
	Label        string
	Abbreviation string
}

var (
	units = []Unit{
		{"item", "i"},
		{"items", "is"},
		{"cup", "c"},
		{"cups", "cs"},
		{"tablespoon", "tbsp"},
		{"teaspoon", "tsp"},
		{"gram", "g"},
		{"grams", "g"},
		{"kilogram", "kg"},
		{"kilograms", "kg"},
	}
	byLabel        = make(map[string]Unit)
	byAbbreviation = make(map[string]Unit)
)

func init() {
	for _, unit := range units {
		byLabel[unit.Label] = unit
		byAbbreviation[unit.Abbreviation] = unit
	}
}

func ValueOfLabel(label string) (Unit, bool) {
	unit, ok := byLabel[label]
	return unit, ok
}

func ValueOfAbbreviation(abbreviation string) (Unit, bool) {
	unit, ok := byAbbreviation[abbreviation]
	return unit, ok
}
