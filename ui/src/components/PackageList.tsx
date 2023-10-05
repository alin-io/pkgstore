import { PackagesEmptyState } from './PackagesEmptyState';
import { Link } from 'react-router-dom';
import { formatDistance } from 'date-fns';
import ChevronRightIcon from '@heroicons/react/20/solid/ChevronRightIcon';
import { LibConfig, PackageItem } from '.';
import { ServiceIcon } from './ServiceIcon';

interface Props {
  packages: PackageItem[];
}

export function PackagesList({ packages }: Props) {
  return (
    <>
      {packages.length === 0 ? (
        <div className="mt-4">
          <PackagesEmptyState />
        </div>
      ) : (
        <ul role="list" className="divide-y divide-gray-100">
          {packages.map((item) => (
            <li key={item.id} className="relative flex justify-between gap-x-6 py-5">
              <div className="flex items-center gap-x-4">
                <ServiceIcon name={item.service} className="h-10 w-10 flex-none rounded-full fill-gray-500" />
                <div className="min-w-0 flex-auto">
                  <p className="text-sm font-semibold leading-6 text-gray-900">
                    <Link to={LibConfig.Routes.Package(item.id)}>
                      <span className="absolute inset-x-0 -top-px bottom-0" />
                      {item.name}
                    </Link>
                  </p>
                  <p className="mt-1 flex text-xs leading-5 text-gray-500">{item.name}</p>
                </div>
              </div>
              <div className="flex items-center gap-x-4">
                <div className="hidden sm:flex sm:flex-col sm:items-end">
                  <p className="text-sm leading-6 text-gray-900 inline-flex items-center">
                    <ServiceIcon name={item.service} className="-ml-0.5 mr-2 h-6 w-6" />
                    {item.service}
                  </p>
                  {item.updated_at ? (
                    <p className="mt-1 text-xs leading-5 text-gray-500">
                      Last Updated <time dateTime={item.updated_at.toString()}>{formatDistance(new Date(item.updated_at), new Date())}</time>
                    </p>
                  ) : (
                    <div className="mt-1 flex items-center gap-x-1.5">
                      <div className="flex-none rounded-full bg-emerald-500/20 p-1">
                        <div className="h-1.5 w-1.5 rounded-full bg-emerald-500" />
                      </div>
                      <p className="text-xs leading-5 text-gray-500">Online</p>
                    </div>
                  )}
                </div>
                <ChevronRightIcon className="h-5 w-5 flex-none text-gray-400" aria-hidden="true" />
              </div>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
