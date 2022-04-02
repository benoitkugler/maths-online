import vue from "@vitejs/plugin-vue";
import vuetify from "@vuetify/vite-plugin";
import { defineConfig } from "vite";

const path = require("path");

// https://vitejs.dev/config/
export default defineConfig(({ command }) => ({
  base: command == "serve" ? "/prof/" : "/static/prof/",
  plugins: [
    vue({
      reactivityTransform: true
    }),
    // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
    vuetify({
      autoImport: true
    })
  ],
  define: {
    "process.env": {
      VERSION: require("./package.json").version
    }
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src")
    }
  }
  /* remove the need to specify .vue files https://vitejs.dev/config/#resolve-extensions
  resolve: {
    extensions: [
      '.js',
      '.json',
      '.jsx',
      '.mjs',
      '.ts',
      '.tsx',
      '.vue',
    ]
  },
  */
}));
