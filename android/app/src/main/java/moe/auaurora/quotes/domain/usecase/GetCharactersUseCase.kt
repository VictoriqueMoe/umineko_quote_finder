package moe.auaurora.quotes.domain.usecase

import moe.auaurora.quotes.data.repository.QuoteRepository

class GetCharactersUseCase(private val repository: QuoteRepository) {
    suspend operator fun invoke(): Result<Map<String, String>> {
        return repository.getCharacters()
    }
}
