package moe.auaurora.quotes.presentation.detail

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.ui.theme.*
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun QuoteDetailScreen(
    viewModel: QuoteDetailViewModel,
    audioPlayer: AudioPlayerManager,
    onViewContext: (String) -> Unit,
    onBack: () -> Unit,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    val config by viewModel.config.collectAsState()
    val currentAudioId by audioPlayer.currentAudioId.collectAsState()
    val isAudioPlaying by audioPlayer.isPlaying.collectAsState()

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(BgVoid)
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        IconButton(onClick = onBack) {
            Icon(
                Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = "Back",
                tint = Gold
            )
        }

        when (val state = uiState) {
            is DetailUiState.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is DetailUiState.Success -> {
                val quote = state.quote

                // Quote detail card
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .drawBehind {
                            drawRect(
                                color = Gold,
                                topLeft = Offset(0f, 0f),
                                size = Size(3.dp.toPx(), size.height)
                            )
                        }
                        .background(BgCard)
                        .border(width = 1.dp, color = PurpleMuted)
                        .padding(20.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    // Character name
                    Text(
                        text = quote.character,
                        style = MaterialTheme.typography.headlineMedium,
                        color = Gold
                    )

                    // Episode label
                    if (quote.episode > 0) {
                        Text(
                            text = "Episode ${quote.episode}",
                            style = MaterialTheme.typography.labelLarge,
                            color = TextMuted
                        )
                    }

                    // Full quote text
                    Text(
                        text = quote.text,
                        style = MaterialTheme.typography.bodyLarge.copy(
                            fontStyle = FontStyle.Italic
                        ),
                        color = when {
                            quote.hasRedTruth -> RedTruth
                            quote.hasBlueTruth -> BlueTruth
                            else -> TextPrimary
                        }
                    )
                }

                // Action buttons
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    if (config.hasAudio && quote.audioId.isNotEmpty()) {
                        val playing = isAudioPlaying && currentAudioId == quote.audioId

                        GothicActionButton(
                            text = if (playing) "STOP" else "PLAY AUDIO",
                            icon = {
                                Icon(
                                    if (playing) {
                                        Icons.Filled.Stop
                                    } else {
                                        Icons.Filled.PlayArrow
                                    },
                                    contentDescription = if (playing) "Stop" else "Play",
                                    tint = Gold
                                )
                            },
                            onClick = {
                                if (playing) {
                                    audioPlayer.stop()
                                } else {
                                    audioPlayer.playSingle(quote.characterId, quote.audioId)
                                }
                            }
                        )
                    }

                    if (quote.audioId.isNotEmpty()) {
                        GothicActionButton(
                            text = "VIEW CONTEXT",
                            onClick = { onViewContext(quote.audioId) }
                        )
                    }
                }
            }

            is DetailUiState.Error -> {
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
private fun GothicActionButton(
    text: String,
    onClick: () -> Unit,
    icon: @Composable (() -> Unit)? = null
) {
    Row(
        modifier = Modifier
            .background(BgCard, RoundedCornerShape(0.dp))
            .border(width = 1.dp, color = Gold, shape = RoundedCornerShape(0.dp))
            .clickable(onClick = onClick)
            .padding(horizontal = 16.dp, vertical = 10.dp),
        horizontalArrangement = Arrangement.spacedBy(8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        if (icon != null) {
            icon()
        }
        Text(
            text = text,
            style = MaterialTheme.typography.labelLarge.copy(
                fontWeight = FontWeight.SemiBold
            ),
            color = Gold
        )
    }
}
