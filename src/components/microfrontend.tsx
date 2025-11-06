import { Card } from '@/src/components/card';
import { DEFAULT_LANGUAGE, type Language } from '@/src/language';
import { isDeployedToProd } from '@/src/utils/environment';

interface Props {
  language: Language;
}

export const Microfrontend = ({ language }: Props) => <Card lang={language} href={getMineKlagerUrl(language)} />;

const getMineKlagerUrl = (lang: Language): string => {
  const url = isDeployedToProd ? 'https://mine-klager.nav.no' : 'https://mine-klager.ansatt.dev.nav.no';

  return lang === DEFAULT_LANGUAGE ? url : `${url}/${lang}`;
};
