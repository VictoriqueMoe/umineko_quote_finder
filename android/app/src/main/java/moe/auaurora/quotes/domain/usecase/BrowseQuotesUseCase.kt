package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.BrowseData
import moe.auaurora.quotes.domain.model.FilterState

class BrowseQuotesUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(
        filters: FilterState,
        limit: Int = 30,
        offset: Int = 0
    ): Result<BrowseData> {
        return repository.browse(
            lang = filters.language,
            limit = limit,
            offset = offset,
            character = filters.character.ifEmpty { null },
            episode = if (filters.episode > 0) filters.episode else null,
            truth = filters.truth.ifEmpty { null }
        )
    }
}
