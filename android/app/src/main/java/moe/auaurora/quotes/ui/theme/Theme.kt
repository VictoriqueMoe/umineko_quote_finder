package moe.auaurora.quotes.ui.theme

import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.runtime.Composable

private val UminekoColorScheme = darkColorScheme(
    primary = Gold,
    onPrimary = BgVoid,
    primaryContainer = GoldDark,
    onPrimaryContainer = GoldLight,
    secondary = Purple,
    onSecondary = TextPrimary,
    secondaryContainer = PurpleMuted,
    onSecondaryContainer = PurpleLight,
    tertiary = Rose,
    onTertiary = TextPrimary,
    background = BgVoid,
    onBackground = TextPrimary,
    surface = BgPurple,
    onSurface = TextPrimary,
    surfaceVariant = BgCard,
    onSurfaceVariant = TextMuted,
    outline = PurpleMuted,
    outlineVariant = PurpleMuted,
    error = RedTruth,
    onError = TextPrimary,
    surfaceContainerLowest = BgVoid,
    surfaceContainerLow = BgVoid,
    surfaceContainer = BgPurple,
    surfaceContainerHigh = BgCard,
    surfaceContainerHighest = BgCardHover,
)

@Composable
fun UminekoQuotesTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = UminekoColorScheme,
        typography = UminekoTypography,
        content = content
    )
}
