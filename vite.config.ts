import { resolve } from 'node:path';
import terser from '@rollup/plugin-terser';
import react from '@vitejs/plugin-react';
import { rollupImportMapPlugin } from 'rollup-plugin-import-map';
import { defineConfig } from 'vite';
import cssInjectedByJsPlugin from 'vite-plugin-css-injected-by-js';
import tsconfigPaths from 'vite-tsconfig-paths';

export default defineConfig({
  plugins: [
    tsconfigPaths(),
    react(),
    cssInjectedByJsPlugin(),
    {
      ...rollupImportMapPlugin([
        {
          imports: {
            react: 'https://www.nav.no/tms-min-side-assets/react/18/esm/index.js',
            'react-dom': 'https://www.nav.no/tms-min-side-assets/react-dom/18/esm/index.js',
          },
        },
      ]),
      enforce: 'pre',
      apply: 'build',
    },
    terser(),
  ],
  build: {
    rollupOptions: {
      input: resolve(__dirname, 'src/microfrontend.tsx'),
      preserveEntrySignatures: 'exports-only',
      output: {
        assetFileNames: 'mine-klager-microfrontend.[hash][extname]',
        entryFileNames: 'mine-klager-microfrontend.[hash].js',
        format: 'esm',
      },
    },
  },
  server: { port: 3001 },
});
