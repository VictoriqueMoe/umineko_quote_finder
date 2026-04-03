import type { FilterState, Game } from "../types/app";
import { normalizeCharacterKey } from "./characterIds";

export function enforceMutuallyExclusiveFilters(filters: FilterState): FilterState {
    if (filters.character && (filters.interactionA || filters.interactionB)) {
        return { ...filters, character: "" };
    }
    return filters;
}

export function normalizeFilterCharacters(filters: FilterState, game?: Game): FilterState {
    return {
        ...filters,
        character: normalizeCharacterKey(filters.character, game),
        interactionA: normalizeCharacterKey(filters.interactionA, game),
        interactionB: normalizeCharacterKey(filters.interactionB, game),
    };
}
