export interface Package {
  id: string;
  name: string;
  service: string;
  auth_id: string;
  latest_version: string;
  versions: PackageVersion[];
  created_at: string;
  updated_at: string;
}

export interface PackageVersion {
  id: string;
  service: string;
  digest: string;
  package_id: string;
  version: string;
  tag: string;
  metadata: unknown;
  created_at: string;
  updated_at: string;
}
