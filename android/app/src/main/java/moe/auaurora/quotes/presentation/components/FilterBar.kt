package moe.auaurora.quotes.presentation.components

import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.unit.dp
import moe.auaurora.quotes.domain.model.FilterState
import moe.auaurora.quotes.ui.theme.*

@OptIn(ExperimentalLayoutApi::class, ExperimentalMaterial3Api::class)
@Composable
fun FilterBar(
    filters: FilterState,
    characters: Map<String, String>,
    onFiltersChanged: (FilterState) -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxWidth()
            .background(BgCard.copy(alpha = 0.5f))
            .border(1.dp, PurpleMuted)
            .padding(12.dp),
        verticalArrangement = Arrangement.spacedBy(10.dp)
    ) {
        // Language and truth filters
        FlowRow(
            horizontalArrangement = Arrangement.spacedBy(6.dp),
        ) {
            GothicChip(
                text = "English",
                selected = filters.language == "en",
                onClick = { onFiltersChanged(filters.copy(language = "en")) }
            )
            GothicChip(
                text = "Japanese",
                selected = filters.language == "ja",
                onClick = { onFiltersChanged(filters.copy(language = "ja")) }
            )
            GothicChip(
                text = "All",
                selected = filters.truth == "",
                onClick = { onFiltersChanged(filters.copy(truth = "")) }
            )
            GothicChip(
                text = "Red",
                selected = filters.truth == "red",
                selectedColour = RedTruth,
                onClick = { onFiltersChanged(filters.copy(truth = "red")) }
            )
            GothicChip(
                text = "Blue",
                selected = filters.truth == "blue",
                selectedColour = BlueTruth,
                onClick = { onFiltersChanged(filters.copy(truth = "blue")) }
            )
        }

        // Episode filters
        FlowRow(
            horizontalArrangement = Arrangement.spacedBy(6.dp),
        ) {
            GothicChip(
                text = "All EPs",
                selected = filters.episode == 0,
                onClick = { onFiltersChanged(filters.copy(episode = 0)) }
            )
            for (ep in 1..8) {
                GothicChip(
                    text = "EP$ep",
                    selected = filters.episode == ep,
                    onClick = { onFiltersChanged(filters.copy(episode = ep)) }
                )
            }
        }

        // Character dropdown
        if (characters.isNotEmpty()) {
            var expanded by remember { mutableStateOf(false) }
            val selectedName = if (filters.character.isEmpty()) {
                "All Characters"
            } else {
                characters[filters.character] ?: filters.character
            }

            ExposedDropdownMenuBox(
                expanded = expanded,
                onExpandedChange = { expanded = it }
            ) {
                OutlinedTextField(
                    value = selectedName,
                    onValueChange = {},
                    readOnly = true,
                    trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = expanded) },
                    modifier = Modifier
                        .fillMaxWidth()
                        .menuAnchor(ExposedDropdownMenuAnchorType.PrimaryNotEditable),
                    textStyle = MaterialTheme.typography.bodySmall.copy(color = TextPrimary),
                    colors = OutlinedTextFieldDefaults.colors(
                        unfocusedBorderColor = PurpleMuted,
                        focusedBorderColor = Gold,
                        unfocusedContainerColor = BgCard,
                        focusedContainerColor = BgCard,
                        cursorColor = Gold,
                        unfocusedTrailingIconColor = GoldDark,
                        focusedTrailingIconColor = Gold
                    )
                )
                ExposedDropdownMenu(
                    expanded = expanded,
                    onDismissRequest = { expanded = false }
                ) {
                    DropdownMenuItem(
                        text = {
                            Text(
                                "All Characters",
                                style = MaterialTheme.typography.bodySmall
                            )
                        },
                        onClick = {
                            onFiltersChanged(filters.copy(character = ""))
                            expanded = false
                        }
                    )
                    characters.entries.sortedBy { it.value }.forEach { (id, name) ->
                        DropdownMenuItem(
                            text = {
                                Text(name, style = MaterialTheme.typography.bodySmall)
                            },
                            onClick = {
                                onFiltersChanged(filters.copy(character = id))
                                expanded = false
                            }
                        )
                    }
                }
            }
        }
    }
}

@Composable
private fun GothicChip(
    text: String,
    selected: Boolean,
    selectedColour: androidx.compose.ui.graphics.Color = Gold,
    onClick: () -> Unit
) {
    val borderColour = if (selected) selectedColour else PurpleMuted
    val bgColour = if (selected) selectedColour.copy(alpha = 0.15f) else BgCard.copy(alpha = 0.3f)
    val textColour = if (selected) selectedColour else TextMuted

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
