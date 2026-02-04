package moe.auaurora.quotes.presentation.context

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.ContextData
import moe.auaurora.quotes.domain.usecase.GetContextUseCase

class ContextViewModel(
    repository: QuoteRepository,
    private val audioId: String
) : ViewModel() {

    private val getContextUseCase = GetContextUseCase(repository)

    private val _uiState = MutableStateFlow<ContextUiState>(ContextUiState.Loading)
    val uiState: StateFlow<ContextUiState> = _uiState.asStateFlow()

    init {
        loadContext()
    }

    private fun loadContext() {
        viewModelScope.launch {
            _uiState.value = ContextUiState.Loading
            getContextUseCase(audioId)
                .onSuccess { _uiState.value = ContextUiState.Success(it) }
                .onFailure { _uiState.value = ContextUiState.Error(it.message ?: "Unknown error") }
        }
    }
}

sealed class ContextUiState {
    data object Loading : ContextUiState()
    data class Success(val data: ContextData) : ContextUiState()
    data class Error(val message: String) : ContextUiState()
}
