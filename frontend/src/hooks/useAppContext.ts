import { useContext } from "react";
import { AppContext, type AppContextValue } from "../context/appContextValue";

export function useAppContext(): AppContextValue {
    const ctx = useContext(AppContext);
    if (!ctx) {
        throw new Error("useAppContext must be used within AppProvider");
    }
    return ctx;
}
