package quote

import (
	"strings"
	"testing"
)

var testParser = NewParser()

func TestPresetColoursParsedFromScript(t *testing.T) {
	p := testParser

	lines := []string{
		"preset_define 0,6,36,#FFFFFF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Japanse text",
		"preset_define 1,1,-1,#FF0000,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Red text",
		"preset_define 2,1,-1,#39C6FF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Blue text",
		"preset_define 7,6,48,#C0FFFF,0,0,0,1,6,#000000,0,0,0,#000000,0,70,1880 ;Chapter/BGM/Hint",
		"preset_define 41,1,-1,#FFAA00,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Gold text",
		"preset_define 42,1,-1,#AA71FF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Purple text",
		"preset_define 3,1,40,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,40,700 ;Menu chars text",
		"preset_define 4,6,30,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,40 ;Menu chars/tips/notes text jp",
		"preset_define 5,1,40,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,20,830 ;Menu tips/notes text",
		"preset_define 13,2,55,#FFFFFF,1,0,0,1,10,#000000,0,0,0,#000000,3,80 ;Menu tips/notes titles",
		"preset_define 11,6,36,#FFFFFF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Menu main buttons",
		"preset_define 12,6,36,#FFFFFF,0,0,0,1,3,#000000,0,0,0,#000000,-2,-1 ;Menu chars kill/change buttons",
		"preset_define 6,6,30,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,70 ;Musicbox bgm titles",
		"preset_define 8,1,75,#FFFFFF,0,0,0,1,10,#000000,0,0,0,#000000,3,80 ;Menu titles and buttons",
		"preset_define 9,1,65,#FFFFFF,0,0,0,1,10,#000000,0,0,0,#000000,3,80 ;Menu titles and buttons (smaller)",
		"preset_define 10,8,75,#FFFFFF,0,0,0,1,10,#000000,0,0,0,#000000,0,80,1920 ;Menu first setting line",
		"preset_define 14,1,30,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,40,1920 ;Menu jump portrait titles",
		"preset_define 15,1,42,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,40,1920 ;Menu jump right titles and confirmation",
		"preset_define 16,8,30,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,0,40,1920 ;Menu jump portrait line",
		"preset_define 18,1,34,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,20,1700 ;Trophy description",
		"preset_define 20,8,40,#FFFFFF,0,0,0,1,10,#000000,0,0,0,#000000,3,40,1920 ;Credits normal",
		"preset_define 21,9,90,#000000,0,0,0,1,10,#fae8e8,0,0,0,#000000,6,140,1920 ;Credits h1",
		"preset_define 22,8,60,#a90303,0,0,0,1,10,#000000,0,0,0,#000000,3,120,1920 ;Credits h2",
		"preset_define 23,8,50,#a90303,0,0,0,1,10,#000000,0,0,0,#000000,3,50,1920 ;Credits h3",
		"preset_define 24,8,60,#a90303,0,0,0,1,10,#000000,0,0,0,#000000,0,120,1920 ;Credits h2-sep",
		"preset_define 25,5,47,#FFFFFF,0,0,0,1,10,#000000,0,0,0,#000000,3,40,1920 ;Credits MVP",
		"preset_define 30,2,45,#FFFFFF,0,0,0,1,5,#000000,0,0,0,#000000,1,40,1920 ;Episode names",
		"preset_define 31,6,36,#FFFFFF,0,0,0,1,3,#000000,0,0,0,#000000,-2,-1 ;Save dates",
		"preset_define 32,1,40,#FFFFFF,0,0,0,1,8,#000000,0,0,0,#000000,2,20,1590 ;Menu murder'",
		"new_episode 1",
		// Red truth - should use class, not colour
		`d [lv 0*"27"*"10100001"]"{p:1:Red truth test line here}."[\]`,
		// Blue truth - should use class, not colour
		`d [lv 0*"10"*"10100002"]"{p:2:Blue truth test line here}."[\]`,
		// Gold - should use inline colour
		`d [lv 0*"10"*"10100003"]"{p:41:Gold truth test line here}."[\]`,
		// Purple - should use inline colour
		`d [lv 0*"17"*"10100004"]"{p:42:Purple truth test line}."[\]`,
		// White preset - should NOT add colour span
		`d [lv 0*"10"*"10100005"]"Text with {p:0:Japanese font preset} here."[\]`,
	}

	quotes := p.ParseAll(lines)
	if len(quotes) != 5 {
		t.Fatalf("expected 5 quotes, got %d", len(quotes))
	}

	// Red truth uses semantic class
	if !strings.Contains(quotes[0].TextHtml, `class="red-truth"`) {
		t.Errorf("red truth should use class: %q", quotes[0].TextHtml)
	}
	if strings.Contains(quotes[0].TextHtml, "color:#FF0000") {
		t.Errorf("red truth should NOT use inline colour: %q", quotes[0].TextHtml)
	}

	// Blue truth uses semantic class
	if !strings.Contains(quotes[1].TextHtml, `class="blue-truth"`) {
		t.Errorf("blue truth should use class: %q", quotes[1].TextHtml)
	}

	// Gold uses inline colour from parsed preset
	if !strings.Contains(quotes[2].TextHtml, `color:#FFAA00`) {
		t.Errorf("gold should use parsed colour: %q", quotes[2].TextHtml)
	}

	// Purple uses inline colour from parsed preset
	if !strings.Contains(quotes[3].TextHtml, `color:#AA71FF`) {
		t.Errorf("purple should use parsed colour: %q", quotes[3].TextHtml)
	}

	// White preset (0) should NOT add any colour span
	if strings.Contains(quotes[4].TextHtml, "color:") {
		t.Errorf("white preset should not add colour: %q", quotes[4].TextHtml)
	}
	if !strings.Contains(quotes[4].Text, "Japanese font preset") {
		t.Errorf("content should be preserved: %q", quotes[4].Text)
	}
}

func TestParseAll_EpisodeAndContentTypes(t *testing.T) {
	p := testParser

	type testCase struct {
		name            string
		enLines         []string
		jaLines         []string
		wantEpisode     int
		wantContentType string
		wantCharID      string
		wantAudioID     string
		enTextContains  []string
		jaTextContains  []string
	}

	tests := []testCase{
		{
			name: "episode 1 regular - Nanjo",
			enLines: []string{
				"new_episode 1",
				"flush 10,167",
				"msgwnd_nan",
				"d [lv 0*\"19\"*\"11900001\"]`\"............Again. `[@][lv 0*\"19\"*\"11900002\"]`...So, you still haven't overcome your love of alcohol?\" `[\\]",
			},
			jaLines: []string{
				"new_episode 1",
				"flush 10,167",
				"msgwnd_nan",
				"d [lv 0*\"19\"*\"11900001\"]`\"\u2026\u2026\u2026\u2026\u2026\u307e\u305f\u3002`[@][lv 0*\"19\"*\"11900002\"]`\u2026\u304a\u9152\u3092\u5603\u307e\u308c\u307e\u3057\u305f\u306a\uff1f\u300d`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "",
			wantCharID:      "19",
			wantAudioID:     "11900001, 11900002",
			enTextContains:  []string{"Again", "alcohol"},
			jaTextContains:  []string{"また", "お酒"},
		},
		{
			name: "episode 5 regular - Beatrice d2 line",
			enLines: []string{
				"new_episode 5",
				"msgwnd_bea",
				"d2 [lv 0*\"27\"*\"50700001\"]`\"You idiot!! `[@][lv 0*\"27\"*\"50700002\"]`Isn't it obvious?!! `[@][#][*][lv 0*\"27\"*\"50700003\"]`It'll be fun to kill her and see your face twist in pain, why else?!! `[@][lv 0*\"27\"*\"50700004\"]`Gahyahahahahahahaha!!\" `[\\]",
			},
			jaLines: []string{
				"new_episode 5",
				"msgwnd_bea",
				"d2 [lv 0*\"27\"*\"50700001\"]`\"\u3053\u306e\u99ac\u9e7f\u304c\u30c3\uff01\uff01`[@][lv 0*\"27\"*\"50700002\"]`\u3000\u6c7a\u307e\u3063\u3066\u304a\u308d\u3046\u304c\u30c3\uff01\uff01`[@][#][*][lv 0*\"27\"*\"50700003\"]`\u3000\u6bba\u3057\u305f\u3089\u305d\u306a\u305f\u304c\u6b6a\u3081\u308b\u3060\u308d\u3046\u305d\u306e\u8868\u60c5\u304c\u697d\u3057\u3044\u304b\u3089\u306e\u4ed6\u306b\u4f55\u306e\u7406\u7531\u304c\u5fc5\u8981\u306a\u306e\u304b\u30c3\uff01\uff01`[@][lv 0*\"27\"*\"50700004\"]`\u3000\u304f\u3063\u3072\u3083\u306f\u306f\u306f\u306f\u306f\u306f\u306f\u306f\u306f\uff01\uff01\u300d`[\\]",
			},
			wantEpisode:     5,
			wantContentType: "",
			wantCharID:      "27",
			wantAudioID:     "50700001, 50700002, 50700003, 50700004",
			enTextContains:  []string{"You idiot", "kill her"},
			jaTextContains:  []string{"この馬鹿", "殺したら"},
		},
		{
			name: "episode 8 regular - Battler",
			enLines: []string{
				"new_episode 8",
				"msgwnd_but",
				"d [lv 0*\"10\"*\"80100001\"]`\"Are you listening, Ange? `[@][lv 0*\"10\"*\"80100002\"]`Take good care of this.\" `[\\]",
			},
			jaLines: []string{
				"new_episode 8",
				"msgwnd_but",
				"d [lv 0*\"10\"*\"80100001\"]`\"\u805e\u3044\u3066\u308b\u304b\u3001\u7e01\u5bff\u3002`[@][lv 0*\"10\"*\"80100002\"]`\u3000\u3053\u308c\u3092\u5927\u4e8b\u306b\u3057\u308d\u300d`[\\]",
			},
			wantEpisode:     8,
			wantContentType: "",
			wantCharID:      "10",
			wantAudioID:     "80100001, 80100002",
			enTextContains:  []string{"listening", "Ange"},
			jaTextContains:  []string{"縁寿"},
		},
		{
			name: "tea party 1 - Battler",
			enLines: []string{
				"new_tea 1",
				"msgwnd_but",
				"d [lv 0*\"10\"*\"90100001\"]`\"Hey, everyone, good job finishing \"Umineko no Naku Koro ni\"! `[@][lv 0*\"10\"*\"90100002\"]`Man, I still didn't have a clue what was going on when the story ended!\" `[\\]",
			},
			jaLines: []string{
				"new_tea 1",
				"msgwnd_but",
				"d [lv 0*\"10\"*\"90100001\"]`\"\u304a\u30fc\u3001\u307f\u3093\u306a\u300e\u3046\u307f\u306d\u3053\u306e\u306a\u304f\u9803\u306b\u300f\u3001\u304a\u75b2\u308c\u3055\u3093\uff01`[@][lv 0*\"10\"*\"90100002\"]`\u3000\u3084\u308c\u3084\u308c\u3001\u308f\u3051\u304c\u308f\u304b\u3093\u306a\u3044\u5185\u306b\u7269\u8a9e\u304c\u7d42\u308f\u3063\u3061\u307e\u3063\u305f\u306a\u3041\uff01\u300d`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "tea",
			wantCharID:      "10",
			wantAudioID:     "90100001, 90100002",
			enTextContains:  []string{"Hey, everyone", "Umineko"},
			jaTextContains:  []string{"おー、みんな", "うみねこ"},
		},
		{
			name: "tea party 1 - Maria",
			enLines: []string{
				"new_tea 1",
				"d [lv 0*\"13\"*\"90400001\"]`\"Uu-. `[@][lv 0*\"13\"*\"90400002\"]`Definitely a bad ending. `[@][lv 0*\"13\"*\"90400003\"]`Uu-.\" `[\\]",
			},
			jaLines: []string{
				"new_tea 1",
				"d [lv 0*\"13\"*\"90400001\"]`\"\u3046\u30fc\u3002`[@][lv 0*\"13\"*\"90400002\"]`\u304d\u3063\u3068\u30d0\u30c3\u30c9\u30a8\u30f3\u30c9\u3002`[@][lv 0*\"13\"*\"90400003\"]`\u3046\u30fc\u300d`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "tea",
			wantCharID:      "13",
			wantAudioID:     "90400001, 90400002, 90400003",
			enTextContains:  []string{"Uu-", "bad ending"},
			jaTextContains:  []string{"うー", "バッドエンド"},
		},
		{
			name: "tea party 5 - Bernkastel",
			enLines: []string{
				"new_tea 5",
				"d [lv 0*\"28\"*\"52100771\"]`\"...Erika, Dlanor.\" `[\\]",
			},
			jaLines: []string{
				"new_tea 5",
				"d [lv 0*\"28\"*\"52100771\"]`\"\u2026\u2026\u30f1\u30ea\u30ab\u3001\u30c9\u30e9\u30ce\u30fc\u30eb\u300d`[\\]",
			},
			wantEpisode:     5,
			wantContentType: "tea",
			wantCharID:      "28",
			wantAudioID:     "52100771",
			enTextContains:  []string{"Erika", "Dlanor"},
			jaTextContains:  []string{"ヱリカ", "ドラノール"},
		},
		{
			name: "ura 1 - Beatrice",
			enLines: []string{
				"new_ura 1",
				"msgwnd_bea",
				"d [lv 0*\"27\"*\"90700088\"]`\"...What sort of tea should I prepare next? `[@][lv 0*\"27\"*\"90700089\"]`You have free choice of any brand that ever was throughout the ages.\" `[\\]",
			},
			jaLines: []string{
				"new_ura 1",
				"msgwnd_bea",
				"d [lv 0*\"27\"*\"90700088\"]`\"\u2026\u6b21\u306f\u4f55\u306e\u7d05\u8336\u3092\u6df9\u308c\u3088\u3046\u304b\uff1f`[@][lv 0*\"27\"*\"90700089\"]`\u3000\u53e4\u4eca\u306e\u3042\u308a\u3068\u3042\u3089\u3086\u308b\u9298\u8336\u3092\u62ab\u9732\u3057\u3088\u3046\u305e\u300d`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "ura",
			wantCharID:      "27",
			wantAudioID:     "90700088, 90700089",
			enTextContains:  []string{"tea", "prepare"},
			jaTextContains:  []string{"紅茶", "淹れよう"},
		},
		{
			name: "ura 1 - Bernkastel",
			enLines: []string{
				"new_ura 1",
				"d [lv 0*\"28\"*\"92100001\"]`\"...Dried plum black tea. `[@][lv 0*\"28\"*\"92100002\"]`...The kind that goes for 200 yen a pack.\" `[\\]",
			},
			jaLines: []string{
				"new_ura 1",
				"d [lv 0*\"28\"*\"92100001\"]`\"\u2026\u2026\u6885\u5e72\u7d05\u8336\u3002`[@][lv 0*\"28\"*\"92100002\"]`\u2026\u2026\u6885\u5e72\u306f\uff11\u30d1\u30c3\u30af\uff12\uff10\uff10\u5186\u306e\u30e4\u30c4\u3088\u300d`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "ura",
			wantCharID:      "28",
			wantAudioID:     "92100001, 92100002",
			enTextContains:  []string{"Dried plum", "200 yen"},
			jaTextContains:  []string{"梅干紅茶", "２００円"},
		},
		{
			name: "episode 6 regular - Featherine",
			enLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000001\"]`\"Splendid. `[@][lv 0*\"50\"*\"65000002\"]`You did well to see through my veil...\" `[\\]",
			},
			jaLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000001\"]`\"\u898b\u4e8b\u306a\u308a\u3002`[@][lv 0*\"50\"*\"65000002\"]`\u3088\u304f\u305e\u3001\u79c1\u3092\u898b\u7834\u3063\u305f\u2026\u300d`[\\]",
			},
			wantEpisode:     6,
			wantContentType: "",
			wantCharID:      "50",
			wantAudioID:     "65000001, 65000002",
			enTextContains:  []string{"Splendid", "veil"},
			jaTextContains:  []string{"見事", "見破った"},
		},
		{
			name: "episode 6 regular - Featherine with font name tag",
			enLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000011\"]`\"Be an observer for me. An observer of the Fragments {f:5:Beatrice} has woven.\" `[\\]",
			},
			jaLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000011\"]`\"\u79c1\u306e\u305f\u3081\u306b\u3001\u30d9\u30a2\u30c8\u30ea\u30fc\u30c1\u30a7\u306e\u7d21\u3050\u30ab\u30b1\u30e9\u306e\u89b3\u6e2c\u8005\u3067\u3042\u308c\u300d`[\\]",
			},
			wantEpisode:     6,
			wantContentType: "",
			wantCharID:      "50",
			wantAudioID:     "65000011",
			enTextContains:  []string{"observer", "Fragments", "Beatrice"},
			jaTextContains:  []string{"観測者", "ベアトリーチェ"},
		},
		{
			name: "episode 6 regular - Featherine multi-segment",
			enLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000003\"]`\"I find you truly intriguing, child of man. `[@][lv 0*\"50\"*\"65000004\"]`...Your charming nature is the perfect medicine for my boredom.\" `[\\]",
			},
			jaLines: []string{
				"new_episode 6",
				"d [lv 0*\"50\"*\"65000003\"]`\"\u9762\u767d\u304d\u304b\u306a\u3001\u4eba\u306e\u5b50\u3088\u3002`[@][lv 0*\"50\"*\"65000004\"]`\u2026\u2026\u6109\u5feb\u306a\u308a\u3001\u305d\u308c\u3067\u3053\u305d\u6211\u304c\u9000\u5c48\u306b\u76f8\u5fdc\u3057\u3044\u300d`[\\]",
			},
			wantEpisode:     6,
			wantContentType: "",
			wantCharID:      "50",
			wantAudioID:     "65000003, 65000004",
			enTextContains:  []string{"intriguing", "child of man", "boredom"},
			jaTextContains:  []string{"面白きかな", "人の子", "退屈"},
		},
		{
			name: "omake 1 - Jessica (English only)",
			enLines: []string{
				"*o1_0",
				"bgmplay 89,71,0",
				"*o1_1",
				"d [lv 0*\"04\"*\"10200442\"]`\"KyaaaaaAAAAAAaaaaaAAaa!!!\"`[\\]",
			},
			wantEpisode:     1,
			wantContentType: "omake",
			wantCharID:      "04",
			wantAudioID:     "10200442",
			enTextContains:  []string{"KyaaaaaAAAAAAaaaaaAAaa"},
		},
		{
			name: "alphanumeric audio ID - awase group voice",
			enLines: []string{
				"new_episode 6",
				"d2 [lv 0*\"00\"*\"awase6100_o\"][ak][text_speed_t 5]`\"\"With your fellow monsters.\"\" `[#][*][\\]",
			},
			jaLines: []string{
				"new_episode 6",
				"d2 [lv 0*\"00\"*\"awase6100_o\"][ak][text_speed_t 5]`\"\u300c\u30d0\u30b1\u30e2\u30ce\u540c\u58eb\u306b\u9650\u308b\u308f\u300d\u300d`[#][*][\\]",
			},
			wantEpisode:     6,
			wantContentType: "",
			wantCharID:      "00",
			wantAudioID:     "awase6100_o",
			enTextContains:  []string{"With your fellow monsters"},
			jaTextContains:  []string{"バケモノ同士"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enQuotes := p.ParseAll(tt.enLines)
			if len(enQuotes) == 0 {
				t.Fatal("EN: expected at least 1 quote, got 0")
			}
			en := enQuotes[0]

			if en.Episode != tt.wantEpisode {
				t.Errorf("EN episode: got %d, want %d", en.Episode, tt.wantEpisode)
			}
			if en.ContentType != tt.wantContentType {
				t.Errorf("EN contentType: got %q, want %q", en.ContentType, tt.wantContentType)
			}
			if en.CharacterID != tt.wantCharID {
				t.Errorf("EN characterID: got %q, want %q", en.CharacterID, tt.wantCharID)
			}
			if en.Character != CharacterNames[tt.wantCharID] {
				t.Errorf("EN character: got %q, want %q", en.Character, CharacterNames[tt.wantCharID])
			}
			if en.AudioID != tt.wantAudioID {
				t.Errorf("EN audioID: got %q, want %q", en.AudioID, tt.wantAudioID)
			}
			for _, sub := range tt.enTextContains {
				if !strings.Contains(en.Text, sub) {
					t.Errorf("EN text %q missing substring %q", en.Text, sub)
				}
			}

			if tt.jaLines != nil {
				jaQuotes := p.ParseAll(tt.jaLines)
				if len(jaQuotes) == 0 {
					t.Fatal("JA: expected at least 1 quote, got 0")
				}
				ja := jaQuotes[0]

				if ja.Episode != tt.wantEpisode {
					t.Errorf("JA episode: got %d, want %d", ja.Episode, tt.wantEpisode)
				}
				if ja.ContentType != tt.wantContentType {
					t.Errorf("JA contentType: got %q, want %q", ja.ContentType, tt.wantContentType)
				}
				if ja.CharacterID != tt.wantCharID {
					t.Errorf("JA characterID: got %q, want %q", ja.CharacterID, tt.wantCharID)
				}
				if ja.Character != CharacterNames[tt.wantCharID] {
					t.Errorf("JA character: got %q, want %q", ja.Character, CharacterNames[tt.wantCharID])
				}
				if ja.AudioID != tt.wantAudioID {
					t.Errorf("JA audioID: got %q, want %q", ja.AudioID, tt.wantAudioID)
				}
				for _, sub := range tt.jaTextContains {
					if !strings.Contains(ja.Text, sub) {
						t.Errorf("JA text %q missing substring %q", ja.Text, sub)
					}
				}
			}
		})
	}
}

func TestParseAll_ContentTypeTransitions(t *testing.T) {
	p := testParser

	// Simulate a sequence: episode → tea → ura → new episode (resets content type)
	lines := []string{
		"new_episode 1",
		"d [lv 0*\"19\"*\"11900001\"]`\"............Again. `[@][lv 0*\"19\"*\"11900002\"]`...So, you still haven't overcome your love of alcohol?\" `[\\]",
		"new_tea 1",
		"d [lv 0*\"10\"*\"90100001\"]`\"Hey, everyone, good job finishing the story! `[@][lv 0*\"10\"*\"90100002\"]`Man, I still didn't have a clue what was going on!\" `[\\]",
		"new_ura 1",
		"d [lv 0*\"27\"*\"90700088\"]`\"...What sort of tea should I prepare next? `[@][lv 0*\"27\"*\"90700089\"]`You have free choice of any brand.\" `[\\]",
		"new_episode 2",
		"d [lv 0*\"10\"*\"20100001\"]`\"The second episode begins here with a long enough line.\" `[\\]",
	}

	quotes := p.ParseAll(lines)
	if len(quotes) != 4 {
		t.Fatalf("expected 4 quotes, got %d", len(quotes))
	}

	// Episode 1 regular
	if quotes[0].Episode != 1 || quotes[0].ContentType != "" {
		t.Errorf("quote 0: got ep=%d ct=%q, want ep=1 ct=\"\"", quotes[0].Episode, quotes[0].ContentType)
	}

	// Tea party 1
	if quotes[1].Episode != 1 || quotes[1].ContentType != "tea" {
		t.Errorf("quote 1: got ep=%d ct=%q, want ep=1 ct=\"tea\"", quotes[1].Episode, quotes[1].ContentType)
	}

	// Ura 1
	if quotes[2].Episode != 1 || quotes[2].ContentType != "ura" {
		t.Errorf("quote 2: got ep=%d ct=%q, want ep=1 ct=\"ura\"", quotes[2].Episode, quotes[2].ContentType)
	}

	// Episode 2 regular — content type must be reset
	if quotes[3].Episode != 2 || quotes[3].ContentType != "" {
		t.Errorf("quote 3: got ep=%d ct=%q, want ep=2 ct=\"\"", quotes[3].Episode, quotes[3].ContentType)
	}
}

func TestParseAll_OmakeTransition(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 8",
		"d [lv 0*\"10\"*\"80100001\"]`\"Episode 8 content with sufficient length for the test.\" `[\\]",
		"*o1_0",
		"*o1_1",
		"d [lv 0*\"04\"*\"10200442\"]`\"KyaaaaaAAAAAAaaaaaAAaa!!!\"`[\\]",
		"*o2_0",
		"*o2_1",
		"d [lv 0*\"10\"*\"20100001\"]`\"Omake episode 2 content with sufficient length here.\" `[\\]",
	}

	quotes := p.ParseAll(lines)
	if len(quotes) != 3 {
		t.Fatalf("expected 3 quotes, got %d", len(quotes))
	}

	if quotes[0].Episode != 8 || quotes[0].ContentType != "" {
		t.Errorf("quote 0: got ep=%d ct=%q, want ep=8 ct=\"\"", quotes[0].Episode, quotes[0].ContentType)
	}
	if quotes[1].Episode != 1 || quotes[1].ContentType != "omake" {
		t.Errorf("quote 1: got ep=%d ct=%q, want ep=1 ct=\"omake\"", quotes[1].Episode, quotes[1].ContentType)
	}
	if quotes[2].Episode != 2 || quotes[2].ContentType != "omake" {
		t.Errorf("quote 2: got ep=%d ct=%q, want ep=2 ct=\"omake\"", quotes[2].Episode, quotes[2].ContentType)
	}
}

func TestParseAll_RedTruth(t *testing.T) {
	p := testParser

	t.Run("english red truth - Beatrice", func(t *testing.T) {
		lines := []string{
			"new_episode 2",
			"d2 [lv 0*\"27\"*\"20700951\"]`\" `[#][*]`{p:1:Everything I speak in red is the truth}! `[@][lv 0*\"27\"*\"20700952\"]`There's absolutely no need to doubt it!\" `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]
		if q.Episode != 2 {
			t.Errorf("episode: got %d, want 2", q.Episode)
		}
		if q.CharacterID != "27" || q.Character != CharacterNames["27"] {
			t.Errorf("character: got %q/%q, want 27/%s", q.CharacterID, q.Character, CharacterNames["27"])
		}
		if !strings.Contains(q.Text, "Everything I speak in red is the truth") {
			t.Errorf("plain text missing red truth content: %q", q.Text)
		}
		if !strings.Contains(q.TextHtml, `<span class="red-truth">Everything I speak in red is the truth</span>`) {
			t.Errorf("HTML missing red truth span: %q", q.TextHtml)
		}
		// Plain text should NOT contain the span tag
		if strings.Contains(q.Text, "red-truth") {
			t.Errorf("plain text should not contain HTML: %q", q.Text)
		}
	})

	t.Run("japanese red truth - Beatrice", func(t *testing.T) {
		lines := []string{
			"new_episode 2",
			"d2 [lv 0*\"27\"*\"20700951\"]`\"\u300c`[#][*]`{p:1:\u59be\u304c\u8d64\u3067\u8a9e\u308b\u3053\u3068\u306f\u5168\u3066\u771f\u5b9f}\uff01`[@][lv 0*\"27\"*\"20700952\"]`\u3000\u7591\u3046\u5fc5\u8981\u304c\u4f55\u3082\u306a\u3044\uff01\u300d`[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]
		if q.Episode != 2 {
			t.Errorf("episode: got %d, want 2", q.Episode)
		}
		if q.CharacterID != "27" {
			t.Errorf("characterID: got %q, want 27", q.CharacterID)
		}
		if !strings.Contains(q.Text, "妾が赤で語ることは全て真実") {
			t.Errorf("JA plain text missing red truth content: %q", q.Text)
		}
		if !strings.Contains(q.TextHtml, `<span class="red-truth">`) {
			t.Errorf("JA HTML missing red truth span: %q", q.TextHtml)
		}
	})
}

func TestParseAll_BlueTruth(t *testing.T) {
	p := testParser

	t.Run("english blue truth - Battler", func(t *testing.T) {
		lines := []string{
			"new_episode 4",
			"d2 [lv 0*\"10\"*\"40100220\"]`\"This is my truth! `[@][lv 0*\"10\"*\"40100221\"]`{p:2:Ushiromiya Kinzo is already dead}. `[#][*][@][lv 0*\"10\"*\"40100222\"]`{p:2:Therefore, the true number of people on the island is 17}! `[#][*][@][lv 0*\"10\"*\"40100223\"]`{p:2:By adding an unknown person X to that, it becomes 18 people}. `[#][*][@][lv 0*\"10\"*\"40100224\"]`{p:2:By supposing that this person X exists, the crime is possible even if all 17 people have alibis}!!\" `[#][*][\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]
		if q.Episode != 4 {
			t.Errorf("episode: got %d, want 4", q.Episode)
		}
		if q.CharacterID != "10" || q.Character != CharacterNames["10"] {
			t.Errorf("character: got %q/%q, want 10/%s", q.CharacterID, q.Character, CharacterNames["10"])
		}
		if !strings.Contains(q.Text, "Ushiromiya Kinzo is already dead") {
			t.Errorf("plain text missing blue truth content: %q", q.Text)
		}
		if !strings.Contains(q.TextHtml, `<span class="blue-truth">Ushiromiya Kinzo is already dead</span>`) {
			t.Errorf("HTML missing blue truth span: %q", q.TextHtml)
		}
		if !strings.Contains(q.TextHtml, `<span class="blue-truth">Therefore, the true number of people on the island is 17</span>`) {
			t.Errorf("HTML missing second blue truth span: %q", q.TextHtml)
		}
		// Verify audio IDs for all 5 segments
		if q.AudioID != "40100220, 40100221, 40100222, 40100223, 40100224" {
			t.Errorf("audioID: got %q, want all 5 IDs", q.AudioID)
		}
	})

	t.Run("japanese blue truth - Battler", func(t *testing.T) {
		lines := []string{
			"new_episode 4",
			"d2 [lv 0*\"10\"*\"40100220\"]`\"\u3053\u308c\u304c\u4fc5\u306e\u771f\u5b9f\u3060\uff01`[@][lv 0*\"10\"*\"40100221\"]`\u3000{p:2:\u53f3\u4ee3\u5bae\u91d1\u8535\u306f\u3059\u3067\u306b\u6b7b\u4ea1\u3057\u3066\u3044\u308b}\u3002`[#][*][@][lv 0*\"10\"*\"40100222\"]`{p:2:\u3088\u3063\u3066\u5cf6\u306e\u672c\u5f53\u306e\u4eba\u6570\u306f\uff11\uff17\u4eba}\uff01`[#][*][@][lv 0*\"10\"*\"40100223\"]`\u3000{p:2:\u305d\u3053\u306b\u672a\u77e5\u306e\u4eba\u7269\uff38\u304c\u52a0\u308f\u308b\u3053\u3068\u3067\uff11\uff18\u4eba\u3068\u306a\u3063\u3066\u3044\u308b}\u3002`[#][*][@][lv 0*\"10\"*\"40100224\"]`{p:2:\u3053\u306e\u4eba\u7269\uff38\u306e\u5b58\u5728\u306e\u4eee\u5b9a\u306b\u3088\u3063\u3066\u3001\uff11\uff17\u4eba\u5168\u54e1\u306b\u30a2\u30ea\u30d0\u30a4\u304c\u3042\u3063\u3066\u3082\u72af\u884c\u306f\u53ef\u80fd\u306b\u306a\u308b}\u30c3\uff01\uff01\u300d`[#][*][\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]
		if q.Episode != 4 {
			t.Errorf("episode: got %d, want 4", q.Episode)
		}
		if q.CharacterID != "10" {
			t.Errorf("characterID: got %q, want 10", q.CharacterID)
		}
		if !strings.Contains(q.Text, "右代宮金蔵はすでに死亡している") {
			t.Errorf("JA plain text missing blue truth content: %q", q.Text)
		}
		if !strings.Contains(q.TextHtml, `<span class="blue-truth">`) {
			t.Errorf("JA HTML missing blue truth span: %q", q.TextHtml)
		}
	})
}

func TestParseAll_ColourFormatting(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 1",
		"d [lv 0*\"19\"*\"11900003\"]`\"............Kinzo\u2010{c:86EF9C:san}. `[gstg 1][@][lv 0*\"19\"*\"11900004\"]`...your body only appears to be well thanks to the effects of the medicine.\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	// Plain text: {c:} tag stripped, content preserved
	if !strings.Contains(q.Text, "san") {
		t.Errorf("plain text missing colour content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{c:") {
		t.Errorf("plain text should not contain raw {c:} tag: %q", q.Text)
	}

	// HTML: colour span generated
	if !strings.Contains(q.TextHtml, `<span style="color:#86EF9C">san</span>`) {
		t.Errorf("HTML missing colour span: %q", q.TextHtml)
	}
}

func TestParseAll_RubyAnnotations(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 1",
		"d `You can't really tell Grandfather's story without covering that pivotal event back before the {ruby:1926-1989:Showa} era. `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	// Plain text: ruby rendered as "text (annotation)"
	if !strings.Contains(q.Text, "Showa (1926-1989)") {
		t.Errorf("plain text missing ruby content: %q", q.Text)
	}

	// HTML: proper ruby markup
	if !strings.Contains(q.TextHtml, "<ruby>Showa<rp>(</rp><rt>1926-1989</rt><rp>)</rp></ruby>") {
		t.Errorf("HTML missing ruby markup: %q", q.TextHtml)
	}
}

func TestParseAll_FontNameFormatting(t *testing.T) {
	p := testParser

	// Ura 1 line with {f:5:Bernkastel}
	lines := []string{
		"new_ura 1",
		"d2 [lv 0*\"27\"*\"90700093\"]`\"...*cackle*cackle* Hostility? I bear nothing of the sort. `[@][#][*][lv 0*\"27\"*\"90700094\"]`...I am just terribly concerned with seeing to it that the legendary witch, Lady {f:5:Bernkastel}, is well attended to.\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.Text, "Bernkastel") {
		t.Errorf("plain text missing font-name content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{f:") {
		t.Errorf("plain text should not contain raw {f:} tag: %q", q.Text)
	}
	// {f:} is a font selection tag (small caps), not a name highlight — should pass through as plain text
	if strings.Contains(q.TextHtml, "quote-name") {
		t.Errorf("HTML should not contain quote-name span for {f:} font tag: %q", q.TextHtml)
	}
	if !strings.Contains(q.TextHtml, "Bernkastel") {
		t.Errorf("HTML missing font content: %q", q.TextHtml)
	}
}

func TestParseAll_FontSelectionPassthrough(t *testing.T) {
	p := testParser

	// Real line: Kinzo calling out to Beatrice with {f:5:} small caps (EP1)
	lines := []string{
		"new_episode 1",
		"d [lv 0*\"01\"*\"11500055\"]`And even if I had a chance to write such a thing, `[|][lv 0*\"01\"*\"11500056\"]`......if I did have such a chance!!\" `[@][lv 0*\"01\"*\"11500057\"]`\"...I'd want to see it. `[@][lv 0*\"01\"*\"11500058\"]`I want to see it one more time! `[@][lv 0*\"01\"*\"11500059\"]`I want to see {f:5:Beatrice}'s smiling face one last time!! `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.Text, "Beatrice") {
		t.Errorf("plain text should contain inner text: %q", q.Text)
	}
	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text should not contain braces: %q", q.Text)
	}
	if strings.Contains(q.TextHtml, "quote-name") {
		t.Errorf("HTML should not contain quote-name for font tag: %q", q.TextHtml)
	}
	if strings.Contains(q.TextHtml, "{") || strings.Contains(q.TextHtml, "}") {
		t.Errorf("HTML should not contain braces: %q", q.TextHtml)
	}
}

func TestParseAll_PresetGoldText(t *testing.T) {
	p := testParser

	// Real line: Battler's gold truth declaration (EP5)
	lines := []string{
		"preset_define 41,1,-1,#FFAA00,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Gold text",
		"new_episode 5",
		"d [lv 0*\"10\"*\"50101465\"]`\"{p:41:I guarantee that this is the corpse of Ushiromiya Kinzo}...!!\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.TextHtml, `<span style="color:#FFAA00">I guarantee that this is the corpse of Ushiromiya Kinzo</span>`) {
		t.Errorf("HTML should contain gold colour span: %q", q.TextHtml)
	}
	if !strings.Contains(q.Text, "I guarantee that this is the corpse of Ushiromiya Kinzo") {
		t.Errorf("plain text should contain gold text content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text should not contain braces: %q", q.Text)
	}
}

func TestParseAll_PresetPurpleText(t *testing.T) {
	p := testParser

	// Real line: Kumasawa's purple declaration in Bernkastel's game (EP8)
	lines := []string{
		"preset_define 42,1,-1,#AA71FF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Purple text",
		"new_episode 8",
		"d [lv 0*\"17\"*\"81700082\"]`\"{p:42:The dining hall was locked up}. `[@][lv 0*\"17\"*\"81700083\"]`Of course, there is a lock on the doors to the dining hall, but they usually remained unlocked. `[@][lv 0*\"17\"*\"81700084\"]`I knocked, but there was no answer...\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.TextHtml, `<span style="color:#AA71FF">The dining hall was locked up</span>`) {
		t.Errorf("HTML should contain purple colour span: %q", q.TextHtml)
	}
	if !strings.Contains(q.Text, "The dining hall was locked up") {
		t.Errorf("plain text should contain purple text content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text should not contain braces: %q", q.Text)
	}
}

func TestParseAll_PresetPassthrough(t *testing.T) {
	p := testParser

	// Real line: Battler introduces his name with {p:0:} Japanese font preset (EP1)
	lines := []string{
		"preset_define 0,6,36,#FFFFFF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1 ;Japanese text (white)",
		"new_episode 1",
		"d `My name is written {p:0:\u53f3\u4ee3\u5bae \u6226\u4eba}. `[@]`Can you read it? `[@]`The first part is my family name, \"Ushiromiya\". `[@]`That's a fairly plausible Japanese pronunciation so far. `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.Text, "\u53f3\u4ee3\u5bae \u6226\u4eba") {
		t.Errorf("plain text should contain Japanese content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text should not contain braces: %q", q.Text)
	}
	if strings.Contains(q.TextHtml, "quote-name") {
		t.Errorf("HTML should not contain quote-name for generic preset: %q", q.TextHtml)
	}
	if strings.Contains(q.TextHtml, "color:") {
		t.Errorf("HTML should not contain colour for generic preset: %q", q.TextHtml)
	}
}

func TestParseAll_StrayBracesStripped(t *testing.T) {
	p := testParser

	// Dlanor's line: {p:1: red truth spans across backtick segments, leaving a stray }
	lines := []string{
		"new_episode 5",
		`d2 [lv 0*"47"*"54600338"][#][*]` + "`" + `"{p:1:Die the death!  ` + "`" + `[|][lv 0*"47"*"54600339"][#][*][!w167]` + "`" + `Sentence to death!!{n}` + "`" + `[|][lv 0*"47"*"54600340"][#][*][!w167]` + "`" + `Great equalizer is the death!!}" ` + "`" + `[\]`,
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text should not contain braces: %q", q.Text)
	}
	if strings.Contains(q.TextHtml, "{") || strings.Contains(q.TextHtml, "}") {
		t.Errorf("HTML should not contain braces: %q", q.TextHtml)
	}
	if !strings.Contains(q.Text, "Die the death") {
		t.Errorf("plain text should contain the quote content: %q", q.Text)
	}
}

func TestParseAll_NoBracesInOutput(t *testing.T) {
	p := testParser

	data, err := dataFS.ReadFile("data/english.txt")
	if err != nil {
		t.Fatalf("failed to read english.txt: %v", err)
	}
	lines := strings.Split(string(data), "\n")
	quotes := p.ParseAll(lines)

	for i, q := range quotes {
		if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
			t.Errorf("quote %d text contains braces: %q", i, q.Text)
		}
		if strings.Contains(q.TextHtml, "{") || strings.Contains(q.TextHtml, "}") {
			t.Errorf("quote %d HTML contains braces: %q", i, q.TextHtml)
		}
	}
}

func TestParseAll_ItalicFormatting(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 1",
		"d `...Wait a sec, we {i:are} totally sure it's not gonna shake, ...right...? `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.Text, "are") {
		t.Errorf("plain text missing italic content: %q", q.Text)
	}
	if strings.Contains(q.Text, "{i:") {
		t.Errorf("plain text should not contain raw {i:} tag: %q", q.Text)
	}
	if !strings.Contains(q.TextHtml, "<em>are</em>") {
		t.Errorf("HTML missing em tag: %q", q.TextHtml)
	}
}

func TestParseAll_LineBreaks(t *testing.T) {
	p := testParser

	t.Run("english narrator with line break tag", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `In the corner of this room, which was larger than most, `[@]`{n}was an expensive-looking bed and a physician conducting an examination. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		// Plain text: {n} replaced with space
		if strings.Contains(q.Text, "{n}") {
			t.Errorf("plain text should not contain raw {n} tag: %q", q.Text)
		}
		// HTML: {n} replaced with <br>
		if !strings.Contains(q.TextHtml, "<br>") {
			t.Errorf("HTML missing <br> for line break: %q", q.TextHtml)
		}
	})

	t.Run("japanese narrator with line break", func(t *testing.T) {
		lines := []string{
			"new_episode 8",
			"d `\u793c\u62dd\u5802\u306e\u4e2d\u306b\u306f\u3001\u2026\u2026\u4e8c\u4eba\u306e\u4eba\u5f71\u3002`[@]`{n}\u305d\u308c\u306f\u3001\u5144\u3068\u3001\u59b9\u306e\u59ff\u3002`[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if !strings.Contains(q.Text, "礼拝堂") {
			t.Errorf("JA text missing expected content: %q", q.Text)
		}
		if strings.Contains(q.Text, "{n}") {
			t.Errorf("JA plain text should not contain raw {n} tag: %q", q.Text)
		}
		if !strings.Contains(q.TextHtml, "<br>") {
			t.Errorf("JA HTML missing <br> for line break: %q", q.TextHtml)
		}
	})
}

func TestParseAll_NarratorLines(t *testing.T) {
	p := testParser

	t.Run("english narrator basic", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `The old physician let out a sigh as he removed the stethoscope. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if q.Episode != 1 {
			t.Errorf("episode: got %d, want 1", q.Episode)
		}
		if q.CharacterID != "narrator" {
			t.Errorf("characterID: got %q, want narrator", q.CharacterID)
		}
		if q.Character != "Narrator" {
			t.Errorf("character: got %q, want Narrator", q.Character)
		}
		if q.AudioID != "" {
			t.Errorf("audioID: got %q, want empty", q.AudioID)
		}
		if !strings.Contains(q.Text, "physician") {
			t.Errorf("text missing expected content: %q", q.Text)
		}
	})

	t.Run("japanese narrator basic", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `\u8074\u8a3a\u5668\u3092\u5916\u3057\u306a\u304c\u3089\u3001\u5e74\u8f29\u306e\u533b\u5e2b\u306f\u6e9c\u3081\u606f\u3092\u6f0f\u3089\u3059\u3002`[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if q.Episode != 1 {
			t.Errorf("episode: got %d, want 1", q.Episode)
		}
		if q.CharacterID != "narrator" {
			t.Errorf("characterID: got %q, want narrator", q.CharacterID)
		}
		if q.Character != "Narrator" {
			t.Errorf("character: got %q, want Narrator", q.Character)
		}
		if !strings.Contains(q.Text, "聴診器") {
			t.Errorf("JA text missing expected content: %q", q.Text)
		}
	})

	t.Run("d2 narrator line", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d2 `After eyeing both the master who demanded the alcohol and the family doctor who forbade it, `[@][#][*]`Genji, the old butler, silently gave a slight nod and carried out his master's orders faithfully. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if q.CharacterID != "narrator" {
			t.Errorf("characterID: got %q, want narrator", q.CharacterID)
		}
		if !strings.Contains(q.Text, "Genji") {
			t.Errorf("text missing expected content: %q", q.Text)
		}
	})

	t.Run("d2 japanese narrator line", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d2 `\u6e90\u6b21\u3068\u547c\u3070\u308c\u305f\u8001\u9f62\u306e\u57f7\u4e8b\u306f\u3001\u9152\u3092\u6c42\u3081\u308b\u4e3b\u3068\u3001\u305d\u308c\u3092\u6b62\u3081\u308b\u4e3b\u6cbb\u533b\u306e\u53cc\u65b9\u3092\u898b\u6bd4\u3079\u305f\u5f8c\u3001`[@][#][*]`\u7121\u8a00\u3067\u5c0f\u3055\u304f\u9837\u304d\u3001\u5df1\u306e\u4e3b\u306e\u547d\u4ee4\u306b\u5fe0\u5b9f\u306b\u5f93\u3046\u306e\u3060\u3063\u305f\u3002`[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if q.CharacterID != "narrator" {
			t.Errorf("characterID: got %q, want narrator", q.CharacterID)
		}
		if !strings.Contains(q.Text, "源次") {
			t.Errorf("JA text missing expected content: %q", q.Text)
		}
	})

	t.Run("english episode 8 narrator with multi-segment", func(t *testing.T) {
		lines := []string{
			"new_episode 8",
			"d `Two silhouettes could be seen in the chapel. `[@]`A big brother...and a little sister. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		q := quotes[0]

		if q.Episode != 8 {
			t.Errorf("episode: got %d, want 8", q.Episode)
		}
		if !strings.Contains(q.Text, "chapel") {
			t.Errorf("text missing expected content: %q", q.Text)
		}
		if !strings.Contains(q.Text, "big brother") {
			t.Errorf("text missing expected content: %q", q.Text)
		}
	})
}

func TestParseAll_ShortLineFiltering(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 1",
		// This should be filtered (text <= 10 chars after processing)
		"d [lv 0*\"10\"*\"10100001\"]`\"Hello.\" `[\\]",
		// This should be kept (text > 10 chars)
		"d [lv 0*\"10\"*\"10100002\"]`\"This is a much longer line of dialogue that passes.\" `[\\]",
	}

	quotes := p.ParseAll(lines)
	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote (short line filtered), got %d", len(quotes))
	}
	if !strings.Contains(quotes[0].Text, "longer line") {
		t.Errorf("wrong quote kept: %q", quotes[0].Text)
	}
}

func TestParseAll_MultipleAudioIDs(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 5",
		"d2 [lv 0*\"27\"*\"50700001\"]`\"You idiot!! `[@][lv 0*\"27\"*\"50700002\"]`Isn't it obvious?!! `[@][#][*][lv 0*\"27\"*\"50700003\"]`It'll be fun to kill her and see your face twist in pain, why else?!! `[@][lv 0*\"27\"*\"50700004\"]`Gahyahahahahahahaha!!\" `[\\]",
	}

	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if q.AudioID != "50700001, 50700002, 50700003, 50700004" {
		t.Errorf("audioID: got %q, want \"50700001, 50700002, 50700003, 50700004\"", q.AudioID)
	}
}

func TestParseAll_DuplicateAudioIDsDeduped(t *testing.T) {
	p := testParser

	// Simulate a line with duplicate audio IDs
	lines := []string{
		"new_episode 1",
		"d [lv 0*\"10\"*\"10100001\"][lv 0*\"10\"*\"10100001\"]`\"This line has duplicate audio IDs that should be deduped.\" `[\\]",
	}

	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	if quotes[0].AudioID != "10100001" {
		t.Errorf("audioID should be deduped: got %q, want \"10100001\"", quotes[0].AudioID)
	}
}

func TestParseAll_SpecialCharTags(t *testing.T) {
	p := testParser

	t.Run("qt tag produces double quote", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `She said {qt}hello there{qt} in a cheerful and bright tone of voice. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		if !strings.Contains(quotes[0].Text, `"hello there"`) {
			t.Errorf("text should have quotes from {qt}: %q", quotes[0].Text)
		}
	})

	t.Run("ob and eb tags are stripped", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `The text contains {ob}curly braces{eb} that need preserving somehow. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		if strings.Contains(quotes[0].Text, "{") || strings.Contains(quotes[0].Text, "}") {
			t.Errorf("text should not contain braces: %q", quotes[0].Text)
		}
		if !strings.Contains(quotes[0].Text, "curly braces") {
			t.Errorf("text should preserve inner content: %q", quotes[0].Text)
		}
	})

	t.Run("os and es tags produce square brackets", func(t *testing.T) {
		lines := []string{
			"new_episode 1",
			"d `The text contains {os}square brackets{es} inside here I think. `[\\]",
		}
		quotes := p.ParseAll(lines)
		if len(quotes) == 0 {
			t.Fatal("expected at least 1 quote")
		}
		if !strings.Contains(quotes[0].Text, "[square brackets]") {
			t.Errorf("text should have square brackets from {os}/{es}: %q", quotes[0].Text)
		}
	})
}

func TestParseAll_EpisodeFromAudioID(t *testing.T) {
	p := testParser

	// Without a new_episode marker, the episode should be derived from the audio ID
	lines := []string{
		"d [lv 0*\"10\"*\"30100001\"]`\"This is a line with enough text for the parser filter.\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	// Audio ID "30100001" starts with '3', so episode = 3
	// But currentEpisode is 0, so the audio-derived episode (3) is used since
	// ParseAll only overrides when currentEpisode > 0
	if quotes[0].Episode != 3 {
		t.Errorf("episode from audioID: got %d, want 3", quotes[0].Episode)
	}
}

func TestParseAll_EpisodeMarkerOverridesAudioID(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 7",
		"d [lv 0*\"10\"*\"30100001\"]`\"This is a line with enough text for the parser filter.\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	// new_episode 7 overrides the audio-derived episode (3)
	if quotes[0].Episode != 7 {
		t.Errorf("episode should be from marker: got %d, want 7", quotes[0].Episode)
	}
}

func TestParseAll_NonDialogueNonNarratorLinesSkipped(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 1",
		"flush 10,167",
		"textoff",
		"waits 167",
		"lbg s0_3,\"black\"",
		"bgmplay 13,71,0",
		"msgwnd_nan",
		"*d1",
		"d [lv 0*\"19\"*\"11900001\"]`\"This is the only actual dialogue line with enough text.\" `[\\]",
		"lss s0_8,\"nan\",\"a1_fumu1\" ;1",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote from mixed lines, got %d", len(quotes))
	}
}

func TestParseAll_EmptyInput(t *testing.T) {
	p := testParser

	quotes := p.ParseAll(nil)
	if len(quotes) != 0 {
		t.Errorf("expected 0 quotes for nil input, got %d", len(quotes))
	}

	quotes = p.ParseAll([]string{})
	if len(quotes) != 0 {
		t.Errorf("expected 0 quotes for empty input, got %d", len(quotes))
	}
}

func TestParseAll_AllEpisodes(t *testing.T) {
	p := testParser

	// Verify all 8 episodes can be parsed
	var lines []string
	for ep := 1; ep <= 8; ep++ {
		lines = append(lines,
			"new_episode "+string(rune('0'+ep)),
			"d [lv 0*\"10\"*\""+string(rune('0'+ep))+"0100001\"]`\"Episode dialogue line that is long enough to pass the filter.\" `[\\]",
		)
	}

	quotes := p.ParseAll(lines)
	if len(quotes) != 8 {
		t.Fatalf("expected 8 quotes, got %d", len(quotes))
	}
	for i, q := range quotes {
		wantEp := i + 1
		if q.Episode != wantEp {
			t.Errorf("quote %d: episode got %d, want %d", i, q.Episode, wantEp)
		}
		if q.ContentType != "" {
			t.Errorf("quote %d: contentType got %q, want empty", i, q.ContentType)
		}
	}
}

func TestParseAll_AllContentTypesPerEpisode(t *testing.T) {
	p := testParser

	// For each episode 1-8, test all content types
	for ep := 1; ep <= 8; ep++ {
		epStr := string(rune('0' + ep))
		t.Run("episode_"+epStr, func(t *testing.T) {
			lines := []string{
				"new_episode " + epStr,
				"d [lv 0*\"10\"*\"" + epStr + "0100001\"]`\"Regular content for this episode that passes length filter.\" `[\\]",
				"new_tea " + epStr,
				"d [lv 0*\"10\"*\"9010000" + epStr + "\"]`\"Tea party content for this episode that passes length filter.\" `[\\]",
				"new_ura " + epStr,
				"d [lv 0*\"10\"*\"9210000" + epStr + "\"]`\"Ura content for this episode that passes the length filter.\" `[\\]",
			}
			quotes := p.ParseAll(lines)
			if len(quotes) != 3 {
				t.Fatalf("expected 3 quotes, got %d", len(quotes))
			}
			if quotes[0].Episode != ep || quotes[0].ContentType != "" {
				t.Errorf("regular: ep=%d ct=%q", quotes[0].Episode, quotes[0].ContentType)
			}
			if quotes[1].Episode != ep || quotes[1].ContentType != "tea" {
				t.Errorf("tea: ep=%d ct=%q", quotes[1].Episode, quotes[1].ContentType)
			}
			if quotes[2].Episode != ep || quotes[2].ContentType != "ura" {
				t.Errorf("ura: ep=%d ct=%q", quotes[2].Episode, quotes[2].ContentType)
			}
		})
	}
}

func TestParseAll_UnclosedTagCleanup(t *testing.T) {
	p := testParser

	// Unclosed {p:1: tag (no closing }) — from real data in episode 7
	lines := []string{
		"new_episode 7",
		"d2 [lv 0*\"46\"*\"54501382\"]`\"`[#][*]`{p:1:During that one hour, you were in the dining hall of the mansion}!! `[@][lv 0*\"46\"*\"54501383\"]`{p:1:This text has an unclosed tag that should be cleaned up`[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	// The closed {p:1:} should be rendered as red truth
	if !strings.Contains(q.TextHtml, `<span class="red-truth">During that one hour, you were in the dining hall of the mansion</span>`) {
		t.Errorf("HTML missing red truth for closed tag: %q", q.TextHtml)
	}
	// Unclosed tag should be cleaned up (removed)
	if strings.Contains(q.TextHtml, "{p:1:") {
		t.Errorf("HTML should not contain unclosed tag: %q", q.TextHtml)
	}
}

func TestParseAll_MixedRedBlueTruth(t *testing.T) {
	p := testParser

	// Real line from episode 5 tea party: Erika uses both red and blue truth
	lines := []string{
		"new_tea 5",
		"d2 [lv 0*\"46\"*\"54501382\"]`\"`[#][*]`{p:2:It's only possible for the crimes to have taken place during the single hour after midnight}!! `[@][lv 0*\"46\"*\"54501383\"]`{p:1:During that one hour, you were in the dining hall of the mansion}!! `[@][lv 0*\"46\"*\"54501384\"]`{p:2:Therefore, it was impossible for you to commit the crimes}!\" `[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if q.Episode != 5 || q.ContentType != "tea" {
		t.Errorf("got ep=%d ct=%q, want ep=5 ct=tea", q.Episode, q.ContentType)
	}
	if q.CharacterID != "46" || q.Character != CharacterNames["46"] {
		t.Errorf("character: got %q/%q, want 46/%s", q.CharacterID, q.Character, CharacterNames["46"])
	}

	// Both red and blue truth should be present in HTML
	if !strings.Contains(q.TextHtml, `class="red-truth"`) {
		t.Errorf("HTML missing red truth: %q", q.TextHtml)
	}
	if !strings.Contains(q.TextHtml, `class="blue-truth"`) {
		t.Errorf("HTML missing blue truth: %q", q.TextHtml)
	}

	// Plain text should contain the content without HTML tags
	if !strings.Contains(q.Text, "During that one hour") {
		t.Errorf("plain text missing red truth content: %q", q.Text)
	}
	if !strings.Contains(q.Text, "impossible for you to commit") {
		t.Errorf("plain text missing blue truth content: %q", q.Text)
	}
}

func TestParseAll_NestedTagsInTruth(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 2",
		`d2 [lv 0*"27"*"20700984"]` + "`" + `"You talk too much, you incompetent fool. ` + "`" + `[@][lv 0*"27"*"20700985"]` + "`" + `Then let me expand on my earlier move. ` + "`" + `[@][lv 0*"27"*"20700986"][#][*]` + "`" + `{p:1:The six definitely entered through {i:this front door}}!!" ` + "`" + `[\]`,
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.TextHtml, `<span class="red-truth">The six definitely entered through <em>this front door</em></span>`) {
		t.Errorf("HTML has broken nested tags: %q", q.TextHtml)
	}

	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text contains stray braces: %q", q.Text)
	}

	if !strings.Contains(q.Text, "The six definitely entered through this front door") {
		t.Errorf("plain text missing truth content: %q", q.Text)
	}
}

func TestParseAll_NestedNobrInTruth(t *testing.T) {
	p := testParser

	lines := []string{
		"new_episode 5",
		`d2 [lv 0*"28"*"52100552"]` + "`" + `"I'll respond. ` + "`" + `[@][lv 0*"28"*"52100553"][#][*]` + "`" + `{p:1:From {nobr:1 a.m.} to {nobr:3 a.m.}, the trio of Erika, Nanjo, and Gohda... ` + "`" + `[@][lv 0*"28"*"52100554"]` + "`" + `spent their time in the lounge on the first floor of the guesthouse}." ` + "`" + `[\]`,
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if !strings.Contains(q.TextHtml, `<span class="red-truth">From 1 a.m. to 3 a.m., the trio of Erika, Nanjo, and Gohda...`) {
		t.Errorf("HTML has broken nested nobr tags: %q", q.TextHtml)
	}
	if !strings.Contains(q.TextHtml, `the first floor of the guesthouse</span>`) {
		t.Errorf("HTML red truth span not closed properly: %q", q.TextHtml)
	}

	if strings.Contains(q.Text, "{") || strings.Contains(q.Text, "}") {
		t.Errorf("plain text contains stray braces: %q", q.Text)
	}

	if !strings.Contains(q.Text, "From 1 a.m. to 3 a.m.") {
		t.Errorf("plain text missing nobr content: %q", q.Text)
	}
}

func TestParseAll_AlphanumericAudioIDs(t *testing.T) {
	p := testParser

	tests := []struct {
		name         string
		line         string
		wantAudioID  string
		wantCharID   string
		textContains string
	}{
		{
			name:         "awase group voice",
			line:         "d2 [lv 0*\"99\"*\"awase0001\"]`\"Eeeeh, mackereeeel?!?! This is long enough surely.\"`[#][*][\\]",
			wantAudioID:  "awase0001",
			wantCharID:   "99",
			textContains: "mackereeeel",
		},
		{
			name:         "announcer voice",
			line:         "d [lv 0*\"99\"*\"anaf1001\"]`\"Our apologies for the delay. Boarding will now commence for Flight 201 to Niijima.\"`[\\]",
			wantAudioID:  "anaf1001",
			wantCharID:   "99",
			textContains: "Boarding",
		},
		{
			name:         "staff voice",
			line:         "d [lv 0*\"99\"*\"staf1001\"]`\"Boarding will now commence. As I call out the names on the passenger list.\"`[\\]",
			wantAudioID:  "staf1001",
			wantCharID:   "99",
			textContains: "passenger list",
		},
		{
			name:         "awase with underscore suffix",
			line:         "d2 [lv 0*\"00\"*\"awase6100_o\"][ak][text_speed_t 5]`\"\"With your fellow monsters.\"\" `[#][*][\\]",
			wantAudioID:  "awase6100_o",
			wantCharID:   "00",
			textContains: "fellow monsters",
		},
		{
			name:         "brackets before lv tag",
			line:         "d2 [ak][text_speed_t 10][text_fade_t 0][lv 1*\"60\"*\"60700401\"][#][*][lv 0*\"47\"*\"54600449\"]`\"*cackle*cackle*cackle*...... Yes, it is. {p:41:You used magic to create a golden flower petal inside an overturned cup.}\"`[|][\\]",
			wantAudioID:  "60700401, 54600449",
			wantCharID:   "60",
			textContains: "golden flower petal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := []string{
				"new_episode 1",
				tt.line,
			}
			quotes := p.ParseAll(lines)
			if len(quotes) == 0 {
				t.Fatal("expected at least 1 quote, got 0")
			}
			q := quotes[0]
			if q.AudioID != tt.wantAudioID {
				t.Errorf("audioID: got %q, want %q", q.AudioID, tt.wantAudioID)
			}
			if q.CharacterID != tt.wantCharID {
				t.Errorf("characterID: got %q, want %q", q.CharacterID, tt.wantCharID)
			}
			if !strings.Contains(q.Text, tt.textContains) {
				t.Errorf("text %q missing substring %q", q.Text, tt.textContains)
			}
		})
	}
}

func TestParseAll_OmakeNarratorLine(t *testing.T) {
	p := testParser

	lines := []string{
		"*o1_0",
		"*o1_2",
		"d `Jessica's piercing shriek rang out across the entire building.`[\\]",
	}
	quotes := p.ParseAll(lines)
	if len(quotes) == 0 {
		t.Fatal("expected at least 1 quote")
	}
	q := quotes[0]

	if q.Episode != 1 {
		t.Errorf("episode: got %d, want 1", q.Episode)
	}
	if q.ContentType != "omake" {
		t.Errorf("contentType: got %q, want omake", q.ContentType)
	}
	if q.CharacterID != "narrator" {
		t.Errorf("characterID: got %q, want narrator", q.CharacterID)
	}
	if !strings.Contains(q.Text, "Jessica") {
		t.Errorf("text missing expected content: %q", q.Text)
	}
}

func TestParseAll_CleanupPatterns(t *testing.T) {
	p := testParser

	tests := []struct {
		name     string
		line     string
		wantText string
		wantHtml string
	}{
		{
			name:     "backtick [@] separates multi-voice segments",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"First segment here. `[@][lv 0*\"10\"*\"10100002\"]`Second segment here too.\" `[\\]",
			wantText: "First segment here. Second segment here too.",
			wantHtml: "First segment here. Second segment here too.",
		},
		{
			name:     "backtick [\\] ends a dialogue line",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"This line ends with the backslash marker.\" `[\\]",
			wantText: "This line ends with the backslash marker.",
			wantHtml: "This line ends with the backslash marker.",
		},
		{
			name:     "backtick [|] acts as a page break marker",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"Before the page break marker. `[|][lv 0*\"10\"*\"10100002\"]`After the page break text here.\" `[\\]",
			wantText: "Before the page break marker. After the page break text here.",
			wantHtml: "Before the page break marker. After the page break text here.",
		},
		{
			name:     "bare [@] without backtick prefix",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"Text before bare marker. [@][lv 0*\"10\"*\"10100002\"]`More text after the marker.\" `[\\]",
			wantText: "Text before bare marker. More text after the marker.",
			wantHtml: "Text before bare marker. More text after the marker.",
		},
		{
			name:     "bare [\\] without backtick prefix",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"Text that uses bare backslash at end.\"[\\]",
			wantText: "Text that uses bare backslash at end.",
			wantHtml: "Text that uses bare backslash at end.",
		},
		{
			name:     "bare [|] without backtick prefix",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"Text with bare pipe break. [|][lv 0*\"10\"*\"10100002\"]`Continued after pipe.\" `[\\]",
			wantText: "Text with bare pipe break. Continued after pipe.",
			wantHtml: "Text with bare pipe break. Continued after pipe.",
		},
		{
			name:     "backtick-quote and quote-backtick pairs stripped",
			line:     "d [lv 0*\"10\"*\"10100001\"]`\"The quoted dialogue text is long enough here.\"`[\\]",
			wantText: "The quoted dialogue text is long enough here.",
			wantHtml: "The quoted dialogue text is long enough here.",
		},
		{
			name:     "all cleanup patterns in one line",
			line:     "d2 [lv 0*\"10\"*\"10100001\"]`\"First part of the text. `[@][lv 0*\"10\"*\"10100002\"]`Second part continues here. `[|][lv 0*\"10\"*\"10100003\"]`Third part after pipe break.\" `[\\]",
			wantText: "First part of the text. Second part continues here. Third part after pipe break.",
			wantHtml: "First part of the text. Second part continues here. Third part after pipe break.",
		},
		{
			name:     "narrator line cleanup patterns",
			line:     "d `The narrator speaks across segments. `[@]`More narration continues here. `[\\]",
			wantText: "The narrator speaks across segments. More narration continues here.",
			wantHtml: "The narrator speaks across segments. More narration continues here.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := []string{
				"new_episode 1",
				tt.line,
			}
			quotes := p.ParseAll(lines)
			if len(quotes) == 0 {
				t.Fatalf("expected at least 1 quote for input: %s", tt.line)
			}
			q := quotes[0]
			if q.Text != tt.wantText {
				t.Errorf("text:\n  got  %q\n  want %q", q.Text, tt.wantText)
			}
			if q.TextHtml != tt.wantHtml {
				t.Errorf("html:\n  got  %q\n  want %q", q.TextHtml, tt.wantHtml)
			}
		})
	}
}
