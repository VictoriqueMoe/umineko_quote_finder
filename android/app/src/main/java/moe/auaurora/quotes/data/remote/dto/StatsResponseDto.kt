package moe.auaurora.quotes.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class StatsResponseDto(
    val topSpeakers: List<SpeakerStatDto> = emptyList(),
    val linesPerEpisode: List<EpisodeCharacterLinesDto>? = null,
    val truthPerEpisode: List<EpisodeTruthDto>? = null,
    val interactions: List<InteractionPairDto> = emptyList(),
    val characterPresence: List<CharacterPresenceDto>? = null,
    val characterNames: Map<String, String> = emptyMap(),
    val episodeNames: Map<String, String> = emptyMap()
)

@Serializable
data class SpeakerStatDto(
    val characterId: String = "",
    val name: String = "",
    val count: Int = 0
)

@Serializable
data class EpisodeTruthDto(
    val episode: Int = 0,
    val red: Int = 0,
    val blue: Int = 0
)

@Serializable
data class InteractionPairDto(
    val charA: String = "",
    val charB: String = "",
    val nameA: String = "",
    val nameB: String = "",
    val count: Int = 0
)

@Serializable
data class EpisodeCharacterLinesDto(
    val episode: Int = 0,
    val episodeName: String = "",
    val characters: Map<String, Int> = emptyMap()
)

@Serializable
data class CharacterPresenceDto(
    val characterId: String = "",
    val name: String = "",
    val episodes: List<Int> = emptyList()
)
