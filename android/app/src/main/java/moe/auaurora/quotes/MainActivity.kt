package moe.auaurora.quotes

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.SystemBarStyle
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.DisposableEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.presentation.navigation.AppNavigation
import moe.auaurora.quotes.presentation.navigation.bottomNavItems
import moe.auaurora.quotes.ui.theme.BgVoid
import moe.auaurora.quotes.ui.theme.Gold
import moe.auaurora.quotes.ui.theme.TextMuted
import moe.auaurora.quotes.ui.theme.UminekoQuotesTheme
import moe.auaurora.quotes.util.AudioPlayerManager

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        val darkBarStyle = SystemBarStyle.dark(android.graphics.Color.parseColor("#0A0612"))
        enableEdgeToEdge(
            statusBarStyle = darkBarStyle,
            navigationBarStyle = darkBarStyle
        )
        super.onCreate(savedInstanceState)
        setContent {
            UminekoQuotesTheme {
                UminekoQuotesApp()
            }
        }
    }
}

@Composable
fun UminekoQuotesApp() {
    val navController = rememberNavController()
    val context = LocalContext.current
    val audioPlayer = remember { AudioPlayerManager(context) }
    val repository = remember { QuoteRepository() }

    DisposableEffect(Unit) {
        onDispose { audioPlayer.release() }
    }

    Scaffold(
        bottomBar = {
            val navBackStackEntry by navController.currentBackStackEntryAsState()
            val currentDestination = navBackStackEntry?.destination

            // Only show bottom bar on top-level screens
            val showBottomBar = bottomNavItems.any { item ->
                currentDestination?.hierarchy?.any { it.route == item.screen.route } == true
            }

            if (showBottomBar) {
                NavigationBar(
                    containerColor = BgVoid,
                    contentColor = Gold
                ) {
                    bottomNavItems.forEach { item ->
                        val selected = currentDestination?.hierarchy?.any {
                            it.route == item.screen.route
                        } == true

                        NavigationBarItem(
                            icon = { Icon(item.icon, contentDescription = item.label) },
                            label = {
                                androidx.compose.material3.Text(
                                    text = item.label,
                                    fontWeight = if (selected) {
                                        FontWeight.SemiBold
                                    } else {
                                        FontWeight.Normal
                                    },
                                    fontSize = 11.sp
                                )
                            },
                            selected = selected,
                            onClick = {
                                navController.navigate(item.screen.route) {
                                    popUpTo(navController.graph.findStartDestination().id) {
                                        saveState = true
                                    }
                                    launchSingleTop = true
                                    restoreState = true
                                }
                            },
                            colors = NavigationBarItemDefaults.colors(
                                selectedIconColor = Gold,
                                selectedTextColor = Gold,
                                unselectedIconColor = TextMuted,
                                unselectedTextColor = TextMuted,
                                indicatorColor = Color.Transparent
                            )
                        )
                    }
                }
            }
        }
    ) { innerPadding ->
        AppNavigation(
            navController = navController,
            audioPlayer = audioPlayer,
            repository = repository,
            modifier = Modifier.padding(innerPadding)
        )
    }
}
