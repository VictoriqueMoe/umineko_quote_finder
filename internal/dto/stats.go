package dto

type (
	SpeakerStat struct {
		// Numeric character identifier
		CharacterID string `json:"characterId" example:"10"`
		// Display name of the character
		Name string `json:"name" example:"Ushiromiya Battler"`
		// Number of dialogue lines
		Count int `json:"count" example:"4304"`
	}

	EpisodeTruth struct {
		// Episode number (1-8)
		Episode int `json:"episode" example:"2"`
		// Number of red truth statements
		Red int `json:"red" example:"38"`
		// Number of blue truth statements
		Blue int `json:"blue" example:"0"`
		// Number of gold truth statements
		Gold int `json:"gold" example:"0"`
		// Number of purple truth statements
		Purple int `json:"purple" example:"0"`
	}

	InteractionPair struct {
		// First character ID
		CharA string `json:"charA" example:"10"`
		// Second character ID
		CharB string `json:"charB" example:"27"`
		// First character name
		NameA string `json:"nameA" example:"Ushiromiya Battler"`
		// Second character name
		NameB string `json:"nameB" example:"Beatrice"`
		// Number of shared dialogue exchanges
		Count int `json:"count" example:"1701"`
	}

	EpisodeCharacterLines struct {
		// Episode number (1-8)
		Episode int `json:"episode" example:"1"`
		// Episode title
		EpisodeName string `json:"episodeName" example:"Legend"`
		// Map of character ID to line count
		Characters map[string]int `json:"characters"`
	}

	CharacterPresence struct {
		// Numeric character identifier
		CharacterID string `json:"characterId" example:"10"`
		// Display name of the character
		Name string `json:"name" example:"Ushiromiya Battler"`
		// List of episode numbers the character appears in
		Episodes []int `json:"episodes"`
	}

	StatsResult struct {
		// Top speakers by dialogue line count
		TopSpeakers []SpeakerStat `json:"topSpeakers"`
		// Dialogue line counts per episode per character
		LinesPerEpisode []EpisodeCharacterLines `json:"linesPerEpisode"`
		// Truth statement counts per episode
		TruthPerEpisode []EpisodeTruth `json:"truthPerEpisode"`
		// Top character interaction pairs
		Interactions []InteractionPair `json:"interactions"`
		// Map of character pair key to interaction count
		InteractionCounts map[string]int `json:"interactionCounts"`
		// Characters and which episodes they appear in
		CharacterPresence []CharacterPresence `json:"characterPresence"`
		// Map of character ID to display name
		CharacterNames map[string]string `json:"characterNames" example:"10:Ushiromiya Battler,27:Beatrice"`
		// Map of episode number to episode title
		EpisodeNames map[int]string `json:"episodeNames"`
	}

	HigurashiStatsResult struct {
		// Top speakers by dialogue line count
		TopSpeakers []SpeakerStat `json:"topSpeakers"`
		// Dialogue line counts per arc per character
		LinesPerArc map[string]map[string]int `json:"linesPerArc"`
		// Top character interaction pairs
		Interactions []InteractionPair `json:"interactions"`
		// Map of character pair key to interaction count
		InteractionCounts map[string]int `json:"interactionCounts"`
		// Map of character ID to display name
		CharacterNames map[string]string `json:"characterNames"`
	}
)
