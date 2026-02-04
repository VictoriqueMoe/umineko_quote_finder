package moe.auaurora.quotes.presentation.stats

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.StatsData
import moe.auaurora.quotes.domain.usecase.GetStatsUseCase

class StatsViewModel(repository: QuoteRepository) : ViewModel() {

    private val getStatsUseCase = GetStatsUseCase(repository)

    private val _uiState = MutableStateFlow<StatsUiState>(StatsUiState.Loading)
    val uiState: StateFlow<StatsUiState> = _uiState.asStateFlow()

    private val _selectedEpisode = MutableStateFlow<Int?>(null)
    val selectedEpisode: StateFlow<Int?> = _selectedEpisode.asStateFlow()

    init {
        loadStats()
    }

    fun selectEpisode(episode: Int?) {
        _selectedEpisode.value = episode
        loadStats()
    }

    private fun loadStats() {
        viewModelScope.launch {
            _uiState.value = StatsUiState.Loading
            getStatsUseCase(_selectedEpisode.value)
                .onSuccess { _uiState.value = StatsUiState.Success(it) }
                .onFailure { _uiState.value = StatsUiState.Error(it.message ?: "Unknown error") }
        }
    }
}

sealed class StatsUiState {
    data object Loading : StatsUiState()
    data class Success(val data: StatsData) : StatsUiState()
    data class Error(val message: String) : StatsUiState()
}
