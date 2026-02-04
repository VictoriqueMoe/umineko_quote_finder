package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.FilterState
import moe.auaurora.quotes.domain.model.Quote

class GetRandomQuoteUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(filters: FilterState): Result<Quote> {
        return repository.random(
            lang = filters.language,
            character = filters.character.ifEmpty { null },
            episode = if (filters.episode > 0) filters.episode else null,
            truth = filters.truth.ifEmpty { null }
        )
    }
}
