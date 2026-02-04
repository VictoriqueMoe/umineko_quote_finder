package moe.auaurora.quotes.presentation.navigation

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.List
import androidx.compose.material.icons.filled.Casino
import androidx.compose.material.icons.filled.QueryStats
import androidx.compose.material.icons.filled.Search
import androidx.compose.ui.graphics.vector.ImageVector

sealed class Screen(val route: String) {
    data object Search : Screen("search")
    data object Browse : Screen("browse")
    data object Random : Screen("random")
    data object Stats : Screen("stats")

    data object QuoteDetail : Screen("quote/{audioId}") {
        fun createRoute(audioId: String) = "quote/$audioId"
    }

    data object Context : Screen("context/{audioId}") {
        fun createRoute(audioId: String) = "context/$audioId"
    }
}

data class BottomNavItem(
    val screen: Screen,
    val label: String,
    val icon: ImageVector
)

val bottomNavItems = listOf(
    BottomNavItem(Screen.Search, "Search", Icons.Filled.Search),
    BottomNavItem(Screen.Browse, "Browse", Icons.AutoMirrored.Filled.List),
    BottomNavItem(Screen.Random, "Random", Icons.Filled.Casino),
    BottomNavItem(Screen.Stats, "Stats", Icons.Filled.QueryStats),
)
