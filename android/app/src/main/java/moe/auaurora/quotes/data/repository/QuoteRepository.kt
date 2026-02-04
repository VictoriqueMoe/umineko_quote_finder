package moe.auaurora.quotes.data.repository

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import moe.auaurora.quotes.data.remote.ApiClient
import moe.auaurora.quotes.data.remote.UminekoApi
import moe.auaurora.quotes.data.remote.dto.*
import moe.auaurora.quotes.domain.model.*

class QuoteRepository(
    private val api: UminekoApi = ApiClient.api
) {

    suspend fun search(
        query: String,
        lang: String,
        limit: Int,
        offset: Int,
        character: String?,
        episode: Int?,
        truth: String?
    ): Result<SearchData> = withContext(Dispatchers.IO) {
        try {
            val response = api.search(query, lang, limit, offset, character, episode, truth)
            if (response.isSuccessful && response.body() != null) {
                val body = response.body()!!
                Result.success(
                    SearchData(
                        results = body.results.map { SearchResult(it.quote.toDomain(), it.score) },
                        total = body.total,
                        limit = body.limit,
                        offset = body.offset
                    )
                )
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun random(
        lang: String,
        character: String?,
        episode: Int?,
        truth: String?
    ): Result<Quote> = withContext(Dispatchers.IO) {
        try {
            val response = api.random(lang, character, episode, truth)
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!.toDomain())
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun browse(
        lang: String,
        limit: Int,
        offset: Int,
        character: String?,
        episode: Int?,
        truth: String?
    ): Result<BrowseData> = withContext(Dispatchers.IO) {
        try {
            val response = api.browse(lang, limit, offset, character, episode, truth)
            if (response.isSuccessful && response.body() != null) {
                val body = response.body()!!
                Result.success(
                    BrowseData(
                        characterId = body.characterId,
                        character = body.character,
                        quotes = body.quotes.map { it.toDomain() },
                        total = body.total,
                        limit = body.limit,
                        offset = body.offset
                    )
                )
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getQuote(
        audioId: String,
        lang: String
    ): Result<Quote> = withContext(Dispatchers.IO) {
        try {
            val response = api.getQuote(audioId, lang)
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!.toDomain())
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getContext(
        audioId: String,
        lang: String,
        lines: Int = 5
    ): Result<ContextData> = withContext(Dispatchers.IO) {
        try {
            val response = api.getContext(audioId, lang, lines)
            if (response.isSuccessful && response.body() != null) {
                val body = response.body()!!
                Result.success(
                    ContextData(
                        before = body.before.map { it.toDomain() },
                        quote = body.quote.toDomain(),
                        after = body.after.map { it.toDomain() }
                    )
                )
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getCharacters(): Result<Map<String, String>> = withContext(Dispatchers.IO) {
        try {
            val response = api.getCharacters()
            if (response.isSuccessful && response.body() != null) {
                Result.success(response.body()!!)
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getStats(episode: Int?): Result<StatsData> = withContext(Dispatchers.IO) {
        try {
            val response = api.getStats(episode)
            if (response.isSuccessful && response.body() != null) {
                val body = response.body()!!
                Result.success(
                    StatsData(
                        topSpeakers = body.topSpeakers.map { it.toDomain() },
                        linesPerEpisode = body.linesPerEpisode?.map { it.toDomain() } ?: emptyList(),
                        truthPerEpisode = body.truthPerEpisode?.map { it.toDomain() } ?: emptyList(),
                        interactions = body.interactions.map { it.toDomain() },
                        characterPresence = body.characterPresence?.map { it.toDomain() } ?: emptyList(),
                        characterNames = body.characterNames,
                        episodeNames = body.episodeNames
                    )
                )
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    suspend fun getConfig(): Result<ConfigData> = withContext(Dispatchers.IO) {
        try {
            val response = api.getConfig()
            if (response.isSuccessful && response.body() != null) {
                Result.success(ConfigData(hasAudio = response.body()!!.hasAudio))
            } else {
                Result.failure(Exception("API error: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}

private fun QuoteDto.toDomain() = Quote(
    text = text,
    textHtml = textHtml,
    characterId = characterId,
    character = character,
    audioId = audioId,
    episode = episode,
    contentType = contentType,
    hasRedTruth = hasRedTruth,
    hasBlueTruth = hasBlueTruth
)

private fun SpeakerStatDto.toDomain() = SpeakerStat(
    characterId = characterId,
    name = name,
    count = count
)

private fun EpisodeTruthDto.toDomain() = EpisodeTruth(
    episode = episode,
    red = red,
    blue = blue
)

private fun InteractionPairDto.toDomain() = InteractionPair(
    charA = charA,
    charB = charB,
    nameA = nameA,
    nameB = nameB,
    count = count
)

private fun EpisodeCharacterLinesDto.toDomain() = EpisodeCharacterLines(
    episode = episode,
    episodeName = episodeName,
    characters = characters
)

private fun CharacterPresenceDto.toDomain() = CharacterPresence(
    characterId = characterId,
    name = name,
    episodes = episodes
)
