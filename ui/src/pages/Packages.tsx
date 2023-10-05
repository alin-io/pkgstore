import { SearchInput } from '../components/SearchInput.tsx';
import { PackagesList } from '../components/PackageList.tsx';

export function PackagesPage() {
  return (
    <div className="py-0.5">
      <SearchInput updateUrl />
      <PackagesList
        packages={[
          {
            service: 'npm',
            updated_at: new Date().toISOString(),
            name: 'react',
            id: 1,
            created_at: new Date().toISOString(),
          },
          {
            service: 'pypi',
            updated_at: new Date().toISOString(),
            name: 'requests',
            id: 3,
            created_at: new Date().toISOString(),
          },
          {
            service: 'container',
            updated_at: new Date().toISOString(),
            name: 'ubuntu',
            id: 4,
            created_at: new Date().toISOString(),
          },
        ]}
      />
    </div>
  );
}
