import { Box, Heading, HStack, Tag, VStack } from '@navikt/ds-react';
import { StrictMode } from 'react';
import { PreviewCard } from '@/src/components/preview-card';
import { Language, type Translation } from '@/src/language';
import '@/src/styles/global.css';

interface PreviewProps {
  name: string;
  width: number;
  className: string;
}

const Preview = ({ name, width, className }: PreviewProps) => (
  <VStack as="article" gap="0" align="center" justify="center">
    <HStack asChild gap="2">
      <Heading level="1" size="xsmall" spacing>
        <span>{name}</span>

        <Tag size="small" variant="info-filled">
          {width}px
        </Tag>
      </Heading>
    </HStack>

    <HStack gap="4" align="stretch" justify="center">
      {Object.values(Language).map((lang) => (
        <VStack as="section" gap="2" key={lang}>
          <Tag size="small" variant="success" className="mx-auto w-fit">
            {LANGUAGE_LABELS[lang]}
          </Tag>

          <Box borderColor="border-default" borderWidth="4" flexGrow="1">
            <VStack>
              <Box background="bg-default" borderColor="border-default" borderWidth="0 0 1 0" padding="1" flexGrow="1">
                https://nav.no/minside/{lang}
              </Box>

              <Box asChild background="bg-subtle">
                <HStack
                  align="center"
                  justify="center"
                  flexGrow="1"
                  paddingInline="1"
                  paddingBlock="4"
                  className={className}
                >
                  <HStack
                    align="center"
                    justify="center"
                    flexGrow="1"
                    style={{ width, containerType: 'inline-size', containerName: 'preview' }}
                  >
                    <PreviewCard lang={lang} href="/" />
                  </HStack>
                </HStack>
              </Box>
            </VStack>
          </Box>
        </VStack>
      ))}
    </HStack>
  </VStack>
);

export const Dev = () => (
  <StrictMode>
    <VStack as="main" gap="24" padding="4" width="100%" height="100%">
      <Preview name="Mobil retningslinje" width={288} className="aspect-9/16" />

      <Preview name="Desktop retningslinje" width={468} className="aspect-9/16" />

      <Preview name="Desktop maks" width={912} className="aspect-video" />
    </VStack>
  </StrictMode>
);

const LANGUAGE_LABELS: Translation = {
  [Language.NB]: 'Bokmål',
  [Language.NN]: 'Nynorsk',
  [Language.EN]: 'English',
};
