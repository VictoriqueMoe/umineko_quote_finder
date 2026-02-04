import { useCallback, useEffect, useState } from "react";
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

export function useQuoteDisplay(quote: Quote) {
    const { language, hasAudio } = useAppContext();
    const [displayHtml, setDisplayHtml] = useState(quote.textHtml || quote.text);
    const [lang, setLang] = useState<Language>(language);

    useEffect(() => {
        setDisplayHtml(quote.textHtml || quote.text);
    }, [quote]);

    useEffect(() => {
        setLang(language);
    }, [language]);

    const handleTextUpdate = useCallback((textHtml: string) => {
        setDisplayHtml(textHtml);
    }, []);

    const handleLangChange = useCallback((newLang: Language) => {
        setLang(newLang);
    }, []);

    return { displayHtml, lang, hasAudio, handleTextUpdate, handleLangChange };
}
