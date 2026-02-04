package moe.auaurora.quotes.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class ContextResponseDto(
    val before: List<QuoteDto> = emptyList(),
    val quote: QuoteDto,
    val after: List<QuoteDto> = emptyList()
)
