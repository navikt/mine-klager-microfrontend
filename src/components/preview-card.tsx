import { CardContainer, CardContent, type CardProps } from '@/src/components/card';

export const PreviewCard = ({ lang, href }: CardProps) => (
  <CardContainer href={href} className="container/preview max-w-[647px]">
    <CardContent lang={lang} />
  </CardContainer>
);
