package moe.auaurora.quotes.presentation.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.automirrored.filled.ArrowForward
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.ui.theme.BgCard
import moe.auaurora.quotes.ui.theme.Gold
import moe.auaurora.quotes.ui.theme.PurpleMuted
import moe.auaurora.quotes.ui.theme.TextMuted

@Composable
fun PaginationControls(
    offset: Int,
    limit: Int,
    total: Int,
    onPrevious: () -> Unit,
    onNext: () -> Unit,
    modifier: Modifier = Modifier
) {
    if (total <= 0) {
        return
    }

    val currentPage = (offset / limit) + 1
    val totalPages = ((total - 1) / limit) + 1

    Row(
        modifier = modifier
            .fillMaxWidth()
            .background(BgCard)
            .border(1.dp, PurpleMuted)
            .padding(horizontal = 8.dp, vertical = 4.dp),
        horizontalArrangement = Arrangement.Center,
        verticalAlignment = Alignment.CenterVertically
    ) {
        IconButton(
            onClick = onPrevious,
            enabled = offset > 0
        ) {
            Icon(
                Icons.AutoMirrored.Filled.ArrowBack,
                contentDescription = "Previous page",
                tint = if (offset > 0) Gold else PurpleMuted
            )
        }

        Text(
            text = "$currentPage / $totalPages",
            style = MaterialTheme.typography.labelLarge,
            color = Gold
        )
        Text(
            text = "  ($total results)",
            style = MaterialTheme.typography.labelSmall,
            color = TextMuted
        )

        IconButton(
            onClick = onNext,
            enabled = offset + limit < total
        ) {
            Icon(
                Icons.AutoMirrored.Filled.ArrowForward,
                contentDescription = "Next page",
                tint = if (offset + limit < total) Gold else PurpleMuted
            )
        }
    }
}
