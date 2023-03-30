import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react({
            babel: {
                plugins: ["babel-plugin-macros", "babel-plugin-styled-components"],
            },
        }),
    ],
    base: "",
    assetsInclude: ["**/*.otf"],
    build: {
        outDir: "../../int/api/dist/plugin-manager",
        emptyOutDir: true,
        assetsDir: "./",
    },
    esbuild: {
        // https://github.com/vitejs/vite/issues/8644#issuecomment-1159308803
        logOverride: { "this-is-undefined-in-esm": "silent" },
    },
});
