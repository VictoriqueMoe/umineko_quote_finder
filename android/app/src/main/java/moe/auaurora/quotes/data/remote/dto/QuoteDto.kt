package moe.auaurora.quotes.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class QuoteDto(
    val text: String = "",
    val textHtml: String = "",
    val characterId: String = "",
    val character: String = "",
    val audioId: String = "",
    val audioCharMap: Map<String, String>? = null,
    val episode: Int = 0,
    val contentType: String = "",
    val hasRedTruth: Boolean = false,
    val hasBlueTruth: Boolean = false
)
