import { describe, expect, test } from 'bun:test';
import { treeShakeTokens } from './tree-shake-tokens';

const EXCESSIVE_NEWLINES_PATTERN = /\n{3,}/;
const SINGLE_TRAILING_NEWLINE_PATTERN = /[^\n]\n$/;

describe('treeShakeTokens', () => {
  test('keeps a directly referenced token', () => {
    const tokens = ':root {\n  --ax-color-red: #f00;\n}\n';
    const consumer = '.foo { color: var(--ax-color-red); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-color-red: #f00;');
  });

  test('removes an unreferenced token', () => {
    const tokens = ':root {\n  --ax-color-red: #f00;\n  --ax-color-blue: #00f;\n}\n';
    const consumer = '.foo { color: var(--ax-color-red); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-color-red');
    expect(result).not.toContain('--ax-color-blue');
  });

  test('keeps transitively referenced tokens', () => {
    const tokens = [
      ':root {',
      '  --ax-surface: var(--ax-blue-500);',
      '  --ax-blue-500: #0060c0;',
      '  --ax-unused: #999;',
      '}',
    ].join('\n');
    const consumer = '.card { background: var(--ax-surface); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-surface');
    expect(result).toContain('--ax-blue-500');
    expect(result).not.toContain('--ax-unused');
  });

  test('keeps deeply transitive token chains', () => {
    const tokens = [':root {', '  --a: var(--b);', '  --b: var(--c);', '  --c: #123;', '  --d: #456;', '}'].join('\n');
    const consumer = '.x { color: var(--a); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--a');
    expect(result).toContain('--b');
    expect(result).toContain('--c');
    expect(result).not.toContain('--d');
  });

  test('handles multiple var references in a single value', () => {
    const tokens = [
      ':root {',
      '  --ax-shadow: 0 2px 4px var(--ax-shadow-color) var(--ax-shadow-alpha);',
      '  --ax-shadow-color: #000;',
      '  --ax-shadow-alpha: 0.2;',
      '  --ax-unused: red;',
      '}',
    ].join('\n');
    const consumer = '.box { box-shadow: var(--ax-shadow); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-shadow-color');
    expect(result).toContain('--ax-shadow-alpha');
    expect(result).not.toContain('--ax-unused');
  });

  test('handles multiple var references in consumer CSS', () => {
    const tokens = [
      ':root {',
      '  --ax-color-red: #f00;',
      '  --ax-color-blue: #00f;',
      '  --ax-color-green: #0f0;',
      '}',
    ].join('\n');
    const consumer = '.a { color: var(--ax-color-red); } .b { color: var(--ax-color-blue); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-color-red');
    expect(result).toContain('--ax-color-blue');
    expect(result).not.toContain('--ax-color-green');
  });

  test('removes empty rule blocks after stripping', () => {
    const tokens = [':root {', '  --ax-used: #f00;', '}', '.dark {', '  --ax-unused: #0f0;', '}'].join('\n');
    const consumer = '.x { color: var(--ax-used); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain(':root {');
    expect(result).not.toContain('.dark');
  });

  test('returns only a newline when no tokens are referenced', () => {
    const tokens = [':root {', '  --ax-color-red: #f00;', '  --ax-color-blue: #00f;', '}'].join('\n');
    const consumer = '.foo { color: red; }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).not.toContain('--ax-color-red');
    expect(result).not.toContain('--ax-color-blue');
  });

  test('preserves non-declaration lines like comments', () => {
    const tokens = ['/* Token definitions */', ':root {', '  --ax-used: #f00;', '  --ax-unused: #00f;', '}'].join('\n');
    const consumer = '.x { color: var(--ax-used); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('/* Token definitions */');
    expect(result).toContain('--ax-used');
    expect(result).not.toContain('--ax-unused');
  });

  test('does not leave excessive blank lines', () => {
    const tokens = [
      ':root {',
      '  --a: #111;',
      '  --b: #222;',
      '  --c: #333;',
      '  --d: #444;',
      '  --e: #555;',
      '}',
    ].join('\n');
    const consumer = '.x { color: var(--a); background: var(--e); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).not.toMatch(EXCESSIVE_NEWLINES_PATTERN);
  });

  test('result always ends with exactly one newline', () => {
    const tokens = ':root {\n  --a: #f00;\n}\n';
    const consumer = '.x { color: var(--a); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toMatch(SINGLE_TRAILING_NEWLINE_PATTERN);
  });

  test('handles var references with fallback values', () => {
    const tokens = [':root {', '  --ax-primary: #f00;', '  --ax-unused: #0f0;', '}'].join('\n');
    const consumer = '.x { color: var(--ax-primary, red); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--ax-primary');
    expect(result).not.toContain('--ax-unused');
  });

  test('handles circular references without infinite loop', () => {
    const tokens = [':root {', '  --a: var(--b);', '  --b: var(--a);', '}'].join('\n');
    const consumer = '.x { color: var(--a); }';

    const result = treeShakeTokens(tokens, consumer);

    expect(result).toContain('--a');
    expect(result).toContain('--b');
  });
});
