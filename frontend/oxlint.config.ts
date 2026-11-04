import { defineConfig } from "oxlint";

export default defineConfig({
    plugins: ["unicorn", "typescript", "oxc", "react"],
    rules: {
        "eslint/curly": ["error", "all"],
        "typescript/no-floating-promises": "off",
    },
});
