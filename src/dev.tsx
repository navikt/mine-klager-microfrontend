import { PreviewCard } from '@app/card';
import { Language, type Translation } from '@app/language';
import { Heading, HStack, Tag, VStack } from '@navikt/ds-react';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { styled } from 'styled-components';

const container = document.getElementById('root');

if (!container) {
  throw new Error('No container found');
}

const root = createRoot(container);

interface PreviewProps {
  name: string;
  width: number;
}

const Preview = ({ name, width }: PreviewProps) => (
  <VStack gap="0" align="center" justify="center">
    <StyledHeading level="1" size="xsmall" spacing>
      <span>{name}</span>

      <Tag size="small" variant="info-filled">
        {width}px
      </Tag>
    </StyledHeading>

    <HStack gap="4" align="start" justify="center">
      {Object.values(Language).map((lang) => (
        <CardContainer key={lang} $width={width}>
          <StyledTag size="small" variant="success">
            {LANGUAGE_LABELS[lang]}
          </StyledTag>

          <PreviewCard lang={lang} href="/" />
        </CardContainer>
      ))}
    </HStack>
  </VStack>
);

const StyledHeading = styled(Heading)`
  display: flex;
  gap: 8px;
`;

const StyledTag = styled(Tag)`
  margin-left: auto;
  margin-right: auto;
  width: fit-content;
`;

interface VStackContainerProps {
  $width: number;
}

const CardContainer = styled.section<VStackContainerProps>`
  display: flex;
  flex-direction: column;
  container-type: inline-size;
  container-name: preview;
  width: ${({ $width }) => `${$width}px`};
  gap: 8px;
`;

const StyledMain = styled.main`
  display: flex;
  flex-direction: column;
  gap: 96px;
  width: 100%;
  height: 100%;
  padding: 16px;
`;

root.render(
  <StrictMode>
    <StyledMain>
      <Preview name="Mobil retningslinje" width={288} />

      <Preview name="Desktop retningslinje" width={468} />

      <Preview name="Desktop maks" width={912} />
    </StyledMain>
  </StrictMode>,
);

const LANGUAGE_LABELS: Translation = {
  [Language.NB]: 'Bokmål',
  [Language.NN]: 'Nynorsk',
  [Language.EN]: 'English',
};
