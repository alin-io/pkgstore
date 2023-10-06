import { ReactNode } from 'react';
import { classNames } from './index.ts';

interface Props {
  title: string | ReactNode;
  message: string | ReactNode;
  variant?: 'error' | 'success';
}

export function Alert({ title, message, variant }: Props) {
  return (
    <div className={classNames('rounded-md p-4', variant === 'error' ? 'bg-red-50' : 'bg-green-50')}>
      <div className="flex">
        <div className="ml-3">
          <h3 className={classNames('text-sm font-medium', variant === 'error' ? 'text-red-800' : 'text-green-800')}>{title}</h3>
          <div className={classNames('mt-2 text-sm', variant === 'error' ? 'text-red-700' : 'text-green-700')}>{message}</div>
        </div>
      </div>
    </div>
  );
}
