package moe.auaurora.quotes.presentation.browse

import android.content.ClipData
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.FilterList
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.ClipEntry
import androidx.compose.ui.platform.LocalClipboard
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.unit.dp
import kotlinx.coroutines.launch
import moe.auaurora.quotes.presentation.components.FilterBar
import moe.auaurora.quotes.presentation.components.PaginationControls
import moe.auaurora.quotes.presentation.components.QuoteCard
import moe.auaurora.quotes.ui.theme.*
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun BrowseScreen(
    viewModel: BrowseViewModel,
    audioPlayer: AudioPlayerManager,
    onQuoteClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
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
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(
                text = "Browse Dialogue",
                style = MaterialTheme.typography.headlineSmall,
                color = Gold
            )
            IconButton(onClick = { showFilters = !showFilters }) {
                Icon(
                    Icons.Filled.FilterList,
                    contentDescription = "Filters",
                    tint = GoldDark
                )
            }
        }

        AnimatedVisibility(visible = showFilters) {
            FilterBar(
                filters = filters,
                characters = characters,
                onFiltersChanged = { viewModel.updateFilters(it) }
            )
        }

        when (val state = uiState) {
            is BrowseUiState.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is BrowseUiState.Success -> {
                if (state.data.quotes.isEmpty()) {
                    Box(
                        modifier = Modifier.fillMaxSize(),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = "No dialogue found",
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
                            items = state.data.quotes,
                            key = { it.firstAudioId.ifEmpty { it.hashCode().toString() } }
                        ) { quote ->
                            QuoteCard(
                                quote = quote,
                                isPlaying = isAudioPlaying && currentAudioId == quote.firstAudioId,
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
                                    if (quote.firstAudioId.isNotEmpty()) {
                                        onQuoteClick(quote.firstAudioId)
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

            is BrowseUiState.Error -> {
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
