import tailwindcss from "@tailwindcss/vite";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig, loadEnv } from "vite";
import configLocal from "./config/config.json";
import configProd from "./config/config.prod.json";

export default defineConfig(({ mode }) => {
  process.env = { ...process.env, ...loadEnv(mode, process.cwd()) };
  const env = process.env.VITE_ENV;
  const config = env === "prod" ? configProd : configLocal;
  return {
    define: {
      __GLOBAL_CONFIG__: config,
    },
    plugins: [tailwindcss(), sveltekit()],
  };
});
