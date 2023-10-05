import { ReactNode } from 'react';
import { classNames } from '.';

interface Props {
  children: ReactNode;
  className?: string;
}

export function Container({ children, className }: Props) {
  return <div className={classNames('mx-auto px-4 py-4 max-w-7xl sm:px-6 lg:px-8', className ?? '')}>{children}</div>;
}
