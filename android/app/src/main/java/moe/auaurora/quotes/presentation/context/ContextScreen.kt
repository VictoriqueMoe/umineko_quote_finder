package moe.auaurora.quotes.presentation.context

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.Stop
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.domain.model.Quote
import moe.auaurora.quotes.ui.theme.*
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun ContextScreen(
    viewModel: ContextViewModel,
    audioPlayer: AudioPlayerManager,
    onQuoteClick: (String) -> Unit,
    onBack: () -> Unit,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    val currentAudioId by audioPlayer.currentAudioId.collectAsState()
    val isAudioPlaying by audioPlayer.isPlaying.collectAsState()

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(BgVoid)
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        IconButton(onClick = onBack) {
            Icon(
                Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = "Back",
                tint = Gold
            )
        }

        Text(
            text = "Scene Context",
            style = MaterialTheme.typography.headlineSmall,
            color = Gold
        )

        when (val state = uiState) {
            is ContextUiState.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is ContextUiState.Success -> {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    verticalArrangement = Arrangement.spacedBy(4.dp)
                ) {
                    items(state.data.before) { quote ->
                        ContextQuoteLine(
                            quote = quote,
                            isHighlighted = false,
                            isPlaying = isAudioPlaying && currentAudioId == quote.firstAudioId,
                            audioPlayer = audioPlayer,
                            onClick = if (quote.firstAudioId.isNotEmpty()) {
                                { onQuoteClick(quote.firstAudioId) }
                            } else {
                                null
                            }
                        )
                    }

                    item {
                        ContextQuoteLine(
                            quote = state.data.quote,
                            isHighlighted = true,
                            isPlaying = isAudioPlaying && currentAudioId == state.data.quote.firstAudioId,
                            audioPlayer = audioPlayer,
                            onClick = null
                        )
                    }

                    items(state.data.after) { quote ->
                        ContextQuoteLine(
                            quote = quote,
                            isHighlighted = false,
                            isPlaying = isAudioPlaying && currentAudioId == quote.firstAudioId,
                            audioPlayer = audioPlayer,
                            onClick = if (quote.firstAudioId.isNotEmpty()) {
                                { onQuoteClick(quote.firstAudioId) }
                            } else {
                                null
                            }
                        )
                    }
                }
            }

            is ContextUiState.Error -> {
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

@Composable
private fun ContextQuoteLine(
    quote: Quote,
    isHighlighted: Boolean,
    isPlaying: Boolean,
    audioPlayer: AudioPlayerManager,
    onClick: (() -> Unit)? = null
) {
    val leftBorderColour = if (isHighlighted) Gold else PurpleMuted
    val bgColour = if (isHighlighted) Gold.copy(alpha = 0.08f) else BgCard.copy(alpha = 0.4f)

    Box(
        modifier = Modifier
            .fillMaxWidth()
            .drawBehind {
                // Left accent border
                drawRect(
                    color = leftBorderColour,
                    topLeft = Offset(0f, 0f),
                    size = Size(3.dp.toPx(), size.height)
                )
            }
            .background(bgColour)
            .border(width = 1.dp, color = if (isHighlighted) Gold.copy(alpha = 0.3f) else PurpleMuted.copy(alpha = 0.5f))
            .then(
                if (onClick != null) {
                    Modifier.clickable(onClick = onClick)
                } else {
                    Modifier
                }
            )
            .padding(start = 6.dp)
            .padding(12.dp)
    ) {
        Column(
            verticalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                // Character label
                Text(
                    text = quote.character,
                    style = MaterialTheme.typography.labelLarge,
                    color = Gold
                )

                // Play/stop button
                if (quote.audioId.isNotEmpty()) {
                    IconButton(
                        onClick = {
                            if (isPlaying) {
                                audioPlayer.stop()
                            } else {
                                audioPlayer.playSingle(quote)
                            }
                        }
                    ) {
                        Icon(
                            if (isPlaying) {
                                Icons.Filled.Stop
                            } else {
                                Icons.Filled.PlayArrow
                            },
                            contentDescription = if (isPlaying) "Stop" else "Play",
                            tint = if (isPlaying) Gold else GoldDark,
                            modifier = Modifier
                        )
                    }
                }
            }

            // Quote text
            Text(
                text = quote.text,
                style = MaterialTheme.typography.bodyMedium.copy(
                    fontStyle = FontStyle.Italic
                ),
                color = when {
                    quote.hasRedTruth -> RedTruth
                    quote.hasBlueTruth -> BlueTruth
                    else -> TextMuted
                }
            )
        }
    }
}
