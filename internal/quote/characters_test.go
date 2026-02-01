package quote

import (
	"testing"
)

func TestGetCharacterName_AllEntries(t *testing.T) {
	expect := map[string]string{
		"00":       "GroupVoices",
		"01":       "Kinzo",
		"02":       "Krauss",
		"03":       "Natsuhi",
		"04":       "Jessica",
		"05":       "Eva",
		"06":       "Hideyoshi",
		"07":       "George",
		"08":       "Rudolf",
		"09":       "Kyrie",
		"10":       "Battler",
		"11":       "Ange",
		"12":       "Rosa",
		"13":       "Maria",
		"14":       "Genji",
		"15":       "Shannon",
		"16":       "Kanon",
		"17":       "Gohda",
		"18":       "KumasawaChiyo",
		"19":       "NanjoTerumasa",
		"20":       "Amakusa",
		"21":       "Okonogi",
		"22":       "Kasumi",
		"23":       "ProfessorOotsuki",
		"24":       "CaptainKawabata",
		"25":       "NanjoMasayuki",
		"26":       "KumasawaSabakichi",
		"27":       "Beatrice",
		"28":       "Bernkastel",
		"29":       "Lambdadelta",
		"30":       "Virgilia",
		"31":       "Ronove",
		"32":       "Gaap",
		"33":       "Sakutarou",
		"34":       "Evatrice",
		"35":       "Chiester45",
		"36":       "Chiester410",
		"37":       "Chiester00",
		"38":       "Lucifer",
		"39":       "Leviathan",
		"40":       "Satan",
		"41":       "Belphegor",
		"42":       "Mammon",
		"43":       "Beelzebub",
		"44":       "Asmodeus",
		"45":       "Goat",
		"46":       "Erika",
		"47":       "Dlanor",
		"48":       "Gertrude",
		"49":       "Cornelia",
		"50":       "Featherine",
		"51":       "Zepar",
		"52":       "Furfur",
		"53":       "Lion",
		"54":       "Willard",
		"55":       "Claire",
		"56":       "Ikuko",
		"57":       "Tohya",
		"58":       "KinzoYoung",
		"59":       "BiceChickBeato",
		"60":       "BeatoElder",
		"99":       "MiscVoices",
		"narrator": "Narrator",
	}

	for id, wantName := range expect {
		got := CharacterNames.GetCharacterName(id)
		if got != wantName {
			t.Errorf("GetCharacterName(%q): got %q, want %q", id, got, wantName)
		}
	}

	if len(CharacterNames) != len(expect) {
		t.Errorf("CharacterNames has %d entries, test expects %d — update the test", len(CharacterNames), len(expect))
	}
}

func TestGetCharacterName_UnknownID(t *testing.T) {
	unknowns := []string{"", "100", "abc", "-1", "61"}
	for _, id := range unknowns {
		got := CharacterNames.GetCharacterName(id)
		if got != "Unknown" {
			t.Errorf("GetCharacterName(%q): got %q, want \"Unknown\"", id, got)
		}
	}
}

func TestGetAllCharacters_ReturnsCopy(t *testing.T) {
	all := CharacterNames.GetAllCharacters()

	if len(all) != len(CharacterNames) {
		t.Fatalf("GetAllCharacters returned %d entries, want %d", len(all), len(CharacterNames))
	}
	for id, wantName := range CharacterNames {
		if all[id] != wantName {
			t.Errorf("GetAllCharacters()[%q]: got %q, want %q", id, all[id], wantName)
		}
	}

	all["test_mutation"] = "should not leak"
	if _, exists := CharacterNames["test_mutation"]; exists {
		t.Error("GetAllCharacters did not return a copy — mutation leaked into CharacterNames")
	}
}
