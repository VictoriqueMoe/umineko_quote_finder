package quote

import (
	"fmt"
	"testing"
)

func TestKanonEp2QuotesNotAttributedToErika(t *testing.T) {
	svc := testService

	kanonID := Kanon.ID()
	erikaID := Erika.ID()

	audioIDs := []string{
		"20600528", "20600529",
		"20600530", "20600531", "20600532", "20600533",
		"20600534", "20600535", "20600536", "20600537", "20600538", "20600539",
		"20600540", "20600541",
		"20600542", "20600543",
	}

	langs := make([]string, 0, len(svc.(*service).quotes))
	for lang := range svc.(*service).quotes {
		langs = append(langs, lang)
	}

	for _, lang := range langs {
		for _, audioID := range audioIDs {
			t.Run(fmt.Sprintf("%s/%s", lang, audioID), func(t *testing.T) {
				q := svc.GetByAudioID(lang, audioID)
				if q == nil {
					t.Fatalf("quote not found for audioID %s in lang %s", audioID, lang)
				}
				if q.CharacterID == erikaID {
					t.Errorf("audioID %s in %s is attributed to Erika (%s), should be Kanon (%s): %q",
						audioID, lang, erikaID, kanonID, q.Text)
				}
				if q.CharacterID != kanonID {
					t.Errorf("audioID %s in %s has characterID %q, expected Kanon (%s)",
						audioID, lang, q.CharacterID, kanonID)
				}
			})
		}
	}
}
