import { useCallback, useState } from "react";
import { useAppContext } from "./useAppContext";
import type { Quote } from "../types/api";
import type { Language } from "../types/app";

const CONTENT_TYPE_LABELS: Record<string, string> = {
    tea: "Tea Party",
    ura: "????",
    omake: "Omake",
};

export function episodeLabel(quote: Quote): string {
    if (!quote.episode) {
        return "";
    }
    let label = `Episode ${quote.episode}`;
    if (quote.contentType && CONTENT_TYPE_LABELS[quote.contentType]) {
        label += ` \u2014 ${CONTENT_TYPE_LABELS[quote.contentType]}`;
    }
    return label;
}

export function useQuoteDisplay(quote: Quote, langOverride?: Exclude<Language, "auto">) {
    const { language, hasAudio, game } = useAppContext();
    const effective = langOverride ?? (language === "auto" ? "en" : language);
    const gameHasAudio = hasAudio && game === "umineko";

    const [textOverride, setTextOverride] = useState<string | null>(null);
    const [langUserOverride, setLangUserOverride] = useState<Language | null>(null);
    const [prevQuote, setPrevQuote] = useState(quote);
    const [prevEffective, setPrevEffective] = useState(effective);

    if (quote !== prevQuote) {
        setPrevQuote(quote);
        setTextOverride(null);
        setLangUserOverride(null);
    }
    if (effective !== prevEffective) {
        setPrevEffective(effective);
        setLangUserOverride(null);
    }

    const activeLang = langUserOverride ?? effective;
    let displayHtml = textOverride ?? (quote.textHtml || quote.text);
    if (!textOverride && (game === "higurashi" || game === "ciconia") && activeLang === "ja" && quote.textJp) {
        displayHtml = quote.textJpHtml || quote.textJp;
    }
    const lang = activeLang;

    const handleTextUpdate = useCallback((textHtml: string) => {
        setTextOverride(textHtml);
    }, []);

    const handleLangChange = useCallback((newLang: Language) => {
        setLangUserOverride(newLang);
    }, []);

    return { displayHtml, lang, hasAudio: gameHasAudio, handleTextUpdate, handleLangChange };
}
