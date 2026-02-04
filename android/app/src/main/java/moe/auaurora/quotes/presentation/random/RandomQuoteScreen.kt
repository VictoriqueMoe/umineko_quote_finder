package moe.auaurora.quotes.presentation.random

import androidx.compose.animation.AnimatedVisibility
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.FilterList
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.presentation.components.FilterBar
import moe.auaurora.quotes.ui.theme.*
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun RandomQuoteScreen(
    viewModel: RandomQuoteViewModel,
    audioPlayer: AudioPlayerManager,
    onQuoteClick: (String) -> Unit,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    val filters by viewModel.filters.collectAsState()
    val characters by viewModel.characters.collectAsState()
    var showFilters by remember { mutableStateOf(false) }

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(BgVoid)
            .padding(16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        Box(modifier = Modifier.fillMaxWidth()) {
            Text(
                text = "Random Quote",
                style = MaterialTheme.typography.headlineSmall,
                color = Gold
            )
            IconButton(
                onClick = { showFilters = !showFilters },
                modifier = Modifier.align(Alignment.CenterEnd)
            ) {
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

        Spacer(modifier = Modifier.height(8.dp))

        when (val state = uiState) {
            is RandomUiState.Loading -> {
                Box(
                    modifier = Modifier.weight(1f),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is RandomUiState.Success -> {
                Box(
                    modifier = Modifier.weight(1f),
                    contentAlignment = Alignment.Center
                ) {
                    Column(
                        modifier = Modifier
                            .fillMaxWidth()
                            .drawBehind {
                                // Gold left border
                                drawRect(
                                    color = Gold,
                                    topLeft = Offset(0f, 0f),
                                    size = Size(3.dp.toPx(), size.height)
                                )
                            }
                            .background(BgCard)
                            .border(width = 1.dp, color = PurpleMuted)
                            .clickable {
                                if (state.quote.firstAudioId.isNotEmpty()) {
                                    onQuoteClick(state.quote.firstAudioId)
                                }
                            }
                            .padding(24.dp),
                        horizontalAlignment = Alignment.CenterHorizontally,
                        verticalArrangement = Arrangement.spacedBy(12.dp)
                    ) {
                        // Label above the quote
                        Text(
                            text = "A Fragment from the Sea",
                            style = MaterialTheme.typography.labelMedium,
                            color = TextMuted,
                            textAlign = TextAlign.Center
                        )

                        // Decorative ornament above
                        Text(
                            text = "\u2767",
                            style = MaterialTheme.typography.headlineMedium,
                            color = Gold.copy(alpha = 0.35f),
                            textAlign = TextAlign.Center
                        )

                        // Character name
                        Text(
                            text = state.quote.character,
                            style = MaterialTheme.typography.labelLarge,
                            color = Gold,
                            textAlign = TextAlign.Center
                        )

                        // Quote text - featured style
                        Text(
                            text = state.quote.text,
                            style = MaterialTheme.typography.bodyLarge.copy(
                                fontStyle = FontStyle.Italic
                            ),
                            color = when {
                                state.quote.hasRedTruth -> RedTruth
                                state.quote.hasBlueTruth -> BlueTruth
                                else -> TextPrimary
                            },
                            textAlign = TextAlign.Center
                        )

                        // Decorative ornament below
                        Text(
                            text = "\u2767",
                            style = MaterialTheme.typography.headlineMedium,
                            color = Gold.copy(alpha = 0.35f),
                            textAlign = TextAlign.Center
                        )

                        // Episode label
                        if (state.quote.episode > 0) {
                            Text(
                                text = "Episode ${state.quote.episode}",
                                style = MaterialTheme.typography.labelSmall,
                                color = TextMuted,
                                textAlign = TextAlign.Center
                            )
                        }
                    }
                }
            }

            is RandomUiState.Error -> {
                Box(
                    modifier = Modifier.weight(1f),
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

        // "New Random Quote" button with gold gradient
        Box(
            modifier = Modifier
                .background(
                    brush = Brush.horizontalGradient(
                        colors = listOf(GoldDark, Gold)
                    ),
                    shape = RoundedCornerShape(0.dp)
                )
                .clickable { viewModel.loadRandom() }
                .padding(horizontal = 24.dp, vertical = 12.dp),
            contentAlignment = Alignment.Center
        ) {
            Text(
                text = "NEW RANDOM QUOTE",
                style = MaterialTheme.typography.labelLarge.copy(
                    fontWeight = FontWeight.SemiBold
                ),
                color = BgVoid
            )
        }
    }
}
