package moe.auaurora.quotes.presentation.stats

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.HorizontalDivider
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.domain.model.StatsData
import moe.auaurora.quotes.ui.theme.*

@OptIn(ExperimentalLayoutApi::class)
@Composable
fun StatsScreen(
    viewModel: StatsViewModel,
    modifier: Modifier = Modifier
) {
    val uiState by viewModel.uiState.collectAsState()
    val selectedEpisode by viewModel.selectedEpisode.collectAsState()

    Column(
        modifier = modifier
            .fillMaxSize()
            .background(BgVoid)
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        Text(
            text = "Statistics",
            style = MaterialTheme.typography.headlineSmall,
            color = Gold
        )

        // Episode filter chips - gothic style
        FlowRow(
            horizontalArrangement = Arrangement.spacedBy(6.dp),
        ) {
            GothicChip(
                text = "All",
                selected = selectedEpisode == null,
                onClick = { viewModel.selectEpisode(null) }
            )
            for (ep in 1..8) {
                GothicChip(
                    text = "EP$ep",
                    selected = selectedEpisode == ep,
                    onClick = { viewModel.selectEpisode(ep) }
                )
            }
        }

        when (val state = uiState) {
            is StatsUiState.Loading -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    CircularProgressIndicator(color = Gold)
                }
            }

            is StatsUiState.Success -> {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    item { TopSpeakersCard(state.data) }
                    item { TruthCard(state.data) }
                    item { InteractionsCard(state.data) }
                }
            }

            is StatsUiState.Error -> {
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
private fun GothicChip(
    text: String,
    selected: Boolean,
    onClick: () -> Unit
) {
    val borderColour = if (selected) Gold else PurpleMuted
    val bgColour = if (selected) Gold.copy(alpha = 0.15f) else BgCard.copy(alpha = 0.3f)
    val textColour = if (selected) Gold else TextMuted

    Box(
        modifier = Modifier
            .padding(vertical = 3.dp)
            .clip(RoundedCornerShape(2.dp))
            .border(1.dp, borderColour, RoundedCornerShape(2.dp))
            .background(bgColour, RoundedCornerShape(2.dp))
            .clickable(onClick = onClick)
            .padding(horizontal = 10.dp, vertical = 5.dp)
    ) {
        Text(
            text = text,
            style = MaterialTheme.typography.labelSmall,
            color = textColour
        )
    }
}

@Composable
private fun TopSpeakersCard(stats: StatsData) {
    Box(
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
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Text(
                text = "Top Speakers",
                style = MaterialTheme.typography.titleMedium,
                color = Gold,
                fontWeight = FontWeight.Bold
            )
            stats.topSpeakers.take(15).forEachIndexed { index, speaker ->
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(8.dp)
                    ) {
                        Text(
                            text = "${index + 1}.",
                            style = MaterialTheme.typography.bodyMedium,
                            color = Gold,
                            fontWeight = FontWeight.SemiBold
                        )
                        Text(
                            text = speaker.name,
                            style = MaterialTheme.typography.bodyMedium,
                            color = TextPrimary
                        )
                    }
                    Text(
                        text = "${speaker.count} lines",
                        style = MaterialTheme.typography.bodyMedium,
                        color = GoldDark
                    )
                }
            }
        }
    }
}

@Composable
private fun TruthCard(stats: StatsData) {
    if (stats.truthPerEpisode.isEmpty()) {
        return
    }

    Box(
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
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Text(
                text = "Truth per Episode",
                style = MaterialTheme.typography.titleMedium,
                color = Gold,
                fontWeight = FontWeight.Bold
            )
            stats.truthPerEpisode.forEach { ep ->
                val epName = stats.episodeNames[ep.episode.toString()] ?: "Episode ${ep.episode}"
                Column {
                    Text(
                        text = epName,
                        style = MaterialTheme.typography.bodyMedium,
                        color = TextPrimary,
                        fontWeight = FontWeight.SemiBold
                    )
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(16.dp)
                    ) {
                        Text(
                            text = "Red: ${ep.red}",
                            color = RedTruth,
                            style = MaterialTheme.typography.bodySmall
                        )
                        Text(
                            text = "Blue: ${ep.blue}",
                            color = BlueTruth,
                            style = MaterialTheme.typography.bodySmall
                        )
                    }
                    HorizontalDivider(
                        modifier = Modifier.padding(top = 4.dp),
                        color = PurpleMuted.copy(alpha = 0.5f)
                    )
                }
            }
        }
    }
}

@Composable
private fun InteractionsCard(stats: StatsData) {
    if (stats.interactions.isEmpty()) {
        return
    }

    Box(
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
    ) {
        Column(
            modifier = Modifier.padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(8.dp)
        ) {
            Text(
                text = "Top Interactions",
                style = MaterialTheme.typography.titleMedium,
                color = Gold,
                fontWeight = FontWeight.Bold
            )
            stats.interactions.take(10).forEachIndexed { index, interaction ->
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween
                ) {
                    Row(
                        horizontalArrangement = Arrangement.spacedBy(8.dp),
                        modifier = Modifier.weight(1f)
                    ) {
                        Text(
                            text = "${index + 1}.",
                            style = MaterialTheme.typography.bodyMedium,
                            color = Gold,
                            fontWeight = FontWeight.SemiBold
                        )
                        Text(
                            text = "${interaction.nameA} & ${interaction.nameB}",
                            style = MaterialTheme.typography.bodyMedium,
                            color = TextPrimary
                        )
                    }
                    Text(
                        text = "${interaction.count}",
                        style = MaterialTheme.typography.bodyMedium,
                        color = GoldDark
                    )
                }
            }
        }
    }
}
