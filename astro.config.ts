import node from '@astrojs/node';
import react from '@astrojs/react';
import tailwindcss from '@tailwindcss/vite';
import type { ViteUserConfig } from 'astro';
import { defineConfig } from 'astro/config';
import prefixer from 'postcss-prefix-selector';
import { rollupImportMapPlugin } from 'rollup-plugin-import-map';

// https://astro.build/config
export default defineConfig({
  build: {
    assetsPrefix: 'https://cdn.nav.no/min-side/mine-klager-microfrontend',
    inlineStylesheets: 'always',
  },
  vite: {
    css: {
      postcss: {
        plugins: [
          prefixer({
            prefix: '.mine-klager-microfrontend',
            ignoreFiles: [/module.css/],
          }),
        ],
      },
    },

    plugins: [tailwindcss()],
  },
  integrations: [
    react(),
    {
      name: 'importmap',
      hooks: {
        'astro:build:setup': ({ vite, target }) => {
          if (target === 'client') {
            const pluginOptions: ViteUserConfig['plugins'] = [
              {
                ...rollupImportMapPlugin({
                  imports: {
                    react: 'https://www.nav.no/tms-min-side-assets/react/18/esm/index.js',
                    'react-dom': 'https://www.nav.no/tms-min-side-assets/react-dom/18/esm/index.js',
                  },
                }),
                enforce: 'pre',
                apply: 'build',
              },
            ];

            if (vite.plugins === undefined) {
              vite.plugins = pluginOptions;
            } else {
              vite.plugins = [...vite.plugins, ...pluginOptions];
            }
          }
        },
      },
    },
  ],
  i18n: {
    defaultLocale: 'nb',
    locales: ['nb', 'nn', 'en'],
    routing: {
      prefixDefaultLocale: true,
    },
  },
  output: 'server',
  adapter: node({
    mode: 'standalone',
  }),
});
