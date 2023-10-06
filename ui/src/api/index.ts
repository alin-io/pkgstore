import { Package, PackageVersion } from './types';
import { useFetch } from '../hooks/useFetch.ts';

export * from './types';

const API_URL = import.meta.env.PROD ? '/api' : 'http://localhost:8080/api';
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
