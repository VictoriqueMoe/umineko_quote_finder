package moe.auaurora.quotes.presentation.browse

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.BrowseData
import moe.auaurora.quotes.domain.model.FilterState
import moe.auaurora.quotes.domain.usecase.BrowseQuotesUseCase
import moe.auaurora.quotes.domain.usecase.GetCharactersUseCase

class BrowseViewModel(repository: QuoteRepository) : ViewModel() {

    private val browseUseCase = BrowseQuotesUseCase(repository)
    private val getCharactersUseCase = GetCharactersUseCase(repository)

    private val _uiState = MutableStateFlow<BrowseUiState>(BrowseUiState.Loading)
    val uiState: StateFlow<BrowseUiState> = _uiState.asStateFlow()

    private val _filters = MutableStateFlow(FilterState())
    val filters: StateFlow<FilterState> = _filters.asStateFlow()

    private val _characters = MutableStateFlow<Map<String, String>>(emptyMap())
    val characters: StateFlow<Map<String, String>> = _characters.asStateFlow()

    private var currentOffset = 0

    init {
        loadCharacters()
        loadBrowse()
    }

    private fun loadCharacters() {
        viewModelScope.launch {
            getCharactersUseCase().onSuccess { _characters.value = it }
        }
    }

    fun loadBrowse() {
        viewModelScope.launch {
            _uiState.value = BrowseUiState.Loading
            browseUseCase(_filters.value, offset = currentOffset)
                .onSuccess { _uiState.value = BrowseUiState.Success(it) }
                .onFailure { _uiState.value = BrowseUiState.Error(it.message ?: "Unknown error") }
        }
    }

    fun loadPage(offset: Int) {
        currentOffset = offset
        loadBrowse()
    }

    fun updateFilters(newFilters: FilterState) {
        _filters.value = newFilters
        currentOffset = 0
        loadBrowse()
    }
}

sealed class BrowseUiState {
    data object Loading : BrowseUiState()
    data class Success(val data: BrowseData) : BrowseUiState()
    data class Error(val message: String) : BrowseUiState()
}
