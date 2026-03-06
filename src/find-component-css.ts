import { readdirSync, readFileSync } from 'node:fs';
import { join } from 'node:path';

const COMPONENT_CSS_DIR = 'node_modules/@navikt/ds-css/dist/component';

const CLASS_ATTR_PATTERN = /class="([^"]*)"/g;
const CSS_CLASS_SELECTOR_PATTERN = /\.([a-zA-Z_][\w-]*)/g;
const WHITESPACE_PATTERN = /\s+/;

/**
 * Extracts all CSS class names from an HTML string by parsing `class="..."` attributes.
 */
const extractClassesFromHtml = (html: string): Set<string> => {
  const classes = new Set<string>();

  for (const [, classList] of html.matchAll(CLASS_ATTR_PATTERN)) {
    if (classList !== undefined) {
      for (const className of classList.split(WHITESPACE_PATTERN)) {
        if (className.length > 0) {
          classes.add(className);
        }
      }
    }
  }

  return classes;
};

/**
 * Extracts all class names referenced as selectors in a CSS string.
 */
const extractClassesFromCss = (css: string): Set<string> => {
  const classes = new Set<string>();

  for (const [, className] of css.matchAll(CSS_CLASS_SELECTOR_PATTERN)) {
    if (className !== undefined) {
      classes.add(className);
    }
  }

  return classes;
};

interface ComponentCssFile {
  name: string;
  content: string;
  classes: Set<string>;
}

/** All `*.min.css` files in the component directory, pre-read and indexed. */
const componentCssFiles: ComponentCssFile[] = readdirSync(COMPONENT_CSS_DIR)
  .filter((file) => file.endsWith('.min.css'))
  .toSorted()
  .map((file) => {
    const content = readFileSync(join(COMPONENT_CSS_DIR, file), 'utf-8');

    return {
      name: file,
      content,
      classes: extractClassesFromCss(content),
    };
  });

/**
 * Given one or more rendered HTML strings, finds all component CSS files from
 * `@navikt/ds-css/dist/component/` whose selectors reference at least one class
 * present in the HTML.
 *
 * This replaces the need to manually specify which component CSS files to include.
 *
 * @returns An object with `files` (the matched file names) and `css` (their concatenated content).
 */
export const findComponentCss = (...htmlStrings: string[]): { files: string[]; css: string } => {
  const htmlClasses = new Set<string>();

  for (const html of htmlStrings) {
    for (const cls of extractClassesFromHtml(html)) {
      htmlClasses.add(cls);
    }
  }

  const matched: ComponentCssFile[] = [];

  for (const file of componentCssFiles) {
    const hasOverlap = file.classes.values().some((cls) => htmlClasses.has(cls));

    if (hasOverlap) {
      matched.push(file);
    }
  }

  return {
    files: matched.map((f) => f.name),
    css: matched.map((f) => f.content).join('\n'),
  };
};
