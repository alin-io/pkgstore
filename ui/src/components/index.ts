export interface PackageItem {
  id: number;
  service: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface PackageVersion {
  id: number;
  package_id: number;
  version: string;
  description: string;
  readme: string;
  created_at: string;
  updated_at: string;
}

export const LibConfig = {
  urlPrefix: '/ui',
  Routes: {
    Home: '/',
    Packages: '/packages',
    Package: (id: number) => `/packages/${id}`,
  },
};

export function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}
