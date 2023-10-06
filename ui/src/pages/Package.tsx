import { useParams } from 'react-router-dom';
import { DetailsList } from '../components/DetailsList.tsx';
import { useGetPackage, useGetPackageVersions } from '../api';
import { TopLoadingBar } from '../components/TopLoadingBar.tsx';
import { Alert } from '../components/Alert.tsx';
import { formatDistance } from 'date-fns';
import { ServiceIcon } from '../components/ServiceIcon.tsx';
import { CheckCircleIcon } from '@heroicons/react/20/solid';
import { classNames } from '../components';
import { PackageCodeSample } from '../components/PackageCodeSample.tsx';

export function PackagePage() {
  const { id } = useParams();
  const { data: pkg, error: pkgError } = useGetPackage(id ?? '');
  const { data: versions, error: versionsError } = useGetPackageVersions(id ?? '');
  const isLoading = !pkg && !pkgError && !versions && !versionsError;
  const error = pkgError ?? versionsError;

  return (
    <div className="py-0.5">
      {isLoading && <TopLoadingBar />}
      {error && <Alert title="Unable to fetch packages" message={error.message} variant="error" />}
      {pkg && (
        <>
          <h1 className="text-2xl font-bold flex mb-6">
            <ServiceIcon name={pkg.service} className="mr-3 text-4xl" />
            {pkg.name}
          </h1>
          <PackageCodeSample name={pkg.name} service={pkg.service} type="pull" />
          <DetailsList
            className="border-t border-gray-100"
            details={[
              { name: 'Name', value: pkg.name },
              { name: 'Service', value: pkg.service },
              { name: 'Latest Version', value: pkg.latest_version },
              {
                name: 'Last Updated',
                value: <time dateTime={pkg.updated_at.toString()}>{formatDistance(new Date(pkg.updated_at), new Date())}</time>,
              },
            ]}
          />
        </>
      )}
      {versions && (
        <div className="mt-10">
          <h4 className="text-2xl">Versions</h4>
          <ul role="list" className="space-y-6 mt-6">
            {versions.map(({ version, created_at, digest, tag }, activityItemIdx) => (
              <li key={version} className="relative flex gap-x-4">
                <div
                  className={classNames(
                    activityItemIdx === versions.length - 1 ? 'h-6' : '-bottom-6',
                    'absolute left-0 top-0 flex w-6 justify-center'
                  )}
                >
                  <div className="w-px bg-gray-200" />
                </div>
                <div className="relative flex h-6 w-6 flex-none items-center justify-center bg-white">
                  {tag === 'latest' ? (
                    <CheckCircleIcon className="h-6 w-6 text-indigo-600" aria-hidden="true" />
                  ) : (
                    <div className="h-1.5 w-1.5 rounded-full bg-gray-100 ring-1 ring-gray-300" />
                  )}
                </div>
                <p className="flex-auto py-0.5 text-sm leading-5 text-gray-500">
                  <span className="font-medium text-gray-900">{version} - </span> Digest <code>[{digest}]</code>
                </p>
                <time dateTime={new Date(created_at).toISOString()} className="flex-none py-0.5 text-xs leading-5 text-gray-500">
                  {formatDistance(new Date(created_at), new Date(), { addSuffix: true })}
                </time>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
