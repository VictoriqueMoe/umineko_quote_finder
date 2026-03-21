package language

type Language string

const (
	Auto       Language = "auto"
	English    Language = "en"
	WitchHunt  Language = "wh"
	Japanese   Language = "ja"
	Russian    Language = "ru"
	Spanish    Language = "es"
	Portuguese Language = "pt"
)

var All = []Language{English, WitchHunt, Japanese, Russian, Spanish, Portuguese}

func (Language) Parse(s string) Language {
	switch s {
	case "auto":
		return Auto
	case "en":
		return English
	case "ja":
		return Japanese
	case "es":
		return Spanish
	case "wh":
		return WitchHunt
	case "ru":
		return Russian
	case "pt":
		return Portuguese
	default:
		return English
	}
}

func (l Language) String() string {
	return string(l)
}
