package moe.auaurora.quotes.data.remote.dto

import kotlinx.serialization.Serializable

@Serializable
data class ConfigDto(
    val hasAudio: Boolean = false
)
