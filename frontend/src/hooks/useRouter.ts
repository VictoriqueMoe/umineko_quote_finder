import { useCallback, useEffect, useRef, useState } from "react";
import type { FilterState, Language, ViewMode } from "../types/app";
import { buildUrl, type ParsedRoute, parseRoute } from "./routeDefinitions";

export type { ParsedRoute };

export interface NavigateOpts {
    searchOffset?: number;
    browseOffset?: number;
    currentAudioId?: string | null;
    searchQuery?: string;
}

interface UseRouterParams {
    language: Language;
    onRouteLoad: (route: ParsedRoute) => void | Promise<void>;
}

export function useRouter({ language, onRouteLoad }: UseRouterParams) {
    const [viewMode, setViewMode] = useState<ViewMode>("featured");
    const initialised = useRef(false);
    const onRouteLoadRef = useRef(onRouteLoad);
    const languageRef = useRef(language);

    useEffect(() => {
        onRouteLoadRef.current = onRouteLoad;
    });
    useEffect(() => {
        languageRef.current = language;
    });

    const loadFromURL = useCallback(() => {
        const route = parseRoute(window.location.search);
        Promise.resolve(onRouteLoadRef.current(route)).then(() => {
            setViewMode(route.viewMode);
            initialised.current = true;
        });
    }, []);

    useEffect(() => {
        loadFromURL();
        window.addEventListener("popstate", loadFromURL);
        return () => {
            window.removeEventListener("popstate", loadFromURL);
        };
    }, [loadFromURL]);

    const navigate = useCallback((vm: ViewMode, filters: FilterState, opts?: NavigateOpts) => {
        if (!initialised.current) {
            return;
        }
        setViewMode(vm);
        const url = buildUrl(
            {
                viewMode: vm,
                filters,
                currentAudioId: opts?.currentAudioId ?? null,
                searchOffset: opts?.searchOffset ?? 0,
                browseOffset: opts?.browseOffset ?? 0,
            },
            languageRef.current,
            opts?.searchQuery ?? "",
        );
        history.pushState(null, "", url);
    }, []);

    return { viewMode, navigate };
}
