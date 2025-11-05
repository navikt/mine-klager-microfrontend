import { Card } from '@app/card';
import { DEFAULT_LANGUAGE, type Language, useLanguage } from '@app/language';

export const Microfrontend = () => {
  const language = useLanguage();

  return <Card lang={language} href={getMineKlagerUrl(language)} />;
};

export default Microfrontend;

const getMineKlagerUrl = (lang: Language): string => {
  const { hostname } = window.location;
  const isProduction = hostname === 'nav.no' || hostname === 'www.nav.no';

  const url = isProduction ? 'https://mine-klager.nav.no' : 'https://mine-klager.intern.dev.nav.no';

  if (lang === DEFAULT_LANGUAGE) {
    return url;
  }

  return `${url}/${lang}`;
};
