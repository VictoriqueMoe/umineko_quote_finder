export function Footer() {
    return (
        <footer className="footer">
            <div className="footer-ornament">{"\u2666 \u2663 \u2665 \u2660"}</div>
            <p className="footer-text">Without love, it cannot be seen.</p>
            <p className="footer-credit">
                {"Umineko no Naku Koro ni \u00A9 "}
                <a href="https://07th-expansion.net/" target="_blank" rel="noopener" className="footer-author">
                    07th Expansion
                </a>
                {" - this is an unofficial fan project"}
            </p>
            <p className="footer-made-by">
                {"Made with \u2764 by "}
                <a href="https://x.com/FeatherineFAA" target="_blank" rel="noopener" className="footer-author">
                    Featherine Augustus Aurora
                </a>
            </p>
            <div className="footer-links">
                <a
                    href="https://github.com/VictoriqueMoe/umineko_quote_finder"
                    target="_blank"
                    rel="noopener"
                    className="footer-link"
                >
                    <svg viewBox="0 0 16 16" aria-hidden="true">
                        <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0016 8c0-4.42-3.58-8-8-8z" />
                    </svg>
                    Source
                </a>
                <a href="/swagger/index.html" target="_blank" rel="noopener" className="footer-link">
                    <svg viewBox="0 0 16 16" aria-hidden="true">
                        <path d="M2 1h12a1 1 0 011 1v12a1 1 0 01-1 1H2a1 1 0 01-1-1V2a1 1 0 011-1zm1 3v2h3V4H3zm0 4v2h3V8H3zm5-4v2h5V4H8zm0 4v2h5V8H8zm-5 4v1h10v-1H3z" />
                    </svg>
                    API Docs
                </a>
            </div>
            <p className="footer-support">
                {"Support 07th Expansion - "}
                <a href="https://store.steampowered.com/app/406550/" target="_blank" rel="noopener" className="footer-author">
                    get the game on Steam
                </a>
            </p>
        </footer>
    );
}
