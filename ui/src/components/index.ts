export const LibConfig = {
  urlPrefix: '/ui',
  Routes: {
    Home: '/',
    Packages: '/packages',
    Package: (id: string) => `/packages/${id}`,
  },
};

export function classNames(...classes: string[]) {
  return classes.filter(Boolean).join(' ');
}
