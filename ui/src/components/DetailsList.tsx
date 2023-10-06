import { ReactNode } from 'react';
import { classNames } from '.';

interface Props {
  details: {
    name: string;
    value: ReactNode | ReactNode[];
    description?: string;
  }[];
  className?: string;
}

export function DetailsList({ details, className }: Props) {
  return (
    <div className={classNames(className ?? '')}>
      <dl className="divide-y">
        {details.map(({ name, value, description }) => (
          <div key={name} className="px-4 py-6 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6">
            <dt className="text-sm font-medium text-gray-900">{name}</dt>
            <dd className="mt-1 text-sm leading-6 text-gray-700 sm:col-span-2 sm:mt-0">
              {value}
              {description && <p className="text-xs text-gray-500 block">{description}</p>}
            </dd>
          </div>
        ))}
      </dl>
    </div>
  );
}
