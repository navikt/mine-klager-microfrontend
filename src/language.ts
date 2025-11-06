export enum Language {
  NB = 'nb',
  NN = 'nn',
  EN = 'en',
}

export type Translation = Record<Language, string>;

const LANGUAGES = Object.values(Language);
export const DEFAULT_LANGUAGE = Language.NB;

export const isLanguage = (value: string): value is Language => LANGUAGES.some((l) => l === value);
