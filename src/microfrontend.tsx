import { LinkCard } from '@navikt/ds-react';
import { Icon } from '@/icon';

export interface MicrofrontendProps {
  /** The title text for the LinkCard */
  title: string;
  /** The description text for the LinkCard */
  description: string;
  /** The user facing URL for the LinkCard */
  url: string;
}

export const Microfrontend = ({ title: heading, description, url }: MicrofrontendProps) => (
  <LinkCard size="medium">
    <LinkCard.Icon>
      <Icon />
    </LinkCard.Icon>
    <LinkCard.Title as="h2">
      <LinkCard.Anchor href={url}>{heading}</LinkCard.Anchor>
    </LinkCard.Title>
    <LinkCard.Description>{description}</LinkCard.Description>
  </LinkCard>
);
