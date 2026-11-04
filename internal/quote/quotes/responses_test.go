package quotes

import "testing"

func TestNewCharactersResult_ExcludesCuratedAndEmpty(t *testing.T) {
	curated := map[string]string{"miyao": "Miyao"}
	speakers := map[string]string{"miyao": "Mitake Miyao", "misc": "Misc", "": "Nobody"}

	result := NewCharactersResult(curated, speakers)

	if result.Characters["miyao"] != "Miyao" {
		t.Errorf("curated name: got %q, want Miyao", result.Characters["miyao"])
	}
	if _, ok := result.Additional["miyao"]; ok {
		t.Error("curated character should not be listed as additional")
	}
	if _, ok := result.Additional[""]; ok {
		t.Error("empty character id should not be listed as additional")
	}
	if result.Additional["misc"] != "Misc" {
		t.Errorf("additional name: got %q, want Misc", result.Additional["misc"])
	}
}
