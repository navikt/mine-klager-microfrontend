import { appendFileSync, mkdirSync, writeFileSync } from 'node:fs';
import { renderToStaticMarkup } from 'react-dom/server';
import { buildCss } from '@/css';
import { formatBytes } from '@/format';
import { Microfrontend, type MicrofrontendProps } from '@/microfrontend';

const BASE_URL_PLACEHOLDER = '{{BASE_URL}}';

const variants: Record<string, MicrofrontendProps> = {
  nb: {
    title: 'Mine saker hos Klageinstans',
    description: 'Her kan du se status på dine saker hos Klageinstans.',
    url: BASE_URL_PLACEHOLDER,
  },
  nn: {
    title: 'Mine saker hjå Klageinstans',
    description: 'Her kan du sjå status på dine saker hjå Klageinstans.',
    url: `${BASE_URL_PLACEHOLDER}/nn`,
  },
  en: {
    title: 'My cases with Nav Complaints Unit',
    description: 'Here you can see the status of your cases with Nav Complaints Unit (Klageinstans).',
    url: `${BASE_URL_PLACEHOLDER}/en`,
  },
};

const DIST_DIR = 'dist';

mkdirSync(DIST_DIR, { recursive: true });

// 1. Render all markup variants first.
const rendered = Object.entries(variants).map(([lang, props]) => ({
  lang,
  markup: renderToStaticMarkup(<Microfrontend {...props} />),
}));

// 2. Dynamically discover which component CSS files are needed based on all rendered HTML.
const allMarkups = rendered.map((r) => r.markup);
const { css, tokensBefore, tokensAfter, tokensReduction, tokensReductionPercent } = buildCss(...allMarkups);

// 3. Write the output files.
const TE = new TextEncoder();

const cssSize = formatBytes(TE.encode(css).byteLength);

interface BuildStats {
  file: string;
  htmlSize: string;
  cssSize: string;
  totalSize: string;
}

const stats: BuildStats[] = [];

console.info(`Generating microfrontend variants in ${DIST_DIR}/...`);

for (const { lang, markup } of rendered) {
  const file = `${lang}.html`;

  const html = `
<mine-klager-microfrontend>
  <template shadowrootmode="open">
  <style>${css}</style>
  ${markup}
  </template>
</mine-klager-microfrontend>
<script>
  (function() {
    var el = document.currentScript.previousElementSibling;
    if (el.shadowRoot === null) {
      el.attachShadow({ mode: "open" }).appendChild(el.firstElementChild.content);
    }
  })()
</script>`.trim();

  writeFileSync(`${DIST_DIR}/${file}`, html);

  const microfrontendSize = formatBytes(TE.encode(markup).byteLength);
  const totalFileSize = formatBytes(TE.encode(html).byteLength);

  stats.push({ file, htmlSize: microfrontendSize, cssSize, totalSize: totalFileSize });

  console.info(`Generated ${DIST_DIR}/${file} (${microfrontendSize} HTML, ${cssSize} CSS, ${totalFileSize} total)`);
}

const summaryFile = process.env.GITHUB_STEP_SUMMARY;

if (summaryFile !== undefined) {
  const statsRows = stats.map((s) => `| \`${s.file}\` | ${s.htmlSize} | ${s.cssSize} | ${s.totalSize} |`).join('\n');

  const summary = `
### 📦 Microfrontend Build Stats

| File | HTML | CSS | Total |
| --- | ---: | ---: | ---: |
${statsRows}

🌳 **Token tree-shaking:** ${formatBytes(tokensBefore)} → ${formatBytes(tokensAfter)} (${formatBytes(tokensReduction)} / ${tokensReductionPercent}% reduction)`.trim();

  appendFileSync(summaryFile, summary);
}
