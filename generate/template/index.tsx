import { LinkCard } from '@navikt/ds-react';
import { renderToStaticMarkup } from 'react-dom/server';
import { Icon } from '@/icon';

const URL_PLACEHOLDER = '{{URL}}';
const TITLE_PLACEHOLDER = '{{TITLE}}';
const DESCRIPTION_PLACEHOLDER = '{{DESCRIPTION}}';

const html = renderToStaticMarkup(
  <LinkCard data-color="neutral">
    <LinkCard.Icon>
      <Icon />
    </LinkCard.Icon>
    <LinkCard.Title as="h2">
      <LinkCard.Anchor href={URL_PLACEHOLDER}>{TITLE_PLACEHOLDER}</LinkCard.Anchor>
    </LinkCard.Title>
    <LinkCard.Description>{DESCRIPTION_PLACEHOLDER}</LinkCard.Description>
  </LinkCard>,
);

await Bun.write('templates/template.html', `${html}\n`);
