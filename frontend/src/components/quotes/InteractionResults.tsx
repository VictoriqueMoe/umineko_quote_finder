import { useAppContext } from "../../hooks/useAppContext";
import type { Quote } from "../../types/api";
import { normalizeCharacterKey, toCharacterId } from "../../utils/characterIds";

interface InteractionResultsProps {
    mode: "browse" | "search";
    quotes: Quote[];
    offset: number;
    total: number;
    interactionA: string;
    interactionB: string;
    characterFilter?: string;
    episodeFilter?: string;
    truthFilter?: string;
    query?: string;
    onContextQuoteClick?: (audioId: string) => void;
}

type InteractionBlock = {
    first: Quote;
    second?: Quote;
    startNumber: number;
    endNumber: number;
};

function firstAudioId(audioId?: string): string {
    return audioId ? audioId.split(", ")[0] : "";
}

function buildBlocks(
    mode: "browse" | "search",
    quotes: Quote[],
    offset: number,
    interactionA: string,
    interactionB: string,
): InteractionBlock[] {
    if (mode === "search") {
        return quotes.map((quote, index) => ({
            first: quote,
            startNumber: offset + index + 1,
            endNumber: offset + index + 1,
        }));
    }

    const pairIds = new Set([toCharacterId(interactionA), toCharacterId(interactionB)]);
    const blocks: InteractionBlock[] = [];
    let i = 0;
    while (i < quotes.length) {
        const first = quotes[i];
        const next = i + 1 < quotes.length ? quotes[i + 1] : undefined;
        const firstNumber = offset + i + 1;
        if (
            next &&
            first.characterId &&
            next.characterId &&
            pairIds.has(first.characterId) &&
            pairIds.has(next.characterId) &&
            first.characterId !== next.characterId
        ) {
            blocks.push({ first, second: next, startNumber: firstNumber, endNumber: firstNumber + 1 });
            i += 2;
            continue;
        }
        blocks.push({ first, startNumber: firstNumber, endNumber: firstNumber });
        i += 1;
    }
    return blocks;
}

function resolveCharacterName(characters: Record<string, string>, value: string): string {
    const key = normalizeCharacterKey(value);
    return characters[key] || characters[value] || value;
}

function buildFilterText(
    characters: Record<string, string>,
    characterFilter?: string,
    episodeFilter?: string,
    truthFilter?: string,
): string {
    const parts: string[] = [];
    if (characterFilter) {
        parts.push(`Character: ${resolveCharacterName(characters, characterFilter)}`);
    }
    if (episodeFilter && episodeFilter !== "0") {
        parts.push(`Episode ${episodeFilter}`);
    }
    if (truthFilter === "red") {
        parts.push("Red Truth");
    } else if (truthFilter === "blue") {
        parts.push("Blue Truth");
    }
    return parts.length > 0 ? ` Filters: ${parts.join(" | ")}.` : "";
}

export function InteractionResults({
    mode,
    quotes,
    offset,
    total,
    interactionA,
    interactionB,
    characterFilter,
    episodeFilter,
    truthFilter,
    query,
    onContextQuoteClick,
}: InteractionResultsProps) {
    const { characters } = useAppContext();
    const blocks = buildBlocks(mode, quotes, offset, interactionA, interactionB);
    const filterText = buildFilterText(characters, characterFilter, episodeFilter, truthFilter);
    const pairLabel = `${resolveCharacterName(characters, interactionA)} x ${resolveCharacterName(characters, interactionB)}`;

    const rangeStart = offset + 1;
    const rangeEnd = offset + quotes.length;
    const title = mode === "browse" ? "Interaction Browse" : "Interaction Search";
    const lineLabel = mode === "browse" ? "Line" : "Match";
    const subtitle =
        mode === "browse"
            ? `Showing lines ${rangeStart}-${rangeEnd} of ${total} for ${pairLabel}.${filterText}`
            : `Showing matches ${rangeStart}-${rangeEnd} of ${total} for "${query || ""}" within ${pairLabel}.${filterText}`;

    return (
        <>
            <div className="browse-header interaction-browse-header">
                <h2 className="browse-title">{title}</h2>
                <p className="browse-subtitle">{subtitle}</p>
            </div>
            <div className="interaction-results interaction-results-transcript">
                {blocks.map((block, index) => {
                    const firstId = firstAudioId(block.first.audioId);
                    const secondId = firstAudioId(block.second?.audioId);
                    return (
                        <article
                            key={`${firstId || `${mode}-${offset}-${index}`}-${block.startNumber}`}
                            className={`interaction-result-card interaction-transcript-card${block.second ? " is-exchange" : " is-single"}`}
                        >
                            <div className="interaction-result-meta interaction-transcript-meta">
                                <span>
                                    {lineLabel} #{block.startNumber}
                                    {block.endNumber !== block.startNumber ? `-${block.endNumber}` : ""}
                                </span>
                                <span className="interaction-transcript-kind">
                                    {block.second ? "Exchange" : mode === "browse" ? "Single Line" : "Matched Line"}
                                </span>
                                <span>{block.first.episode ? `Episode ${block.first.episode}` : ""}</span>
                            </div>

                            <div className="context-line interaction-result-line interaction-transcript-line">
                                <span className="context-character">{block.first.character}</span>
                                <span
                                    className="context-text"
                                    dangerouslySetInnerHTML={{ __html: block.first.textHtml || block.first.text }}
                                />
                            </div>

                            {block.second && (
                                <div className="context-line interaction-result-line interaction-transcript-line">
                                    <span className="context-character">{block.second.character}</span>
                                    <span
                                        className="context-text"
                                        dangerouslySetInnerHTML={{ __html: block.second.textHtml || block.second.text }}
                                    />
                                </div>
                            )}

                            {(firstId || secondId) && onContextQuoteClick && (
                                <div className="interaction-result-actions">
                                    {firstId && (
                                        <button className="context-btn" onClick={() => onContextQuoteClick(firstId)}>
                                            Open {block.first.character} Quote
                                        </button>
                                    )}
                                    {secondId && block.second && (
                                        <button className="context-btn" onClick={() => onContextQuoteClick(secondId)}>
                                            Open {block.second.character} Quote
                                        </button>
                                    )}
                                </div>
                            )}
                        </article>
                    );
                })}
            </div>
        </>
    );
}
