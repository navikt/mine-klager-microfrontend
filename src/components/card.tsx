import { ChevronRightIcon } from '@navikt/aksel-icons';
import { BodyLong, Box, Heading, HStack, VStack } from '@navikt/ds-react';
import { CardIcon } from '@/src/components/icon';
import type { Translation } from '@/src/language';
import { Language } from '@/src/language';

export interface CardProps {
  lang: Language;
  href: string;
}

export const Card = ({ lang, href }: CardProps) => (
  <CardContainer href={href} className="mine-klager-microfrontend">
    <CardContent lang={lang} />
  </CardContainer>
);

interface CardContainerProps {
  href: string;
  children: React.ReactNode;
  className?: string;
}

export const CardContainer = ({ href, children, className }: CardContainerProps) => (
  <Box
    asChild
    borderRadius="large"
    background="bg-default"
    shadow="xsmall"
    padding={{ xs: '4', sm: '4', md: '5', lg: '5', xl: '5', '2xl': '5' }}
    className={`hover:shadow-md text-neutral-950 ${className}`}
  >
    <HStack as="a" align="center" wrap={false} href={href} className="no-underline cursor-pointer group/container">
      {children}
    </HStack>
  </Box>
);

interface CardContentProps {
  lang: Language;
}

export const CardContent = ({ lang }: CardContentProps) => (
  <>
    <HStack gap="5" align="center" flexGrow="1" wrap={false}>
      <CardIcon title={HEADING[lang]} />

      <VStack flexShrink="1">
        <Heading level="3" size="small" className="group-hover/container:underline">
          {HEADING[lang]}
        </Heading>

        <BodyLong size="medium">{DESCRIPTION[lang]}</BodyLong>
      </VStack>
    </HStack>

    <ChevronRightIcon className="h-full w-6 shrink-0 transition-transform ease-in-out duration-200 group-hover/container:translate-x-[3px]" />
  </>
);

const HEADING: Translation = {
  [Language.NB]: 'Mine saker hos Klageinstans',
  [Language.NN]: 'Mine saker hjå Klageinstans',
  [Language.EN]: 'My cases with Nav Complaints Unit',
};

const DESCRIPTION: Translation = {
  [Language.NB]: 'Her kan du se status på dine saker hos Klageinstans.',
  [Language.NN]: 'Her kan du sjå status på dine saker hjå Klageinstans.',
  [Language.EN]: 'Here you can see the status of your cases with Nav Complaints Unit (Klageinstans).',
};
