import { readdirSync, readFileSync } from 'node:fs';
import { findComponentCss } from '@/find-component-css';
import { formatBytes } from '@/format';
import { treeShakeTokens } from '@/tree-shake-tokens';

const GLOBAL_CSS_DIR = 'node_modules/@navikt/ds-css/dist/global';

const globalMinCssFiles = readdirSync(GLOBAL_CSS_DIR)
  .filter((file) => file.endsWith('.min.css') && !file.startsWith('tokens.'))
  .toSorted();

const globalCssByName = new Map<string, string>();

for (const file of globalMinCssFiles) {
  const content = readFileSync(`${GLOBAL_CSS_DIR}/${file}`, 'utf-8');
  globalCssByName.set(file, content);
}

const nonTokenGlobalCss = globalMinCssFiles
  .map((file) => globalCssByName.get(file))
  .filter((content) => content !== undefined);

const globalCss = nonTokenGlobalCss.join('\n');

console.info(`Included global CSS files: ${globalMinCssFiles.join(', ')}`);

const tokens = readFileSync(`${GLOBAL_CSS_DIR}/tokens.css`, 'utf-8');

export interface BuildCssResult {
  css: string;
  tokensBefore: number;
  tokensAfter: number;
  tokensReduction: number;
  tokensReductionPercent: number;
  componentFiles: string[];
}

/**
 * Builds the final CSS for the microfrontend by:
 * 1. Dynamically discovering which component CSS files are needed based on the classes in the HTML markup.
 * 2. Combining global CSS (excluding tokens) with the matched component CSS.
 * 3. Tree-shaking tokens to only keep custom properties referenced by the combined CSS.
 * 4. Rewriting `:root` to `:host` for Shadow DOM scoping.
 *
 * @param htmlMarkups - One or more rendered HTML strings to scan for CSS class usage.
 */
export const buildCss = (...htmlMarkups: string[]): BuildCssResult => {
  const { files: componentFiles, css: componentCss } = findComponentCss(...htmlMarkups);

  console.info(`Matched component CSS files: ${componentFiles.join(', ')}`);

  const allComponentCss = [globalCss, componentCss].join('\n');

  const treeShakenTokens = treeShakeTokens(tokens, allComponentCss);

  const tokensBefore = new TextEncoder().encode(tokens).byteLength;
  const tokensAfter = new TextEncoder().encode(treeShakenTokens).byteLength;
  const tokensReduction = tokensBefore - tokensAfter;
  const tokensReductionPercent = Math.round((1 - tokensAfter / tokensBefore) * 100);

  console.info(
    `Tree-shaken tokens.css: ${formatBytes(tokensBefore)} → ${formatBytes(tokensAfter)} (${formatBytes(tokensReduction)} / ${tokensReductionPercent}% reduction)`,
  );

  const rawCss = [treeShakenTokens, allComponentCss].join('\n');

  // Inside Shadow DOM, :root refers to the document root — not the shadow host.
  // Rewrite :root to :host so custom properties are scoped to the shadow boundary.
  // Then deduplicate any resulting `:host, :host` pairs (the source CSS already contains :host).
  const css = rawCss.replace(/:root/g, ':host').replace(/(:host)(?:\s*,\s*:host)+/g, '$1');

  return {
    css,
    tokensBefore,
    tokensAfter,
    tokensReduction,
    tokensReductionPercent,
    componentFiles,
  };
};
