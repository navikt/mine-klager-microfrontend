import { Card } from '@/src/components/card';
import { DEFAULT_LANGUAGE, type Language } from '@/src/language';
import { isProduction } from '@/src/utils/environment';

interface Props {
  language: Language;
}

export const Microfrontend = ({ language }: Props) => <Card lang={language} href={getMineKlagerUrl(language)} />;

const getMineKlagerUrl = (lang: Language): string => {
  const url = isProduction ? 'https://mine-klager.nav.no' : 'https://mine-klager.intern.dev.nav.no';

  if (lang === DEFAULT_LANGUAGE) {
    return url;
  }

  return `${url}/${lang}`;
};
