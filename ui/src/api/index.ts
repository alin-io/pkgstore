import { Package, PackageVersion } from './types';
import { useFetch } from '../hooks/useFetch.ts';

export * from './types';

export const SERVER_HOST = import.meta.env.VITE_SERVER_HOST || 'http://localhost:8080';
const API_URL = `${SERVER_HOST}/api`;
const FetchOptions: RequestInit = {
  credentials: 'include',
  headers: {
    Authorization: `Bearer username:secret`, // TODO: replace with real token
  },
};

export function useGetPackages(q?: string) {
  return useFetch<Package[]>(`${API_URL}/packages?q=${q ?? ''}`, FetchOptions);
}

export function useGetPackage(id: string) {
  return useFetch<Package>(`${API_URL}/packages/${id}`, FetchOptions);
}

export function useGetPackageVersions(id: string) {
  return useFetch<PackageVersion[]>(`${API_URL}/packages/${id}/versions`, FetchOptions);
}
