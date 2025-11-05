import { styled } from 'styled-components';

interface IconProps {
  title: string;
  className?: string;
}

const Icon = ({ title, className }: IconProps) => (
  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 64 64" className={className}>
    <title>{title}</title>

    <path
      fill="#D8F9FF"
      fillRule="evenodd"
      d="M60 20H17v-8h43v8Zm-8 18H9v-8h43v8ZM0 56h43v-8H0v8Z"
      clipRule="evenodd"
    />
    <path
      fill="#23262A"
      fillRule="evenodd"
      d="M2 26V2h31v24H2ZM0 1l1-1h32l2 1v26l-2 1H1l-1-1V1Zm5 6h1l1 1 3-3a1 1 0 1 1 2 2l-5 3a1 1 0 0 1-1 0L5 8V7Zm9 13h2v1l4-3a1 1 0 0 1 1 2l-5 3a1 1 0 0 1-1 0l-1-2v-1Zm1-13a1 1 0 1 0 0 2h4a1 1 0 1 0 0-2h-4ZM4 15l1-1h4a1 1 0 1 1 0 2H5l-1-1Zm1 5a1 1 0 1 0 0 2h4a1 1 0 1 0 0-2H5Zm9-5 1-1h4a1 1 0 1 1 0 2h-4l-1-1Zm11-1a1 1 0 1 0 0 2h4a1 1 0 1 0 0-2h-4Zm-1-6 1-1h4a1 1 0 1 1 0 2h-4l-1-1Zm11 33h26v1l-1 4-15 5-12-5 2-4v-1Zm-2 7-3 12v1h26l3-13-14 5h-1l-11-5Zm0-7 2-2h26c1 0 3 2 2 4l-5 18c0 2-1 2-3 2H30c-2 0-3-1-2-3l5-19Z"
      clipRule="evenodd"
    />
  </svg>
);

export const CardIcon = styled(Icon)`
  width: 56px;
  flex-shrink: 0;
`;
