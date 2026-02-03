(async () => {
    const API_BASE = '/api/v1';

    const searchInput = document.getElementById('searchInput');
    const searchBtn = document.getElementById('searchBtn');
    const audioIdInput = document.getElementById('audioIdInput');
    const audioIdBtn = document.getElementById('audioIdBtn');
    const randomBtn = document.getElementById('randomBtn');
    const clearBtn = document.getElementById('clearBtn');
    const statsBtn = document.getElementById('statsBtn');
    const characterSelect = document.getElementById('characterSelect');
    const episodeSelect = document.getElementById('episodeSelect');
    const browseBtn = document.getElementById('browseBtn');
    const truthSelect = document.getElementById('truthSelect');
    const resultsContainer = document.getElementById('resultsContainer');
    const langBtns = document.querySelectorAll('.lang-btn');

    let currentLang = 'en';
    let currentAudioId = null;
    let browseMode = false;
    let statsMode = false;
    let quoteMode = false;
    let statsCache = {};
    let statsCharts = [];
    let browseCharacter = '';
    let browseEpisode = 0;
    let browseOffset = 0;
    let browseTotal = 0;
    let hasAudio = true;

    const contentTypeLabels = { tea: 'Tea Party', ura: '????', omake: 'Omake' };

    function episodeLabel(quote) {
        if (!quote.episode) {
            return '';
        }
        let label = `Episode ${quote.episode}`;
        if (quote.contentType && contentTypeLabels[quote.contentType]) {
            label += ` \u2014 ${contentTypeLabels[quote.contentType]}`;
        }
        return label;
    }

    function createButterflies() {
        const container = document.getElementById('butterflies');
        const count = 8;

        for (let i = 0; i < count; i++) {
            const butterfly = document.createElement('div');
            butterfly.className = 'butterfly';
            butterfly.innerHTML = '🦋';
            butterfly.style.setProperty('--start-x', `${Math.random() * 100}vw`);
            butterfly.style.setProperty('--duration', `${20 + Math.random() * 15}s`);
            butterfly.style.setProperty('--delay', `${Math.random() * 20}s`);
            butterfly.style.left = `${Math.random() * 100}%`;
            container.appendChild(butterfly);
        }
    }

    async function loadCharacters() {
        try {
            const response = await fetch(`${API_BASE}/characters`);
            const characters = await response.json();

            const sorted = Object.entries(characters).sort((a, b) => a[1].localeCompare(b[1]));

            for (const [id, name] of sorted) {
                const option = document.createElement('option');
                option.value = id;
                option.textContent = name;
                characterSelect.appendChild(option);
            }
        } catch (error) {
            console.error('Failed to load characters:', error);
        }
    }

    function showLoading() {
        resultsContainer.innerHTML = `
            <div class="loading">
                <div class="loading-spinner"></div>
                <p class="loading-text">Searching through the fragments...</p>
            </div>
        `;
    }

    function showEmpty(message = 'No quotes found in this fragment.') {
        resultsContainer.innerHTML = `
            <div class="empty-state">
                <div class="empty-icon">🦋</div>
                <h3 class="empty-title">The Golden Land remains silent</h3>
                <p class="empty-subtitle">${message}</p>
            </div>
        `;
    }

    let currentQuery = '';
    let currentOffset = 0;
    let currentTotal = 0;
    const PAGE_SIZE = 30;

    function renderQuotes(results, query, total, offset) {
        stopAudio();
        currentAudioId = null;
        if (!results || results.length === 0) {
            showEmpty();
            return;
        }

        const start = offset + 1;
        const end = offset + results.length;

        const header = query ? `
            <div class="results-header">
                <span class="results-count">Showing <span>${start}-${end}</span> of <span>${total}</span> fragments for "${escapeHtml(query)}"</span>
            </div>
        ` : '';

        const quotes = results.map((item, index) => {
            const quote = item.quote || item;
            return `
                <article class="quote-card" style="--index: ${index}">
                    <span class="quote-mark">"</span>
                    <p class="quote-text">${quote.textHtml || escapeHtml(quote.text)}</p>
                    <div class="quote-meta">
                        <span class="quote-character">— ${escapeHtml(quote.character)}</span>
                        <div class="quote-details">
                            ${quote.episode ? `<span class="quote-episode">${episodeLabel(quote)}</span>` : ''}

                            ${langToggleHTML(quote.audioId)}
                        </div>
                    </div>
                    ${audioPlayerHTML(quote.audioId, quote.characterId)}
                    <div class="quote-actions">
                        ${contextBtnHTML(quote.audioId)}
                        ${shareBtnHTML(quote.audioId)}
                    </div>
                </article>
            `;
        }).join('');

        const totalPages = Math.ceil(total / PAGE_SIZE);
        const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

        const pagination = total > PAGE_SIZE ? `
            <div class="pagination">
                <button class="pagination-btn" id="prevPage" ${offset === 0 ? 'disabled' : ''}>◀ Previous</button>
                <span class="pagination-info">Page <span>${currentPage}</span> of <span>${totalPages}</span></span>
                <button class="pagination-btn" id="nextPage" ${end >= total ? 'disabled' : ''}>Next ▶</button>
            </div>
        ` : '';

        resultsContainer.innerHTML = header + `<div class="quotes-grid">${quotes}</div>` + pagination;

        if (total > PAGE_SIZE) {
            document.getElementById('prevPage')?.addEventListener('click', () => {
                if (currentOffset >= PAGE_SIZE) {
                    currentOffset -= PAGE_SIZE;
                    search(currentQuery, currentOffset);
                }
            });
            document.getElementById('nextPage')?.addEventListener('click', () => {
                if (currentOffset + PAGE_SIZE < currentTotal) {
                    currentOffset += PAGE_SIZE;
                    search(currentQuery, currentOffset);
                }
            });
        }
    }

    function renderFeaturedQuote(quote) {
        stopAudio();
        currentAudioId = quote.audioId ? quote.audioId.split(', ')[0] : null;
        resultsContainer.innerHTML = `
            <article class="featured-quote">
                <div class="featured-label">✦ A Fragment from the Sea ✦</div>
                <p class="featured-text">"${quote.textHtml || escapeHtml(quote.text)}"</p>
                <p class="featured-character">— ${escapeHtml(quote.character)}</p>
                ${quote.episode ? `<p class="featured-episode">${episodeLabel(quote)}</p>` : ''}
                ${audioPlayerHTML(quote.audioId, quote.characterId)}
                ${langToggleHTML(quote.audioId)}
                <div class="quote-actions">
                    ${contextBtnHTML(quote.audioId)}
                    ${shareBtnHTML(quote.audioId)}
                </div>
            </article>
        `;
    }

    function formatAudioIds(audioId) {
        return audioId.split(', ').map(id => id + '.ogg').join(', ');
    }

    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    function langToggleHTML(audioId) {
        if (!audioId) {
            return '';
        }
        const firstId = escapeHtml(audioId.split(', ')[0]);
        const enActive = currentLang === 'en' ? ' active' : '';
        const jaActive = currentLang === 'ja' ? ' active' : '';
        return `<span class="lang-card-toggle" data-audio-id="${firstId}"><button class="lang-card-btn${enActive}" data-lang="en">EN</button><button class="lang-card-btn${jaActive}" data-lang="ja">JA</button></span>`;
    }

    function shareBtnHTML(audioId) {
        if (!audioId) {
            return '';
        }
        const firstId = escapeHtml(audioId.split(', ')[0]);
        return `<button class="share-btn" data-audio-id="${firstId}">Share this Fragment</button>`;
    }

    function contextBtnHTML(audioId) {
        if (!audioId) {
            return '';
        }
        const firstId = escapeHtml(audioId.split(', ')[0]);
        return `<button class="context-btn" data-audio-id="${firstId}">Show Context</button>`;
    }

    function audioPlayerHTML(audioId, characterId) {
        if (!audioId || !hasAudio) {
            return '';
        }
        const charId = escapeHtml(characterId || '');
        const ids = audioId.split(', ');
        let individualClips = '';
        for (let i = 0; i < ids.length; i++) {
            const eid = escapeHtml(ids[i]);
            individualClips += `<button class="audio-clip-btn" data-audio-id="${eid}" data-char-id="${charId}">&#9654; ${eid}.ogg</button>`;
        }
        let clipsHTML;
        if (ids.length > 1) {
            const allIds = ids.map(id => escapeHtml(id)).join(',');
            clipsHTML = `<div class="audio-clips"><button class="audio-clip-btn audio-combined-btn" data-char-id="${charId}" data-audio-ids="${allIds}">&#9654; Combined (${ids.length} clips)</button><button class="audio-expand-btn">&#9662; Individual</button></div><div class="audio-individual-clips">${individualClips}</div>`;
        } else {
            clipsHTML = `<div class="audio-clips">${individualClips}</div>`;
        }
        const savedVolume = localStorage.getItem('uminekoVolume') ?? "1.0";
        const volume = Math.floor(parseFloat(savedVolume) * 100);
        return `<div class="audio-player">
            ${clipsHTML}
            <div class="audio-controls">
                <div class="audio-track"><div class="audio-progress"></div></div>
                <div class="audio-volume"><span>VOL</span><input class="audio-volume-slider" type="range" min="0" max="100" step="1" value="${volume}"></div>
                <span class="audio-time">0:00 / 0:00</span>
            </div>
        </div>`;
    }

    let activeAudio = null;
    let activeBtn = null;
    let activePlayer = null;

    function formatTime(sec) {
        if (!sec || !isFinite(sec)) {
            return '0:00';
        }
        const m = Math.floor(sec / 60);
        const s = Math.floor(sec % 60);
        return m + ':' + (s < 10 ? '0' : '') + s;
    }

    function stopAudio() {
        if (activeAudio) {
            activeAudio.pause();
            activeAudio.removeAttribute('src');
            activeAudio.load();
        }
        if (activeBtn) {
            activeBtn.classList.remove('active');
        }
        if (activePlayer) {
            activePlayer.classList.remove('playing');
            const controls = activePlayer.querySelector('.audio-controls');
            if (controls) {
                controls.classList.remove('visible');
            }
            const progress = activePlayer.querySelector('.audio-progress');
            if (progress) {
                progress.style.width = '0%';
            }
            const timeEl = activePlayer.querySelector('.audio-time');
            if (timeEl) {
                timeEl.textContent = '0:00 / 0:00';
            }
        }
        activeBtn = null;
        activePlayer = null;
    }

    function updateSliderFill(slider) {
        const pct = slider.value;
        slider.style.background = `linear-gradient(to right, #d4a84b 0%, #d4a84b ${pct}%, #3d2a5c ${pct}%, #3d2a5c 100%)`;
    }

    function setVolume(volume) {
        activeAudio.volume = volume;
        localStorage.setItem('uminekoVolume', volume.toString());
    }

    function playAudioUrl(url, btn) {
        const player = btn.closest('.audio-player');
        if (activeAudio && activeBtn === btn) {
            if (activeAudio.paused) {
                activeAudio.play();
                player.classList.add('playing');
            } else {
                activeAudio.pause();
                player.classList.remove('playing');
            }
            return;
        }

        stopAudio();

        if (!activeAudio) {
            activeAudio = new Audio();
            const savedVolume = localStorage.getItem('uminekoVolume');
            if (savedVolume) {
                activeAudio.volume = parseFloat(savedVolume);
            }
            activeAudio.addEventListener('timeupdate', () => {
                const progress = activePlayer.querySelector('.audio-progress');
                const timeEl = activePlayer.querySelector('.audio-time');
                if (progress && activeAudio.duration) {
                    progress.style.width = (activeAudio.currentTime / activeAudio.duration * 100) + '%';
                }
                if (timeEl) {
                    timeEl.textContent = formatTime(activeAudio.currentTime) + ' / ' + formatTime(activeAudio.duration);
                }
            });
            activeAudio.addEventListener('ended', () => {
                stopAudio();
            });
            activeAudio.addEventListener('loadedmetadata', () => {
                if (activePlayer) {
                    const timeEl = activePlayer.querySelector('.audio-time');
                    if (timeEl) {
                        timeEl.textContent = '0:00 / ' + formatTime(activeAudio.duration);
                    }
                }
            });
        }

        activeBtn = btn;
        activePlayer = player;
        btn.classList.add('active');
        const volumeSlider = player.querySelector('.audio-volume-slider');
        if (volumeSlider) {
            updateSliderFill(volumeSlider);
            volumeSlider.oninput = () => {
                const v = parseFloat(volumeSlider.value) / 100;
                setVolume(v);
                updateSliderFill(volumeSlider);
            };
        }
        const controls = player.querySelector('.audio-controls');
        if (controls) {
            controls.classList.add('visible');
        }
        activeAudio.src = url;
        activeAudio.play();
        player.classList.add('playing');
    }

    function playAudio(charId, audioId, btn) {
        playAudioUrl(`${API_BASE}/audio/${charId}/${audioId}`, btn);
    }

    function renderContextSection(data) {
        const lines = [...data.before, data.quote, ...data.after];
        const quoteAudioId = data.quote.audioId || '';
        return lines.map(line => {
            const isHighlight = line.audioId === quoteAudioId && quoteAudioId !== '';
            const highlightClass = isHighlight ? ' context-highlight' : '';
            const firstId = line.audioId ? escapeHtml(line.audioId.split(', ')[0]) : '';
            const clickable = firstId && !isHighlight ? ` context-clickable" data-audio-id="${firstId}` : '';
            return `<div class="context-line${highlightClass}${clickable}">
                <span class="context-character">${escapeHtml(line.character)}</span>
                <span class="context-text">${line.textHtml || escapeHtml(line.text)}</span>
            </div>`;
        }).join('');
    }

    resultsContainer.addEventListener('click', async (e) => {
        const contextBtn = e.target.closest('.context-btn');
        if (contextBtn) {
            const card = contextBtn.closest('.quote-card') || contextBtn.closest('.featured-quote');
            const existing = card.querySelector('.context-section');
            if (existing) {
                existing.remove();
                contextBtn.textContent = 'Show Context';
                return;
            }
            const audioId = contextBtn.dataset.audioId;
            const activeLangBtn = card ? card.querySelector('.lang-card-btn.active') : null;
            const lang = activeLangBtn ? activeLangBtn.dataset.lang : currentLang;
            contextBtn.textContent = 'Loading...';
            contextBtn.disabled = true;
            try {
                const response = await fetch(`${API_BASE}/context/${audioId}?lang=${lang}&lines=5`);
                const data = await response.json();
                if (data.error) {
                    contextBtn.textContent = 'Show Context';
                    contextBtn.disabled = false;
                    return;
                }
                const section = document.createElement('div');
                section.className = 'context-section';
                section.innerHTML = renderContextSection(data);
                card.appendChild(section);
                contextBtn.textContent = 'Hide Context';
            } catch (err) {
                console.error('Failed to load context:', err);
                contextBtn.textContent = 'Show Context';
            } finally {
                contextBtn.disabled = false;
            }
            return;
        }
        const contextLine = e.target.closest('.context-clickable');
        if (contextLine) {
            const audioId = contextLine.dataset.audioId;
            if (audioId) {
                quoteMode = true;
                getQuoteByAudioId(audioId);
            }
            return;
        }
        const shareBtn = e.target.closest('.share-btn');
        if (shareBtn) {
            const audioId = shareBtn.dataset.audioId;
            const card = shareBtn.closest('.featured-quote') || shareBtn.closest('.quote-card');
            const activeLangBtn = card ? card.querySelector('.lang-card-btn.active') : null;
            const lang = activeLangBtn ? activeLangBtn.dataset.lang : currentLang;
            let url = window.location.origin + '/?quote=' + audioId;
            if (lang !== 'en') {
                url += '&lang=' + lang;
            }
            navigator.clipboard.writeText(url).then(() => {
                shareBtn.textContent = 'Link Copied';
                setTimeout(() => {
                    shareBtn.textContent = 'Share this Fragment';
                }, 2000);
            });
            return;
        }
        const expandBtn = e.target.closest('.audio-expand-btn');
        if (expandBtn) {
            const player = expandBtn.closest('.audio-player');
            const individual = player.querySelector('.audio-individual-clips');
            if (individual) {
                const open = individual.classList.toggle('visible');
                expandBtn.innerHTML = open ? '&#9652; Individual' : '&#9662; Individual';
            }
            return;
        }
        const combinedBtn = e.target.closest('.audio-combined-btn');
        if (combinedBtn) {
            const charId = combinedBtn.dataset.charId;
            const audioIds = combinedBtn.dataset.audioIds;
            playAudioUrl(`${API_BASE}/audio/${charId}/combined?ids=${audioIds}`, combinedBtn);
            return;
        }
        const clipBtn = e.target.closest('.audio-clip-btn');
        if (clipBtn) {
            playAudio(clipBtn.dataset.charId, clipBtn.dataset.audioId, clipBtn);
            return;
        }
        const track = e.target.closest('.audio-track');
        if (track && activeAudio && activeAudio.duration) {
            const player = track.closest('.audio-player');
            if (player === activePlayer) {
                const rect = track.getBoundingClientRect();
                const ratio = (e.clientX - rect.left) / rect.width;
                activeAudio.currentTime = ratio * activeAudio.duration;
            }
        }
    });

    async function search(query, offset = 0) {
        if (!query.trim()) {
            showEmpty('Enter a search term to find quotes.');
            return;
        }

        browseMode = false;
        statsMode = false;
        quoteMode = false;
        destroyStatsCharts();
        showLoading();

        try {
            const characterId = characterSelect.value;
            const episode = episodeSelect.value;
            const truth = truthSelect.value;
            let url = `${API_BASE}/search?q=${encodeURIComponent(query)}&limit=${PAGE_SIZE}&offset=${offset}&lang=${currentLang}`;
            if (characterId) {
                url += `&character=${characterId}`;
            }
            if (episode && episode !== '0') {
                url += `&episode=${episode}`;
            }
            if (truth) {
                url += `&truth=${truth}`;
            }

            const response = await fetch(url);
            const data = await response.json();

            currentQuery = query;
            currentOffset = data.offset;
            currentTotal = data.total;

            const results = data.results || [];

            renderQuotes(results, query, data.total, data.offset);
            updateURL();
        } catch (error) {
            console.error('Search failed:', error);
            showEmpty('Failed to search. Please try again.');
        }
    }

    async function getRandomQuote() {
        browseMode = false;
        statsMode = false;
        quoteMode = false;
        destroyStatsCharts();
        showLoading();

        try {
            const characterId = characterSelect.value;
            const episode = episodeSelect.value;
            const truth = truthSelect.value;
            let url = `${API_BASE}/random?lang=${currentLang}`;
            if (characterId) {
                url += `&character=${characterId}`;
            }
            if (episode && episode !== '0') {
                url += `&episode=${episode}`;
            }
            if (truth) {
                url += `&truth=${truth}`;
            }
            const response = await fetch(url);
            const quote = await response.json();
            if (quote.error) {
                showEmpty('No quotes found for this character.');
                return;
            }
            renderFeaturedQuote(quote);
            updateURL();
        } catch (error) {
            console.error('Failed to get random quote:', error);
            showEmpty('Failed to retrieve a quote.');
        }
    }

    async function getQuoteByAudioId(audioId) {
        browseMode = false;
        statsMode = false;
        destroyStatsCharts();
        showLoading();

        try {
            const response = await fetch(`${API_BASE}/quote/${audioId}?lang=${currentLang}`);
            const quote = await response.json();
            if (quote.error) {
                showEmpty(`No quote found for audio ID "${escapeHtml(audioId)}".`);
                return;
            }
            renderFeaturedQuote(quote);
            updateURL();
        } catch (error) {
            console.error('Failed to get quote by audioId:', error);
            showEmpty('Failed to look up audio ID. Please try again.');
        }
    }

    searchBtn.addEventListener('click', () => {
        currentOffset = 0;
        search(searchInput.value, 0);
    });
    searchInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            currentOffset = 0;
            search(searchInput.value, 0);
        }
    });

    function lookupAudioId() {
        const id = audioIdInput.value.trim();
        if (!id) {
            return;
        }
        getQuoteByAudioId(id);
    }
    audioIdBtn.addEventListener('click', lookupAudioId);
    audioIdInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            lookupAudioId();
        }
    });

    randomBtn.addEventListener('click', getRandomQuote);

    async function loadStats() {
        browseMode = false;
        quoteMode = false;
        showLoading();
        try {
            const ep = parseInt(episodeSelect.value) || 0;
            const cacheKey = 'ep' + ep;
            if (!statsCache[cacheKey]) {
                let url = `${API_BASE}/stats`;
                if (ep > 0) {
                    url += `?episode=${ep}`;
                }
                const response = await fetch(url);
                statsCache[cacheKey] = await response.json();
            }
            statsMode = true;
            renderStats(statsCache[cacheKey]);
            updateURL();
        } catch (error) {
            console.error('Failed to load stats:', error);
            showEmpty('Failed to load statistics.');
        }
    }

    statsBtn.addEventListener('click', loadStats);

    clearBtn.addEventListener('click', () => {
        stopAudio();
        destroyStatsCharts();
        searchInput.value = '';
        audioIdInput.value = '';
        characterSelect.value = '';
        episodeSelect.value = '0';
        truthSelect.value = '';
        resultsContainer.innerHTML = '';
        browseMode = false;
        statsMode = false;
        quoteMode = false;
        updateBrowseBtn();
        updateURL();
    });

    function updateBrowseBtn() {
        browseBtn.disabled = !characterSelect.value && !truthSelect.value;
    }

    characterSelect.addEventListener('change', () => {
        updateBrowseBtn();
        if (!browseMode) {
            currentOffset = 0;
            if (searchInput.value.trim()) {
                search(searchInput.value, 0);
            }
        }
    });

    episodeSelect.addEventListener('change', () => {
        if (statsMode) {
            loadStats();
        } else if (browseMode) {
            browseEpisode = parseInt(episodeSelect.value) || 0;
            browseOffset = 0;
            browseDialogue(browseCharacter, browseOffset, browseEpisode);
        } else if (searchInput.value.trim()) {
            currentOffset = 0;
            search(searchInput.value, 0);
        }
    });

    truthSelect.addEventListener('change', () => {
        updateBrowseBtn();
        if (browseMode) {
            browseOffset = 0;
            browseDialogue(browseCharacter, browseOffset, browseEpisode);
        } else if (searchInput.value.trim()) {
            currentOffset = 0;
            search(searchInput.value, 0);
        }
    });

    browseBtn.addEventListener('click', () => {
        const characterId = characterSelect.value;
        const truth = truthSelect.value;
        if (!characterId && !truth) {
            return;
        }

        browseMode = true;
        quoteMode = false;
        browseCharacter = characterId;
        browseEpisode = parseInt(episodeSelect.value) || 0;
        browseOffset = 0;
        searchInput.value = '';
        browseDialogue(characterId, 0, browseEpisode);
    });

    resultsContainer.addEventListener('click', async (e) => {
        const btn = e.target.closest('.lang-card-btn');
        if (!btn || btn.classList.contains('active')) {
            return;
        }

        const toggle = btn.closest('.lang-card-toggle');
        const audioId = toggle.dataset.audioId;
        const newLang = btn.dataset.lang;

        for (const b of toggle.querySelectorAll('.lang-card-btn')) {
            b.disabled = true;
        }
        try {
            const response = await fetch(`${API_BASE}/quote/${audioId}?lang=${newLang}`);
            const quote = await response.json();
            if (quote.error) {
                return;
            }

            const card = btn.closest('.quote-card') || btn.closest('.featured-quote');
            const textEl = card.querySelector('.quote-text') || card.querySelector('.featured-text');
            if (textEl) {
                const isFeatured = textEl.classList.contains('featured-text');
                textEl.innerHTML = isFeatured
                    ? `"${quote.textHtml || escapeHtml(quote.text)}"`
                    : (quote.textHtml || escapeHtml(quote.text));
            }

            for (const b of toggle.querySelectorAll('.lang-card-btn')) {
                b.classList.toggle('active', b.dataset.lang === newLang);
            }

            const contextSection = card.querySelector('.context-section');
            if (contextSection) {
                const contextAudioId = card.querySelector('.context-btn')?.dataset.audioId;
                if (contextAudioId) {
                    try {
                        const ctxResp = await fetch(`${API_BASE}/context/${contextAudioId}?lang=${newLang}&lines=5`);
                        const ctxData = await ctxResp.json();
                        if (!ctxData.error) {
                            contextSection.innerHTML = renderContextSection(ctxData);
                        }
                    } catch (ctxErr) {
                        console.error('Failed to refresh context:', ctxErr);
                    }
                }
            }
        } catch (error) {
            console.error('Failed to toggle language:', error);
        } finally {
            for (const b of toggle.querySelectorAll('.lang-card-btn')) {
                b.disabled = false;
            }
        }
    });

    async function browseDialogue(characterId, offset = 0, episode = 0) {
        showLoading();

        try {
            const truth = truthSelect.value;
            let url = `${API_BASE}/browse?limit=${PAGE_SIZE}&offset=${offset}&lang=${currentLang}`;
            if (characterId) {
                url += `&character=${characterId}`;
            }
            if (episode > 0) {
                url += `&episode=${episode}`;
            }
            if (truth) {
                url += `&truth=${truth}`;
            }

            const response = await fetch(url);
            const data = await response.json();

            browseOffset = data.offset;
            browseTotal = data.total;

            renderBrowseResults(data);
            updateURL();
        } catch (error) {
            console.error('Browse failed:', error);
            showEmpty('Failed to load dialogue.');
        }
    }

    function renderBrowseResults(data) {
        stopAudio();
        currentAudioId = null;
        if (!data.quotes || data.quotes.length === 0) {
            showEmpty('No dialogue found for this character.');
            return;
        }

        const browseEpLabel = browseEpisode > 0 ? ` — Episode ${browseEpisode}` : '';
        const truthLabel = truthSelect.value === 'red' ? ' — Red Truth' : (truthSelect.value === 'blue' ? ' — Blue Truth' : '');
        const titleName = data.character || 'All Characters';
        const header = `
            <div class="browse-header">
                <h2 class="browse-title">${escapeHtml(titleName)}${browseEpLabel}${truthLabel}</h2>
                <p class="browse-subtitle">Showing lines ${data.offset + 1}-${data.offset + data.quotes.length} of ${data.total} in story order</p>
            </div>
        `;

        const quotes = data.quotes.map((quote, index) => {
            const lineNum = data.offset + index + 1;
            return `
                <article class="quote-card" style="--index: ${index}">
                    <span class="quote-number">#${lineNum}</span>
                    <span class="quote-mark">"</span>
                    <p class="quote-text">${quote.textHtml || escapeHtml(quote.text)}</p>
                    <div class="quote-meta">
                        <span class="quote-character">— ${escapeHtml(quote.character)}</span>
                        <div class="quote-details">
                            ${quote.episode ? `<span class="quote-episode">${episodeLabel(quote)}</span>` : ''}

                            ${langToggleHTML(quote.audioId)}
                        </div>
                    </div>
                    ${audioPlayerHTML(quote.audioId, quote.characterId)}
                    <div class="quote-actions">
                        ${contextBtnHTML(quote.audioId)}
                        ${shareBtnHTML(quote.audioId)}
                    </div>
                </article>
            `;
        }).join('');

        const totalPages = Math.ceil(data.total / PAGE_SIZE);
        const currentPage = Math.floor(data.offset / PAGE_SIZE) + 1;

        const pagination = data.total > PAGE_SIZE ? `
            <div class="pagination">
                <button class="pagination-btn" id="browsePrev" ${data.offset === 0 ? 'disabled' : ''}>◀ Previous</button>
                <span class="pagination-info">Page <span>${currentPage}</span> of <span>${totalPages}</span></span>
                <button class="pagination-btn" id="browseNext" ${data.offset + data.quotes.length >= data.total ? 'disabled' : ''}>Next ▶</button>
            </div>
        ` : '';

        resultsContainer.innerHTML = header + `<div class="quotes-grid">${quotes}</div>` + pagination;

        if (data.total > PAGE_SIZE) {
            document.getElementById('browsePrev')?.addEventListener('click', () => {
                if (browseOffset >= PAGE_SIZE) {
                    browseOffset -= PAGE_SIZE;
                    browseDialogue(browseCharacter, browseOffset, browseEpisode);
                }
            });
            document.getElementById('browseNext')?.addEventListener('click', () => {
                if (browseOffset + PAGE_SIZE < browseTotal) {
                    browseOffset += PAGE_SIZE;
                    browseDialogue(browseCharacter, browseOffset, browseEpisode);
                }
            });
        }
    }

    function updateURL() {
        const params = new URLSearchParams();

        if (statsMode) {
            params.set('stats', '1');
        } else if (browseMode) {
            params.set('browse', browseCharacter || '1');
        } else if (searchInput.value.trim()) {
            params.set('q', searchInput.value.trim());
            if (characterSelect.value) {
                params.set('character', characterSelect.value);
            }
        } else if (quoteMode && currentAudioId) {
            params.set('quote', currentAudioId);
        }

        const episode = episodeSelect.value;
        if (episode && episode !== '0') {
            params.set('episode', episode);
        }

        const truth = truthSelect.value;
        if (truth) {
            params.set('truth', truth);
        }

        const offset = browseMode ? browseOffset : currentOffset;
        if (offset > 0) {
            params.set('offset', offset);
        }

        if (currentLang !== 'en') {
            params.set('lang', currentLang);
        }

        const qs = params.toString();
        const newURL = qs ? `?${qs}` : window.location.pathname;
        history.pushState(null, '', newURL);
    }

    function loadFromURL() {
        const params = new URLSearchParams(window.location.search);

        const lang = params.get('lang') || 'en';
        currentLang = lang;
        for (const b of langBtns) {
            b.classList.toggle('active', b.dataset.lang === lang);
        }

        const episode = params.get('episode') || '0';
        episodeSelect.value = episode;

        const truth = params.get('truth') || '';
        truthSelect.value = truth;

        const browse = params.get('browse');
        const q = params.get('q');
        const offset = parseInt(params.get('offset')) || 0;
        const isStats = params.get('stats') === '1';

        if (isStats) {
            loadStats();
            return;
        }

        const quoteId = params.get('quote');
        if (quoteId) {
            quoteMode = true;
            getQuoteByAudioId(quoteId);
            return;
        }

        if (browse) {
            const isCharacter = browse !== '1';
            if (isCharacter) {
                characterSelect.value = browse;
            }
            updateBrowseBtn();
            browseMode = true;
            browseCharacter = isCharacter ? browse : '';
            browseEpisode = parseInt(episode) || 0;
            browseOffset = offset;
            searchInput.value = '';
            browseDialogue(browseCharacter, offset, browseEpisode);
        } else if (q) {
            searchInput.value = q;
            const character = params.get('character') || '';
            characterSelect.value = character;
            updateBrowseBtn();
            currentOffset = offset;
            search(q, offset);
        } else {
            characterSelect.value = '';
            updateBrowseBtn();
            searchInput.value = '';
            getRandomQuote();
        }
    }

    window.addEventListener('popstate', loadFromURL);

    for (const btn of langBtns) {
        btn.addEventListener('click', () => {
            const newLang = btn.dataset.lang;
            if (newLang === currentLang) {
                return;
            }

            currentLang = newLang;
            for (const b of langBtns) {
                b.classList.toggle('active', b.dataset.lang === newLang);
            }

            if (browseMode) {
                browseDialogue(browseCharacter, browseOffset, browseEpisode);
            } else if (searchInput.value.trim()) {
                search(searchInput.value, currentOffset);
            } else if (currentAudioId) {
                getQuoteByAudioId(currentAudioId);
            } else {
                getRandomQuote();
            }
        });
    }

    function destroyStatsCharts() {
        for (let i = 0; i < statsCharts.length; i++) {
            statsCharts[i].destroy();
        }
        statsCharts = [];
        document.querySelector('.container').classList.remove('stats-active');
    }

    function statsCardHTML(id, title, tall, wide) {
        const tallClass = tall ? ' stats-chart-tall' : '';
        const wideClass = wide ? ' stats-card-wide' : '';
        return `
            <div class="stats-card${wideClass}">
                <div class="stats-card-header">
                    <h3 class="stats-card-title">${title}</h3>
                    <button class="stats-zoom-reset" data-chart-id="${id}">Reset Zoom</button>
                </div>
                <div class="stats-chart-container${tallClass}"><canvas id="${id}"></canvas></div>
                <p class="stats-zoom-hint">Scroll to zoom &middot; drag to pan</p>
            </div>
        `;
    }

    function renderStats(data) {
        destroyStatsCharts();
        stopAudio();
        currentAudioId = null;
        document.querySelector('.container').classList.add('stats-active');

        const hasAllEpisodes = data.linesPerEpisode && data.linesPerEpisode.length > 0;
        const ep = parseInt(episodeSelect.value) || 0;
        const epLabel = ep > 0
            ? (data.episodeNames[ep] ? 'Episode ' + ep + ' — ' + data.episodeNames[ep] : 'Episode ' + ep)
            : 'All Episodes';

        let cards = statsCardHTML('chartTopSpeakers', 'Top Speakers', true, true);
        if (hasAllEpisodes) {
            cards += statsCardHTML('chartLinesPerEpisode', 'Lines per Episode', false, false);
            cards += statsCardHTML('chartTruth', 'Red Truth &amp; Blue Truth', false, false);
        }
        cards += statsCardHTML('chartInteractions', 'Character Interactions', true, true);
        if (hasAllEpisodes) {
            cards += statsCardHTML('chartPresence', 'Character Presence by Episode', false, true);
        }

        const html = `
            <div class="stats-header">
                <h2 class="stats-title">Umineko Statistics</h2>
                <p class="stats-subtitle">${epLabel} — English script lines</p>
            </div>
            <div class="stats-grid">
                ${cards}
            </div>
        `;

        resultsContainer.innerHTML = html;

        Chart.defaults.font.family = "'Cormorant Garamond', serif";
        Chart.defaults.color = '#a89bb8';

        renderTopSpeakersChart(data);
        if (hasAllEpisodes) {
            renderLinesPerEpisodeChart(data);
            renderTruthChart(data);
        }
        renderInteractionsChart(data);
        if (hasAllEpisodes) {
            renderPresenceChart(data);
        }

        resultsContainer.addEventListener('click', (e) => {
            const resetBtn = e.target.closest('.stats-zoom-reset');
            if (!resetBtn) {
                return;
            }
            const chartId = resetBtn.dataset.chartId;
            for (let i = 0; i < statsCharts.length; i++) {
                if (statsCharts[i].canvas.id === chartId) {
                    statsCharts[i].resetZoom();
                    break;
                }
            }
        });
    }

    const PALETTE = [
        '#d4a84b', '#9d7bc9', '#ff3333', '#3399ff', '#6b4c9a',
        '#f0d590', '#8b2942', '#a67c2e', '#3d2a5c', '#e8e0f0',
        '#c97bb4', '#7bc9a3'
    ];

    const zoomConfig = {
        zoom: {
            wheel: { enabled: true },
            pinch: { enabled: true },
            mode: 'xy'
        },
        pan: {
            enabled: true,
            mode: 'xy'
        }
    };

    function renderTopSpeakersChart(data) {
        const labels = [];
        const counts = [];
        for (let i = 0; i < data.topSpeakers.length; i++) {
            labels.push(data.topSpeakers[i].name);
            counts.push(data.topSpeakers[i].count);
        }

        const ctx = document.getElementById('chartTopSpeakers').getContext('2d');
        const chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Lines',
                    data: counts,
                    backgroundColor: '#d4a84b',
                    borderColor: '#a67c2e',
                    borderWidth: 1
                }]
            },
            options: {
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    zoom: zoomConfig
                },
                scales: {
                    x: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    },
                    y: {
                        grid: { display: false },
                        ticks: { color: '#e8e0f0' }
                    }
                }
            }
        });
        statsCharts.push(chart);
    }

    function renderLinesPerEpisodeChart(data) {
        const epLabels = [];
        for (let i = 0; i < data.linesPerEpisode.length; i++) {
            epLabels.push('EP' + data.linesPerEpisode[i].episode + ' ' + data.linesPerEpisode[i].episodeName);
        }

        const charSet = new Set();
        for (let i = 0; i < data.linesPerEpisode.length; i++) {
            const chars = data.linesPerEpisode[i].characters;
            for (const key of Object.keys(chars)) {
                charSet.add(key);
            }
        }

        const charIds = Array.from(charSet).filter(id => id !== 'other');
        charIds.push('other');

        const datasets = [];
        for (let ci = 0; ci < charIds.length; ci++) {
            const id = charIds[ci];
            const label = id === 'other' ? 'Other' : (data.characterNames[id] || id);
            const epData = [];
            for (let ei = 0; ei < data.linesPerEpisode.length; ei++) {
                epData.push(data.linesPerEpisode[ei].characters[id] || 0);
            }
            datasets.push({
                label: label,
                data: epData,
                backgroundColor: PALETTE[ci % PALETTE.length]
            });
        }

        const ctx = document.getElementById('chartLinesPerEpisode').getContext('2d');
        const chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: epLabels,
                datasets: datasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: { color: '#a89bb8', boxWidth: 12 }
                    },
                    zoom: zoomConfig
                },
                scales: {
                    x: {
                        stacked: true,
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    },
                    y: {
                        stacked: true,
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    }
                }
            }
        });
        statsCharts.push(chart);
    }

    function renderTruthChart(data) {
        const labels = [];
        const redData = [];
        const blueData = [];
        for (let i = 0; i < data.truthPerEpisode.length; i++) {
            labels.push('EP' + data.truthPerEpisode[i].episode);
            redData.push(data.truthPerEpisode[i].red);
            blueData.push(data.truthPerEpisode[i].blue);
        }

        const ctx = document.getElementById('chartTruth').getContext('2d');
        const chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'Red Truth',
                        data: redData,
                        backgroundColor: '#ff3333',
                        borderColor: '#cc0000',
                        borderWidth: 1
                    },
                    {
                        label: 'Blue Truth',
                        data: blueData,
                        backgroundColor: '#3399ff',
                        borderColor: '#0066cc',
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: { color: '#a89bb8' }
                    },
                    zoom: zoomConfig
                },
                scales: {
                    x: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    },
                    y: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    }
                }
            }
        });
        statsCharts.push(chart);
    }

    function renderInteractionsChart(data) {
        const labels = [];
        const counts = [];
        for (let i = 0; i < data.interactions.length; i++) {
            labels.push(data.interactions[i].nameA + ' & ' + data.interactions[i].nameB);
            counts.push(data.interactions[i].count);
        }

        const ctx = document.getElementById('chartInteractions').getContext('2d');
        const chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: labels,
                datasets: [{
                    label: 'Adjacent Lines',
                    data: counts,
                    backgroundColor: '#9d7bc9',
                    borderColor: '#6b4c9a',
                    borderWidth: 1
                }]
            },
            options: {
                indexAxis: 'y',
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    zoom: zoomConfig
                },
                scales: {
                    x: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    },
                    y: {
                        grid: { display: false },
                        ticks: { color: '#e8e0f0', font: { size: 11 } }
                    }
                }
            }
        });
        statsCharts.push(chart);
    }

    function renderPresenceChart(data) {
        const epLabels = [];
        for (let ep = 1; ep <= 8; ep++) {
            epLabels.push('EP' + ep);
        }

        const datasets = [];
        for (let i = 0; i < data.characterPresence.length; i++) {
            const cp = data.characterPresence[i];
            datasets.push({
                label: cp.name,
                data: cp.episodes,
                backgroundColor: PALETTE[i % PALETTE.length],
                borderColor: PALETTE[i % PALETTE.length],
                borderWidth: 1
            });
        }

        const ctx = document.getElementById('chartPresence').getContext('2d');
        const chart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: epLabels,
                datasets: datasets
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: { color: '#a89bb8', boxWidth: 12 }
                    },
                    zoom: zoomConfig
                },
                scales: {
                    x: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    },
                    y: {
                        grid: { color: 'rgba(61, 42, 92, 0.4)' },
                        ticks: { color: '#a89bb8' }
                    }
                }
            }
        });
        statsCharts.push(chart);
    }

    async function loadConfig() {
        try {
            const response = await fetch(`${API_BASE}/config`);
            if (response.ok) {
                const config = await response.json();
                hasAudio = config.hasAudio === true;
            }
        } catch (e) {
            console.warn('Failed to load config:', e);
        }
    }

    createButterflies();
    await loadConfig();
    await loadCharacters();
    loadFromURL();
})();
