import { SearchInput } from '../components/SearchInput.tsx';
import { PackagesList } from '../components/PackageList.tsx';
import { useGetPackages } from '../api';
import { TopLoadingBar } from '../components/TopLoadingBar.tsx';
import { Alert } from '../components/Alert.tsx';

export function PackagesPage() {
  const { data, error } = useGetPackages('');
  const isLoading = !data && !error;

  return (
    <div className="py-0.5">
      <SearchInput updateUrl />
      {isLoading && <TopLoadingBar />}
      {data && <PackagesList packages={data} />}
      {error && <Alert title="Unable to fetch packages" message={error.message} variant="error" />}
    </div>
  );
}
