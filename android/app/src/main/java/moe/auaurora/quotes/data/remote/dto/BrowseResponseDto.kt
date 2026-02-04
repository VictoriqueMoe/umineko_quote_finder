package moe.auaurora.quotes.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class BrowseResponseDto(
    val characterId: String = "",
    val character: String = "",
    val quotes: List<QuoteDto> = emptyList(),
    val total: Int = 0,
    val limit: Int = 0,
    val offset: Int = 0
)
