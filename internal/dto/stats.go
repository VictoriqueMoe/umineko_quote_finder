package dto

type SpeakerStat struct {
	// Numeric character identifier
	CharacterID string `json:"characterId" example:"10"`
	// Display name of the character
	Name string `json:"name" example:"Ushiromiya Battler"`
	// Number of dialogue lines
	Count int `json:"count" example:"4304"`
}

type EpisodeTruth struct {
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

type InteractionPair struct {
	// First character ID
	CharA string `json:"charA" example:"10"`
	// Second character ID
	CharB string `json:"charB" example:"27"`
	// First character name
	NameA string `json:"nameA" example:"Ushiromiya Battler"`
	// Second character name
	NameB string `json:"nameB" example:"Beatrice"`
	// Number of adjacent dialogue exchanges
	Count int `json:"count" example:"1701"`
}

type EpisodeCharacterLines struct {
	// Episode number (1-8)
	Episode int `json:"episode" example:"1"`
	// Episode title
	EpisodeName string `json:"episodeName" example:"Legend"`
	// Map of character ID to line count
	Characters map[string]int `json:"characters"`
}

type CharacterPresence struct {
	// Numeric character identifier
	CharacterID string `json:"characterId" example:"10"`
	// Display name of the character
	Name string `json:"name" example:"Ushiromiya Battler"`
	// Line counts per episode (index 0 = ep 1)
	Episodes []int `json:"episodes"`
}

type StatsResult struct {
	// Characters ranked by dialogue line count
	TopSpeakers []SpeakerStat `json:"topSpeakers"`
	// Per-episode line breakdown by character
	LinesPerEpisode []EpisodeCharacterLines `json:"linesPerEpisode"`
	// Truth counts per episode
	TruthPerEpisode []EpisodeTruth `json:"truthPerEpisode"`
	// Most frequent character dialogue pairings
	Interactions []InteractionPair `json:"interactions"`
	// Full interaction map keyed by "charA|charB"
	InteractionCounts map[string]int `json:"interactionCounts"`
	// Top characters with per-episode line counts
	CharacterPresence []CharacterPresence `json:"characterPresence"`
	// Map of character ID to display name
	CharacterNames map[string]string `json:"characterNames" example:"10:Ushiromiya Battler,27:Beatrice"`
	// Map of episode number to title
	EpisodeNames map[int]string `json:"episodeNames"`
}
