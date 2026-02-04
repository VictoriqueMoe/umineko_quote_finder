package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.ContextData

class GetContextUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(
        audioId: String,
        lang: String = "en",
        lines: Int = 5
    ): Result<ContextData> {
        return repository.getContext(audioId, lang, lines)
    }
}
