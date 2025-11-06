import { CardContainer, CardContent, type CardProps } from '@/src/components/card';
import styles from '@/src/components/preview-card.module.css';

export const PreviewCard = ({ lang, href }: CardProps) => (
  <CardContainer href={href} className={styles.container}>
    <CardContent lang={lang} />
  </CardContainer>
);
