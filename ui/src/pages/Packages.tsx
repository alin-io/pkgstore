import { SearchInput } from '../components/SearchInput.tsx';
import { PackagesList } from '../components/PackageList.tsx';
import { useGetPackages } from '../api';
import { TopLoadingBar } from '../components/TopLoadingBar.tsx';
import { Alert } from '../components/Alert.tsx';
import { useSearchParams } from 'react-router-dom';

export function PackagesPage() {
  const [searchParams] = useSearchParams();
  const q = searchParams.get('q');

  const { data, error } = useGetPackages(q ?? '');
  const isLoading = !data && !error;

  return (
    <div className="py-0.5">
      <SearchInput updateUrl value={q ?? ''} />
      {isLoading && <TopLoadingBar />}
      {data && <PackagesList packages={data} />}
      {error && <Alert title="Unable to fetch packages" message={error.message} variant="error" />}
    </div>
  );
}
