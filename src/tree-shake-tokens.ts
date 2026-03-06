const DECL_LINE_PATTERN = /^\s*(--[\w-]+)\s*:/;
const EMPTY_BLOCK_PATTERN = /^[^\n{}]*\{\s*\}\s*$/gm;
const EXCESS_NEWLINES_PATTERN = /\n{3,}/g;

/**
 * Tree-shakes a CSS tokens file to only keep custom property declarations
 * that are directly or indirectly referenced by the consumer CSS.
 *
 * @param tokensCss - The full CSS tokens source (e.g. design-system tokens).
 * @param consumerCss - The CSS that actually consumes tokens via `var()`.
 * @returns The tree-shaken tokens CSS containing only referenced declarations.
 */
export const treeShakeTokens = (tokensCss: string, consumerCss: string): string => {
  const needed = extractVarRefs(consumerCss);

  const depsOf = new Map<string, Set<string>>();

  for (const { property, value } of parseDeclarations(tokensCss)) {
    const refs = extractVarRefs(value);

    if (refs.size > 0) {
      const deps = depsOf.get(property) ?? new Set<string>();
      for (const r of refs) {
        deps.add(r);
      }
      depsOf.set(property, deps);
    }
  }

  const queue = [...needed];
  let prop = queue.pop();

  while (prop !== undefined) {
    for (const dep of depsOf.get(prop) ?? []) {
      if (!needed.has(dep)) {
        needed.add(dep);
        queue.push(dep);
      }
    }

    prop = queue.pop();
  }

  const kept = tokensCss
    .split('\n')
    .filter((line) => {
      const property = line.match(DECL_LINE_PATTERN)?.[1];

      return property === undefined || needed.has(property);
    })
    .join('\n');

  let result = kept;
  let prev: string;

  do {
    prev = result;
    result = result.replace(EMPTY_BLOCK_PATTERN, '').replace(EXCESS_NEWLINES_PATTERN, '\n\n');
  } while (result !== prev);

  return `${result.trim()}\n`;
};

const VAR_REF_PATTERN = /var\((--[\w-]+)/g;

const extractVarRefs = (css: string): Set<string> => {
  const refs = new Set<string>();

  for (const [, ref] of css.matchAll(VAR_REF_PATTERN)) {
    if (ref !== undefined) {
      refs.add(ref);
    }
  }

  return refs;
};

const DECLARATION_PATTERN = /^\s*(--[\w-]+)\s*:\s*(.+?)\s*;/gm;

const parseDeclarations = (css: string) => {
  const declarations = [];

  for (const [, property, value] of css.matchAll(DECLARATION_PATTERN)) {
    if (property !== undefined && value !== undefined) {
      declarations.push({ property, value });
    }
  }

  return declarations;
};
