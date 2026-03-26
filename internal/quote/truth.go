package quote

type Truth string

const (
	TruthAll    Truth = ""
	TruthRed    Truth = "red"
	TruthBlue   Truth = "blue"
	TruthGold   Truth = "gold"
	TruthPurple Truth = "purple"
)

func (Truth) Parse(s string) Truth {
	switch s {
	case "red":
		return TruthRed
	case "blue":
		return TruthBlue
	case "gold":
		return TruthGold
	case "purple":
		return TruthPurple
	default:
		return TruthAll
	}
}
