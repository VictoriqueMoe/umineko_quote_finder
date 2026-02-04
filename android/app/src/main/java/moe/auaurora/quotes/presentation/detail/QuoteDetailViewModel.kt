package moe.auaurora.quotes.presentation.detail

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.ConfigData
import moe.auaurora.quotes.domain.model.Quote

class QuoteDetailViewModel(
    private val repository: QuoteRepository,
    private val audioId: String
) : ViewModel() {

    private val _uiState = MutableStateFlow<DetailUiState>(DetailUiState.Loading)
    val uiState: StateFlow<DetailUiState> = _uiState.asStateFlow()

    private val _config = MutableStateFlow(ConfigData(hasAudio = false))
    val config: StateFlow<ConfigData> = _config.asStateFlow()

    init {
        loadQuote()
        loadConfig()
    }

    private fun loadQuote() {
        viewModelScope.launch {
            _uiState.value = DetailUiState.Loading
            repository.getQuote(audioId, "en")
                .onSuccess { _uiState.value = DetailUiState.Success(it) }
                .onFailure { _uiState.value = DetailUiState.Error(it.message ?: "Unknown error") }
        }
    }

    private fun loadConfig() {
        viewModelScope.launch {
            repository.getConfig().onSuccess { _config.value = it }
        }
    }
}

sealed class DetailUiState {
    data object Loading : DetailUiState()
    data class Success(val quote: Quote) : DetailUiState()
    data class Error(val message: String) : DetailUiState()
}
