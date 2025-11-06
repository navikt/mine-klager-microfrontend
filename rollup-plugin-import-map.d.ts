import type { Plugin } from 'rollup';

declare module 'rollup-plugin-import-map' {
  interface ImportMap {
    imports: Record<string, string>;
  }

  type ImportMapPath = string;

  type Maps = ImportMap | ImportMapPath;

  export function rollupImportMapPlugin(maps: Maps | Maps[]): Plugin;
}
