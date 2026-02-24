package character

type (
	Character        string
	characterMapping map[Character]string
)

const (
	GroupVoices Character = "group_voices"
	Kinzo       Character = "kinzo"
	Krauss      Character = "krauss"
	Natsuhi     Character = "natsuhi"
	Jessica     Character = "jessica"
	Eva         Character = "eva"
	Hideyoshi   Character = "hideyoshi"
	George      Character = "george"
	Rudolf      Character = "rudolf"
	Kyrie       Character = "kyrie"
	Battler     Character = "battler"
	Ange        Character = "ange"
	Rosa        Character = "rosa"
	Maria       Character = "maria"
	Genji       Character = "genji"
	Shannon     Character = "shannon"
	Kanon       Character = "kanon"
	Gohda       Character = "gohda"
	Kumasawa    Character = "kumasawa"
	Nanjo       Character = "nanjo"
	Amakusa     Character = "amakusa"
	Okonogi     Character = "okonogi"
	Kasumi      Character = "kasumi"
	Professor   Character = "professor"
	Kawabata    Character = "kawabata"
	NanjoSon    Character = "nanjo_son"
	KumasawaSon Character = "kumasawa_son"
	Beatrice    Character = "beatrice"
	Bernkastel  Character = "bernkastel"
	Lambdadelta Character = "lambdadelta"
	Virgilia    Character = "virgilia"
	Ronove      Character = "ronove"
	Gaap        Character = "gaap"
	Sakutarou   Character = "sakutarou"
	EvaBeatrice Character = "eva_beatrice"
	Chiester45  Character = "chiester_45"
	Chiester410 Character = "chiester_410"
	Chiester00  Character = "chiester_00"
	Lucifer     Character = "lucifer"
	Leviathan   Character = "leviathan"
	Satan       Character = "satan"
	Belphegor   Character = "belphegor"
	Mammon      Character = "mammon"
	Beelzebub   Character = "beelzebub"
	Asmodeus    Character = "asmodeus"
	Goat        Character = "goat"
	Erika       Character = "erika"
	Dlanor      Character = "dlanor"
	Gertrude    Character = "gertrude"
	Cornelia    Character = "cornelia"
	Featherine  Character = "featherine"
	Zepar       Character = "zepar"
	Furfur      Character = "furfur"
	Lion        Character = "lion"
	Will        Character = "will"
	Clair       Character = "clair"
	Ikuko       Character = "ikuko"
	Tohya       Character = "tohya"
	KinzoYoung  Character = "kinzo_young"
	Bice        Character = "bice"
	BeatoElder  Character = "beato_elder"
	MiscVoices  Character = "misc_voices"
	Narrator    Character = "narrator"
)

var (
	characterIDs = map[Character]string{
		GroupVoices: "00",
		Kinzo:       "01",
		Krauss:      "02",
		Natsuhi:     "03",
		Jessica:     "04",
		Eva:         "05",
		Hideyoshi:   "06",
		George:      "07",
		Rudolf:      "08",
		Kyrie:       "09",
		Battler:     "10",
		Ange:        "11",
		Rosa:        "12",
		Maria:       "13",
		Genji:       "14",
		Shannon:     "15",
		Kanon:       "16",
		Gohda:       "17",
		Kumasawa:    "18",
		Nanjo:       "19",
		Amakusa:     "20",
		Okonogi:     "21",
		Kasumi:      "22",
		Professor:   "23",
		Kawabata:    "24",
		NanjoSon:    "25",
		KumasawaSon: "26",
		Beatrice:    "27",
		Bernkastel:  "28",
		Lambdadelta: "29",
		Virgilia:    "30",
		Ronove:      "31",
		Gaap:        "32",
		Sakutarou:   "33",
		EvaBeatrice: "34",
		Chiester45:  "35",
		Chiester410: "36",
		Chiester00:  "37",
		Lucifer:     "38",
		Leviathan:   "39",
		Satan:       "40",
		Belphegor:   "41",
		Mammon:      "42",
		Beelzebub:   "43",
		Asmodeus:    "44",
		Goat:        "45",
		Erika:       "46",
		Dlanor:      "47",
		Gertrude:    "48",
		Cornelia:    "49",
		Featherine:  "50",
		Zepar:       "51",
		Furfur:      "52",
		Lion:        "53",
		Will:        "54",
		Clair:       "55",
		Ikuko:       "56",
		Tohya:       "57",
		KinzoYoung:  "58",
		Bice:        "59",
		BeatoElder:  "60",
		MiscVoices:  "99",
	}
	idToCharacter  map[string]Character
	CharacterNames = characterMapping{
		GroupVoices: "Group Voices",
		Kinzo:       "Ushiromiya Kinzo",
		Krauss:      "Ushiromiya Krauss",
		Natsuhi:     "Ushiromiya Natsuhi",
		Jessica:     "Ushiromiya Jessica",
		Eva:         "Ushiromiya Eva",
		Hideyoshi:   "Ushiromiya Hideyoshi",
		George:      "Ushiromiya George",
		Rudolf:      "Ushiromiya Rudolf",
		Kyrie:       "Ushiromiya Kyrie",
		Battler:     "Ushiromiya Battler",
		Ange:        "Ushiromiya Ange",
		Rosa:        "Ushiromiya Rosa",
		Maria:       "Ushiromiya Maria",
		Genji:       "Ronoue Genji",
		Shannon:     "Shannon",
		Kanon:       "Kanon",
		Gohda:       "Gohda Toshiro",
		Kumasawa:    "Kumasawa Chiyo",
		Nanjo:       "Nanjo Terumasa",
		Amakusa:     "Amakusa Juuza",
		Okonogi:     "Okonogi Tetsuro",
		Kasumi:      "Sumadera Kasumi",
		Professor:   "Professor Ootsuki",
		Kawabata:    "Captain Kawabata",
		NanjoSon:    "Nanjo Masayuki",
		KumasawaSon: "Kumasawa Sabakichi",
		Beatrice:    "Beatrice",
		Bernkastel:  "Bernkastel",
		Lambdadelta: "Lambdadelta",
		Virgilia:    "Virgilia",
		Ronove:      "Ronove",
		Gaap:        "Gaap",
		Sakutarou:   "Sakutarou",
		EvaBeatrice: "Eva Beatrice",
		Chiester45:  "Chiester 45",
		Chiester410: "Chiester 410",
		Chiester00:  "Chiester 00",
		Lucifer:     "Lucifer",
		Leviathan:   "Leviathan",
		Satan:       "Satan",
		Belphegor:   "Belphegor",
		Mammon:      "Mammon",
		Beelzebub:   "Beelzebub",
		Asmodeus:    "Asmodeus",
		Goat:        "Goat",
		Erika:       "Furudo Erika",
		Dlanor:      "Dlanor A. Knox",
		Gertrude:    "Gertrude",
		Cornelia:    "Cornelia",
		Featherine:  "Featherine",
		Zepar:       "Zepar",
		Furfur:      "Furfur",
		Lion:        "Ushiromiya Lion",
		Will:        "Willard H. Wright",
		Clair:       "Clair",
		Ikuko:       "Hachijo Ikuko",
		Tohya:       "Hachijo Tohya",
		KinzoYoung:  "Ushiromiya Kinzo",
		Bice:        "Bice",
		BeatoElder:  "Beato the Elder",
		MiscVoices:  "Misc Voices",
		Narrator:    "Narrator",
	}
)

func init() {
	idToCharacter = make(map[string]Character, len(characterIDs))
	for char, id := range characterIDs {
		idToCharacter[id] = char
	}
	idToCharacter["narrator"] = Narrator
}

func (c Character) ID() string {
	if id, ok := characterIDs[c]; ok {
		return id
	}
	return string(c)
}

func CharacterFromID(id string) Character {
	if c, ok := idToCharacter[id]; ok {
		return c
	}
	return Character(id)
}

func (c characterMapping) GetCharacterName(ch Character) string {
	if name, ok := c[ch]; ok {
		return name
	}
	return "Unknown"
}

func (c characterMapping) GetAllCharacters() map[Character]string {
	out := make(map[Character]string, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}
