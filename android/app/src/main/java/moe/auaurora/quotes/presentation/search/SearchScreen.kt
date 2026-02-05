package moe.auaurora.quotes.presentation.search

import android.content.ClipData
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardActions
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.FilterList
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.ClipEntry
import androidx.compose.ui.platform.LocalClipboard
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.input.ImeAction
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import moe.auaurora.quotes.presentation.components.FilterBar
import moe.auaurora.quotes.presentation.components.PaginationControls
import moe.auaurora.quotes.presentation.components.QuoteCard
import moe.auaurora.quotes.ui.theme.*
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun SearchScreen(
    viewModel: SearchViewModel,
    audioPlayer: AudioPlayerManager,
    onQuoteClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    val query by viewModel.query.collectAsState()
    val filters by viewModel.filters.collectAsState()
    val characters by viewModel.characters.collectAsState()
    val currentAudioId by audioPlayer.currentAudioId.collectAsState()
    val isAudioPlaying by audioPlayer.isPlaying.collectAsState()
    val clipboard = LocalClipboard.current
    val scope = rememberCoroutineScope()
    var showFilters by remember { mutableStateOf(false) }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(BgVoid)
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        OutlinedTextField(
            value = query,
            onValueChange = { viewModel.updateQuery(it) },
            placeholder = {
                Text(
                    text = "Search quotes\u2026",
                    style = MaterialTheme.typography.bodyMedium.copy(
                        fontStyle = FontStyle.Italic
                    ),
                    color = TextMuted.copy(alpha = 0.6f)
                )
            },
            leadingIcon = {
                Icon(
                    Icons.Filled.Search,
                    contentDescription = null,
                    tint = Gold
                )
            },
            trailingIcon = {
                IconButton(onClick = { showFilters = !showFilters }) {
                    Icon(
                        Icons.Filled.FilterList,
                        contentDescription = "Filters",
                        tint = GoldDark
                    )
                }
            },
            keyboardOptions = KeyboardOptions(imeAction = ImeAction.Search),
            keyboardActions = KeyboardActions(onSearch = { viewModel.search() }),
            singleLine = true,
            textStyle = MaterialTheme.typography.bodyMedium.copy(color = TextPrimary),
            shape = RoundedCornerShape(0.dp),
            colors = OutlinedTextFieldDefaults.colors(
                unfocusedContainerColor = BgCard,
                focusedContainerColor = BgCard,
                unfocusedBorderColor = PurpleMuted,
                focusedBorderColor = Gold,
                cursorColor = Gold
            ),
            modifier = Modifier.fillMaxWidth()
        )

        AnimatedVisibility(visible = showFilters) {
            FilterBar(
                filters = filters,
                characters = characters,
                onFiltersChanged = { viewModel.updateFilters(it) }
            )
        }

        when (val state = uiState) {
            is SearchUiState.Idle -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = "Search for Umineko dialogue",
                        style = MaterialTheme.typography.bodyLarge.copy(
                            fontStyle = FontStyle.Italic
                        ),
                        color = TextMuted.copy(alpha = 0.6f)
                    )
                }
            }

            is SearchUiState.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is SearchUiState.Success -> {
                if (state.data.results.isEmpty()) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = "No results found",
                            style = MaterialTheme.typography.bodyLarge.copy(
                                fontStyle = FontStyle.Italic
                            ),
                            color = TextMuted
                        )
                    }
                } else {
                    LazyColumn(
                        modifier = Modifier.weight(1f),
                        verticalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        items(
                            items = state.data.results,
                            key = { it.quote.firstAudioId.ifEmpty { it.hashCode().toString() } }
                        ) { result ->
                            QuoteCard(
                                quote = result.quote,
                                isPlaying = isAudioPlaying && currentAudioId == result.quote.firstAudioId,
                                onPlayAudio = { q ->
                                    audioPlayer.playSingle(q)
                                },
                                onStopAudio = { audioPlayer.stop() },
                                onShare = { audioId ->
                                    scope.launch {
                                        clipboard.setClipEntry(ClipEntry(ClipData.newPlainText("url", "https://quotes.auaurora.moe/?quote=$audioId")))
                                    }
                                },
                                onClick = {
                                    if (result.quote.firstAudioId.isNotEmpty()) {
                                        onQuoteClick(result.quote.firstAudioId)
                                    }
                                }
                            )
                        }
                    }

                    PaginationControls(
                        offset = state.data.offset,
                        limit = state.data.limit,
                        total = state.data.total,
                        onPrevious = { viewModel.loadPage(state.data.offset - state.data.limit) },
                        onNext = { viewModel.loadPage(state.data.offset + state.data.limit) }
                    )
                }
            }

            is SearchUiState.Error -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Text(
                        text = state.message,
                        style = MaterialTheme.typography.bodyMedium,
                        color = RedTruth
                    )
                }
            }
        }
    }
}
