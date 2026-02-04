package moe.auaurora.quotes.domain.model

data class StatsData(
    val topSpeakers: List<SpeakerStat>,
    val linesPerEpisode: List<EpisodeCharacterLines>,
    val truthPerEpisode: List<EpisodeTruth>,
    val interactions: List<InteractionPair>,
    val characterPresence: List<CharacterPresence>,
    val characterNames: Map<String, String>,
    val episodeNames: Map<String, String>
)

data class SpeakerStat(
    val characterId: String,
    val name: String,
    val count: Int
)

data class EpisodeTruth(
    val episode: Int,
    val red: Int,
    val blue: Int
)

data class InteractionPair(
    val charA: String,
    val charB: String,
    val nameA: String,
    val nameB: String,
    val count: Int
)

data class EpisodeCharacterLines(
    val episode: Int,
    val episodeName: String,
    val characters: Map<String, Int>
)

data class CharacterPresence(
    val characterId: String,
    val name: String,
    val episodes: List<Int>
)
