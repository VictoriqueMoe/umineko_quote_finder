package moe.auaurora.quotes.presentation.search

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.domain.model.FilterState
import moe.auaurora.quotes.domain.model.SearchData
import moe.auaurora.quotes.domain.usecase.GetCharactersUseCase
import moe.auaurora.quotes.domain.usecase.SearchQuotesUseCase

class SearchViewModel(repository: QuoteRepository) : ViewModel() {

    private val searchUseCase = SearchQuotesUseCase(repository)
    private val getCharactersUseCase = GetCharactersUseCase(repository)

    private val _uiState = MutableStateFlow<SearchUiState>(SearchUiState.Idle)
    val uiState: StateFlow<SearchUiState> = _uiState.asStateFlow()

    private val _filters = MutableStateFlow(FilterState())
    val filters: StateFlow<FilterState> = _filters.asStateFlow()

    private val _query = MutableStateFlow("")
    val query: StateFlow<String> = _query.asStateFlow()

    private val _characters = MutableStateFlow<Map<String, String>>(emptyMap())
    val characters: StateFlow<Map<String, String>> = _characters.asStateFlow()

    private var currentOffset = 0

    init {
        loadCharacters()
    }

    private fun loadCharacters() {
        viewModelScope.launch {
            getCharactersUseCase().onSuccess { _characters.value = it }
        }
    }

    fun updateQuery(newQuery: String) {
        _query.value = newQuery
    }

    fun search() {
        val q = _query.value.trim()
        if (q.isEmpty()) {
            return
        }
        currentOffset = 0
        executeSearch()
    }

    fun loadPage(offset: Int) {
        currentOffset = offset
        executeSearch()
    }

    fun updateFilters(newFilters: FilterState) {
        _filters.value = newFilters
        if (_query.value.isNotEmpty()) {
            currentOffset = 0
            executeSearch()
        }
    }

    private fun executeSearch() {
        viewModelScope.launch {
            _uiState.value = SearchUiState.Loading
            searchUseCase(_query.value, _filters.value, offset = currentOffset)
                .onSuccess { _uiState.value = SearchUiState.Success(it) }
                .onFailure { _uiState.value = SearchUiState.Error(it.message ?: "Unknown error") }
        }
    }
}

sealed class SearchUiState {
    data object Idle : SearchUiState()
    data object Loading : SearchUiState()
    data class Success(val data: SearchData) : SearchUiState()
    data class Error(val message: String) : SearchUiState()
}
