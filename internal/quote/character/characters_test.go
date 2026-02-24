package character

import (
	"testing"
)

func TestGetCharacterName_AllEntries(t *testing.T) {
	expect := map[Character]string{
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

	for ch, wantName := range expect {
		got := CharacterNames.GetCharacterName(ch)
		if got != wantName {
			t.Errorf("GetCharacterName(%q): got %q, want %q", ch, got, wantName)
		}
	}

	if len(CharacterNames) != len(expect) {
		t.Errorf("CharacterNames has %d entries, test expects %d — update the test", len(CharacterNames), len(expect))
	}
}

func TestGetCharacterName_UnknownID(t *testing.T) {
	unknowns := []Character{"", Character("100"), Character("abc"), Character("-1"), Character("61")}
	for _, ch := range unknowns {
		got := CharacterNames.GetCharacterName(ch)
		if got != "Unknown" {
			t.Errorf("GetCharacterName(%q): got %q, want \"Unknown\"", ch, got)
		}
	}
}

func TestGetAllCharacters_ReturnsCopy(t *testing.T) {
	all := CharacterNames.GetAllCharacters()

	if len(all) != len(CharacterNames) {
		t.Fatalf("GetAllCharacters returned %d entries, want %d", len(all), len(CharacterNames))
	}
	for ch, wantName := range CharacterNames {
		if all[ch] != wantName {
			t.Errorf("GetAllCharacters()[%q]: got %q, want %q", ch, all[ch], wantName)
		}
	}

	all[Character("test_mutation")] = "should not leak"
	if _, exists := CharacterNames[Character("test_mutation")]; exists {
		t.Error("GetAllCharacters did not return a copy — mutation leaked into CharacterNames")
	}
}

func TestCharacterFromID(t *testing.T) {
	if CharacterFromID("10") != Battler {
		t.Errorf("CharacterFromID(\"10\"): got %q, want %q", CharacterFromID("10"), Battler)
	}
	if CharacterFromID("27") != Beatrice {
		t.Errorf("CharacterFromID(\"27\"): got %q, want %q", CharacterFromID("27"), Beatrice)
	}
	if CharacterFromID("narrator") != Narrator {
		t.Errorf("CharacterFromID(\"narrator\"): got %q, want %q", CharacterFromID("narrator"), Narrator)
	}
	if CharacterFromID("unknown") != Character("unknown") {
		t.Errorf("CharacterFromID(\"unknown\"): got %q, want %q", CharacterFromID("unknown"), Character("unknown"))
	}
}
