import { PackageCodeSample } from './PackageCodeSample.tsx';

export function PackagesEmptyState() {
  return (
    <div className="text-center">
      <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path
          vectorEffect="non-scaling-stroke"
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z"
        />
      </svg>
      <h3 className="mt-2 text-sm font-semibold text-gray-900">Publish and Store private packages!</h3>
      <p className="mt-1 text-sm text-gray-500">
        Alin.io Package Store allows to self-host package registry for NPM, Pypi and Docker Container Images
      </p>
      <div className="mb-16">
        <div className="my-4 flex flex-col items-center justify-center">
          <h3 className="my-2 text-sm font-semibold text-gray-900">Publish and Store private packages!</h3>
          <PackageCodeSample name="test" service="container" type="push" />
        </div>
        <div className="my-4 flex flex-col items-center justify-center">
          <h3 className="my-2 text-sm font-semibold text-gray-900">Publish and Store private packages!</h3>
          <PackageCodeSample name="test" service="npm" type="push" />
        </div>
        <div className="my-4 flex flex-col items-center justify-center">
          <h3 className="my-2 text-sm font-semibold text-gray-900">Publish and Store private packages!</h3>
          <PackageCodeSample name="test" service="pypi" type="push" />
        </div>
      </div>
    </div>
  );
}
