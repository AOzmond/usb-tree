import { defineConfig } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import svg from "@poppanator/sveltekit-svg"
import { resolve } from "path"

const projectRootDir = import.meta.dirname

// https://vitejs.dev/config/
/** @type {import('vite').UserConfig} */
export default defineConfig({
  plugins: [svelte(), svg()],
  resolve: {
    alias: {
      $style: resolve(projectRootDir, "src/style"),
      $lib: resolve(projectRootDir, "src/lib"),
      $wailsjs: resolve(projectRootDir, "wailsjs"),
    },
  },
})
