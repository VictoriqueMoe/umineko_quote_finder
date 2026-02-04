package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.FilterState
import moe.auaurora.quotes.domain.model.SearchData

class SearchQuotesUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(
        query: String,
        filters: FilterState,
        limit: Int = 30,
        offset: Int = 0
    ): Result<SearchData> {
        return repository.search(
            query = query,
            lang = filters.language,
            limit = limit,
            offset = offset,
            character = filters.character.ifEmpty { null },
            episode = if (filters.episode > 0) filters.episode else null,
            truth = filters.truth.ifEmpty { null }
        )
    }
}
