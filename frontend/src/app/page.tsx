import {SearchInput} from "@/components/SearchInput";
import {PackagesList} from "@/components/PackageList";

export default function Home() {
  return (
    <div className="py-0.5">
      <SearchInput updateUrl />
      <PackagesList packages={[
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
          id: 3,
          created_at: new Date().toISOString(),
        },
      ]} />
    </div>
  )
}
