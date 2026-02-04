import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { BarElement, CategoryScale, Chart as ChartJS, Legend, LinearScale, Title, Tooltip } from "chart.js";
import zoomPlugin from "chartjs-plugin-zoom";
import App from "./App";
import { AppProvider } from "./context/AppContext";
import "./styles/variables.css";
import "./styles/global.css";

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend, zoomPlugin);

ChartJS.defaults.font.family = "'Cormorant Garamond', serif";
ChartJS.defaults.color = "#a89bb8";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <AppProvider>
            <App />
        </AppProvider>
    </StrictMode>,
);
