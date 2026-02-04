package moe.auaurora.quotes.presentation.navigation

import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.navArgument
import moe.auaurora.quotes.data.repository.QuoteRepository
import moe.auaurora.quotes.presentation.browse.BrowseScreen
import moe.auaurora.quotes.presentation.browse.BrowseViewModel
import moe.auaurora.quotes.presentation.context.ContextScreen
import moe.auaurora.quotes.presentation.context.ContextViewModel
import moe.auaurora.quotes.presentation.detail.QuoteDetailScreen
import moe.auaurora.quotes.presentation.detail.QuoteDetailViewModel
import moe.auaurora.quotes.presentation.random.RandomQuoteScreen
import moe.auaurora.quotes.presentation.random.RandomQuoteViewModel
import moe.auaurora.quotes.presentation.search.SearchScreen
import moe.auaurora.quotes.presentation.search.SearchViewModel
import moe.auaurora.quotes.presentation.stats.StatsScreen
import moe.auaurora.quotes.presentation.stats.StatsViewModel
import moe.auaurora.quotes.util.AudioPlayerManager

@Composable
fun AppNavigation(
    navController: NavHostController,
    audioPlayer: AudioPlayerManager,
    repository: QuoteRepository,
    modifier: Modifier = Modifier
) {
    NavHost(
        modifier = modifier,
        navController = navController,
        startDestination = Screen.Search.route
    ) {
        composable(Screen.Search.route) {
            val viewModel: SearchViewModel = viewModel(
                factory = ViewModelFactory { SearchViewModel(repository) }
            )
            SearchScreen(
                viewModel = viewModel,
                audioPlayer = audioPlayer,
                onQuoteClick = { audioId ->
                    navController.navigate(Screen.QuoteDetail.createRoute(audioId))
                }
            )
        }

        composable(Screen.Browse.route) {
            val viewModel: BrowseViewModel = viewModel(
                factory = ViewModelFactory { BrowseViewModel(repository) }
            )
            BrowseScreen(
                viewModel = viewModel,
                audioPlayer = audioPlayer,
                onQuoteClick = { audioId ->
                    navController.navigate(Screen.QuoteDetail.createRoute(audioId))
                }
            )
        }

        composable(Screen.Random.route) {
            val viewModel: RandomQuoteViewModel = viewModel(
                factory = ViewModelFactory { RandomQuoteViewModel(repository) }
            )
            RandomQuoteScreen(
                viewModel = viewModel,
                audioPlayer = audioPlayer,
                onQuoteClick = { audioId ->
                    navController.navigate(Screen.QuoteDetail.createRoute(audioId))
                }
            )
        }

        composable(Screen.Stats.route) {
            val viewModel: StatsViewModel = viewModel(
                factory = ViewModelFactory { StatsViewModel(repository) }
            )
            StatsScreen(viewModel = viewModel)
        }

        composable(
            route = Screen.QuoteDetail.route,
            arguments = listOf(navArgument("audioId") { type = NavType.StringType })
        ) { backStackEntry ->
            val audioId = backStackEntry.arguments?.getString("audioId") ?: return@composable
            val viewModel: QuoteDetailViewModel = viewModel(
                factory = ViewModelFactory { QuoteDetailViewModel(repository, audioId) }
            )
            QuoteDetailScreen(
                viewModel = viewModel,
                audioPlayer = audioPlayer,
                onViewContext = { id ->
                    navController.navigate(Screen.Context.createRoute(id))
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.Context.route,
            arguments = listOf(navArgument("audioId") { type = NavType.StringType })
        ) { backStackEntry ->
            val audioId = backStackEntry.arguments?.getString("audioId") ?: return@composable
            val viewModel: ContextViewModel = viewModel(
                factory = ViewModelFactory { ContextViewModel(repository, audioId) }
            )
            ContextScreen(
                viewModel = viewModel,
                audioPlayer = audioPlayer,
                onQuoteClick = { audioId ->
                    navController.navigate(Screen.QuoteDetail.createRoute(audioId))
                },
                onBack = { navController.popBackStack() }
            )
        }
    }
}
