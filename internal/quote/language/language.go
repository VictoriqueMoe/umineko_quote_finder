package language

type Language string

const (
	English    Language = "en"
	Japanese   Language = "ja"
	Spanish    Language = "es"
	Portuguese Language = "pt"
)

func (Language) Parse(s string) Language {
	switch s {
	case "en":
		return English
	case "ja":
		return Japanese
	case "es":
		return Spanish
	case "pt":
		return Portuguese
	default:
		return English
	}
}

func (l Language) String() string {
	return string(l)
}
