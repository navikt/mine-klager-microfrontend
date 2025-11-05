import { CardIcon } from '@app/icons/icon';
import type { Translation } from '@app/language';
import { Language } from '@app/language';
import { logNavigereEvent } from '@app/utils/amplitude';
import { ChevronRightIcon } from '@navikt/aksel-icons';
import { BodyLong, Heading, HStack, VStack } from '@navikt/ds-react';
import styled from 'styled-components';

interface CardProps {
  lang: Language;
  href: string;
}

export const PreviewCard = ({ lang, href }: CardProps) => (
  <PreviewContainer href={href}>
    <CardContent lang={lang} />
  </PreviewContainer>
);

const StyledChevron = styled(ChevronRightIcon)`
  height: 100%;
  width: 24px;
  flex-shrink: 0;
  transition: transform 0.2s ease-in-out;
`;

const ContainerStack = styled.a`
  display: flex;
  flex-wrap: nowrap;
  align-items: center;
  text-decoration: none;
  cursor: pointer;
  border-radius: var(--a-border-radius-large);
  background-color: var(--a-bg-default);
  color: var(--a-text-default);
  box-shadow: var(--a-shadow-xsmall);
  padding: 20px;

  &:hover {
    box-shadow: var(--a-shadow-small);

    .mine-klager-card-heading {
      text-decoration: underline;
    }

    ${StyledChevron} {
      transform: translateX(3px);
    }
  }
  
  @media (max-width: 647px) {
    padding: 16px;
  }
`;

const PreviewContainer = styled(ContainerStack)`    
  @container preview (max-width: 647px) {
    padding: 16px;
  }
`;

export const Card = ({ lang, href }: CardProps) => (
  <Container href={href}>
    <CardContent lang={lang} />
  </Container>
);

interface ContainerProps {
  className?: string;
  href: string;
  children: React.ReactNode;
}

const Container = ({ className = '', href, children }: ContainerProps) => (
  <ContainerStack href={href} className={className} onClick={() => logNavigereEvent()}>
    {children}
  </ContainerStack>
);

interface CardContentProps {
  lang: Language;
}

const CardContent = ({ lang }: CardContentProps) => (
  <>
    <Content gap="5" align="center" wrap={false}>
      <CardIcon />

      <InnerContent>
        <Heading level="3" size="small" className="mine-klager-card-heading">
          {HEADING[lang]}
        </Heading>

        <BodyLong size="medium">{DESCRIPTION[lang]}</BodyLong>
      </InnerContent>
    </Content>

    <StyledChevron />
  </>
);

const Content = styled(HStack)`
  flex-grow: 1;
`;

const InnerContent = styled(VStack)`
  flex-shrink: 1;
`;

const HEADING: Translation = {
  [Language.NB]: 'Mine klager og anker',
  [Language.NN]: 'Mine klager og anker',
  [Language.EN]: 'My complaints and appeals',
};

const DESCRIPTION: Translation = {
  [Language.NB]: 'Her kan du se status på dine klager og anker hos klageinstansen.',
  [Language.NN]: 'Her kan du sjå status på dine klager og anker hjå klageinstansen.',
  [Language.EN]:
    'Here you can see the status of your complaints and appeals with Nav Complaints Unit (Nav klageinstans).',
};
