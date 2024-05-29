package time_units

type TimeUnit struct {
	Label        string
	Abbreviation string
}

var (
	timeUnits = []TimeUnit{
		{"second", "s"},
		{"seconds", "secs"},
		{"minute", "min"},
		{"minutes", "mins"},
		{"hour", "h"},
		{"hours", "hs"},
	}
	byLabel        = make(map[string]TimeUnit)
	byAbbreviation = make(map[string]TimeUnit)
)

func init() {
	for _, unit := range timeUnits {
		byLabel[unit.Label] = unit
		byAbbreviation[unit.Abbreviation] = unit
	}
}

func ValueOfLabel(label string) (TimeUnit, bool) {
	unit, ok := byLabel[label]
	return unit, ok
}

func ValueOfAbbreviation(abbreviation string) (TimeUnit, bool) {
	unit, ok := byAbbreviation[abbreviation]
	return unit, ok
}
