package moe.auaurora.quotes.domain.model

data class Quote(
    val text: String,
    val textHtml: String,
    val characterId: String,
    val character: String,
    val audioId: String,
    val episode: Int,
    val contentType: String,
    val hasRedTruth: Boolean = false,
    val hasBlueTruth: Boolean = false
)

data class SearchData(
    val results: List<SearchResult>,
    val total: Int,
    val limit: Int,
    val offset: Int
)

data class SearchResult(
    val quote: Quote,
    val score: Int
)

data class BrowseData(
    val characterId: String,
    val character: String,
    val quotes: List<Quote>,
    val total: Int,
    val limit: Int,
    val offset: Int
)

data class ContextData(
    val before: List<Quote>,
    val quote: Quote,
    val after: List<Quote>
)

data class ConfigData(
    val hasAudio: Boolean
)

data class FilterState(
    val language: String = "en",
    val character: String = "",
    val episode: Int = 0,
    val truth: String = ""
)
