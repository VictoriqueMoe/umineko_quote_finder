package moe.auaurora.quotes.util

import android.content.Context
import androidx.media3.common.MediaItem
import androidx.media3.common.Player
import androidx.media3.exoplayer.ExoPlayer
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import moe.auaurora.quotes.BuildConfig

class AudioPlayerManager(context: Context) {

    private var player: ExoPlayer? = ExoPlayer.Builder(context).build()

    private val _isPlaying = MutableStateFlow(false)
    val isPlaying: StateFlow<Boolean> = _isPlaying.asStateFlow()

    private val _currentAudioId = MutableStateFlow<String?>(null)
    val currentAudioId: StateFlow<String?> = _currentAudioId.asStateFlow()

    init {
        player?.addListener(object : Player.Listener {
            override fun onIsPlayingChanged(playing: Boolean) {
                _isPlaying.value = playing
            }

            override fun onPlaybackStateChanged(playbackState: Int) {
                if (playbackState == Player.STATE_ENDED) {
                    _isPlaying.value = false
                    _currentAudioId.value = null
                }
            }
        })
    }

    fun playSingle(charId: String, audioId: String) {
        val ids = audioId.split(",").map { it.trim() }.filter { it.isNotEmpty() }
        if (ids.isEmpty()) {
            return
        }

        _currentAudioId.value = ids.first()

        val url = if (ids.size == 1) {
            "${BuildConfig.BASE_URL}/api/v1/audio/$charId/${ids.first()}"
        } else {
            "${BuildConfig.BASE_URL}/api/v1/audio/$charId/combined?ids=${ids.joinToString(",")}"
        }
        play(url)
    }

    fun playCombined(charId: String, audioIds: List<String>) {
        val url = "${BuildConfig.BASE_URL}/api/v1/audio/$charId/combined?ids=${audioIds.joinToString(",")}"
        _currentAudioId.value = audioIds.firstOrNull()
        play(url)
    }

    private fun play(url: String) {
        player?.apply {
            stop()
            setMediaItem(MediaItem.fromUri(url))
            prepare()
            playWhenReady = true
        }
    }

    fun pause() {
        player?.pause()
    }

    fun resume() {
        player?.play()
    }

    fun stop() {
        player?.stop()
        _isPlaying.value = false
        _currentAudioId.value = null
    }

    fun release() {
        player?.release()
        player = null
        _isPlaying.value = false
        _currentAudioId.value = null
    }
}
