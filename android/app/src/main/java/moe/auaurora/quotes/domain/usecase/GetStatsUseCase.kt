package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.StatsData

class GetStatsUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(episode: Int? = null): Result<StatsData> {
        return repository.getStats(episode)
    }
}
