package moe.auaurora.quotes.data.remote.dto

data class CharactersResponseDto(
    val characters: Map<String, String>,
    val additional: Map<String, String>
)
