package moe.auaurora.quotes.data.remote

import moe.auaurora.quotes.data.remote.dto.*
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Path
import retrofit2.http.Query

interface UminekoApi {

    @GET("/api/v1/search")
    suspend fun search(
        @Query("q") query: String,
        @Query("lang") lang: String = "en",
        @Query("limit") limit: Int = 30,
        @Query("offset") offset: Int = 0,
        @Query("character") character: String? = null,
        @Query("episode") episode: Int? = null,
        @Query("truth") truth: String? = null
    ): Response<SearchResponseDto>

    @GET("/api/v1/random")
    suspend fun random(
        @Query("lang") lang: String = "en",
        @Query("character") character: String? = null,
        @Query("episode") episode: Int? = null,
        @Query("truth") truth: String? = null
    ): Response<QuoteDto>

    @GET("/api/v1/browse")
    suspend fun browse(
        @Query("lang") lang: String = "en",
        @Query("limit") limit: Int = 30,
        @Query("offset") offset: Int = 0,
        @Query("character") character: String? = null,
        @Query("episode") episode: Int? = null,
        @Query("truth") truth: String? = null
    ): Response<BrowseResponseDto>

    @GET("/api/v1/quote/{audioId}")
    suspend fun getQuote(
        @Path("audioId") audioId: String,
        @Query("lang") lang: String = "en"
    ): Response<QuoteDto>

    @GET("/api/v1/context/{audioId}")
    suspend fun getContext(
        @Path("audioId") audioId: String,
        @Query("lang") lang: String = "en",
        @Query("lines") lines: Int = 5
    ): Response<ContextResponseDto>

    @GET("/api/v1/characters")
    suspend fun getCharacters(): Response<Map<String, String>>

    @GET("/api/v1/stats")
    suspend fun getStats(
        @Query("episode") episode: Int? = null
    ): Response<StatsResponseDto>

    @GET("/api/v1/config")
    suspend fun getConfig(): Response<ConfigDto>
}
