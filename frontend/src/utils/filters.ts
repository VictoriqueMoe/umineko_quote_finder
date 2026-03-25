import type { FilterState } from "../types/app";
import { normalizeCharacterKey } from "./characterIds";

export function enforceMutuallyExclusiveFilters(filters: FilterState): FilterState {
    if (filters.character && (filters.interactionA || filters.interactionB)) {
        return { ...filters, character: "" };
    }
    return filters;
}

export function normalizeFilterCharacters(filters: FilterState): FilterState {
    return {
        ...filters,
        character: normalizeCharacterKey(filters.character),
        interactionA: normalizeCharacterKey(filters.interactionA),
        interactionB: normalizeCharacterKey(filters.interactionB),
    };
}
