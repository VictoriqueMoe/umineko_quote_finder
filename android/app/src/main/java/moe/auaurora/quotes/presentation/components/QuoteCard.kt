package moe.auaurora.quotes.presentation.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.PlayArrow
import androidx.compose.material.icons.filled.Share
import androidx.compose.material.icons.filled.Stop
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextOverflow
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import moe.auaurora.quotes.domain.model.Quote
import moe.auaurora.quotes.ui.theme.*

@Composable
fun QuoteCard(
    quote: Quote,
    isPlaying: Boolean = false,
    onPlayAudio: ((String, String) -> Unit)? = null,
    onStopAudio: (() -> Unit)? = null,
    onShare: ((String) -> Unit)? = null,
    onClick: (() -> Unit)? = null,
    modifier: Modifier = Modifier
) {
    val cardBrush = Brush.linearGradient(
        colors = listOf(BgCard, BgCard.copy(alpha = 0.8f)),
        start = Offset(0f, 0f),
        end = Offset(Float.POSITIVE_INFINITY, Float.POSITIVE_INFINITY)
    )

    Box(
        modifier = modifier
            .fillMaxWidth()
            .drawBehind {
                // Gold left border
                drawRect(
                    color = Gold,
                    topLeft = Offset(0f, 0f),
                    size = androidx.compose.ui.geometry.Size(3.dp.toPx(), size.height)
                )
            }
            .background(cardBrush)
            .border(width = 1.dp, color = PurpleMuted)
            .then(
                if (onClick != null) {
                    Modifier.clickable(onClick = onClick)
                } else {
                    Modifier
                }
            )
            .padding(start = 6.dp)
    ) {
        // Decorative quotation mark
        Text(
            text = "\u201C",
            style = MaterialTheme.typography.headlineLarge.copy(
                fontSize = 48.sp,
                color = Gold.copy(alpha = 0.12f)
            ),
            modifier = Modifier
                .align(Alignment.TopStart)
                .padding(start = 4.dp)
        )

        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = quote.character,
                    style = MaterialTheme.typography.labelLarge,
                    color = Gold
                )

                Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                    if (quote.episode > 0) {
                        EpisodeBadge(episode = quote.episode)
                    }
                    if (quote.hasRedTruth) {
                        TruthBadge(text = "Red", colour = RedTruth)
                    }
                    if (quote.hasBlueTruth) {
                        TruthBadge(text = "Blue", colour = BlueTruth)
                    }
                }
            }

            Text(
                text = quote.text,
                style = MaterialTheme.typography.bodyMedium.copy(
                    fontFamily = CormorantGaramond,
                    fontStyle = FontStyle.Italic,
                    fontWeight = if (quote.hasRedTruth || quote.hasBlueTruth) {
                        FontWeight.SemiBold
                    } else {
                        FontWeight.Normal
                    }
                ),
                color = when {
                    quote.hasRedTruth -> RedTruth
                    quote.hasBlueTruth -> BlueTruth
                    else -> TextPrimary
                },
                maxLines = 6,
                overflow = TextOverflow.Ellipsis
            )

            if (quote.firstAudioId.isNotEmpty()) {
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.End
                ) {
                    if (onShare != null) {
                        IconButton(onClick = { onShare(quote.firstAudioId) }) {
                            Icon(
                                Icons.Filled.Share,
                                contentDescription = "Share",
                                tint = GoldDark,
                                modifier = Modifier.size(20.dp)
                            )
                        }
                    }

                    if (onPlayAudio != null) {
                        if (isPlaying) {
                            IconButton(onClick = { onStopAudio?.invoke() }) {
                                Icon(
                                    Icons.Filled.Stop,
                                    contentDescription = "Stop",
                                    tint = Gold,
                                    modifier = Modifier.size(20.dp)
                                )
                            }
                        } else {
                            IconButton(
                                onClick = { onPlayAudio(quote.characterId, quote.audioId) }
                            ) {
                                Icon(
                                    Icons.Filled.PlayArrow,
                                    contentDescription = "Play",
                                    tint = GoldDark,
                                    modifier = Modifier.size(20.dp)
                                )
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
private fun EpisodeBadge(episode: Int) {
    Box(
        modifier = Modifier
            .border(1.dp, PurpleMuted, RoundedCornerShape(2.dp))
            .padding(horizontal = 6.dp, vertical = 2.dp)
    ) {
        Text(
            text = "EP$episode",
            style = MaterialTheme.typography.labelSmall,
            color = TextMuted
        )
    }
}

@Composable
private fun TruthBadge(text: String, colour: Color) {
    Box(
        modifier = Modifier
            .border(1.dp, colour.copy(alpha = 0.4f), RoundedCornerShape(2.dp))
            .background(colour.copy(alpha = 0.1f), RoundedCornerShape(2.dp))
            .padding(horizontal = 6.dp, vertical = 2.dp)
    ) {
        Text(
            text = text,
            style = MaterialTheme.typography.labelSmall,
            color = colour
        )
    }
}
