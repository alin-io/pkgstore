import MagnifyingGlassIcon from '@heroicons/react/20/solid/MagnifyingGlassIcon';
import { classNames } from '.';
import { ChangeEvent, useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';

interface Props extends React.DetailedHTMLProps<React.InputHTMLAttributes<HTMLInputElement>, HTMLInputElement> {
  updateUrl?: boolean;
}

export function SearchInput({ className, updateUrl, ...props }: Props) {
  const [value, setValue] = useState<string>((props.value ?? '').toString());
  const [searchParams, setSearchParams] = useSearchParams();

  useEffect(() => {
    setValue((props.value ?? '').toString());
  }, [props.value]);

  const setDataValue = (e: ChangeEvent<HTMLInputElement>) => {
    setValue(e.target.value);
    if (updateUrl) {
      const params = new URLSearchParams(searchParams as unknown as URLSearchParams);
      if (e.target.value.length > 0) {
        params.set('q', e.target.value);
      } else {
        params.delete('q');
      }
      setSearchParams(params);
    }
  };

  return (
    <div className="relative my-2 flex items-center">
      <input
        type="text"
        placeholder="Search in your Apps"
        className={classNames(
          'block w-full rounded-md border-0 py-1.5 pr-14 pl-4 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6',
          className ?? ''
        )}
        value={value}
        onChange={setDataValue}
        {...props}
      />
      <div className="absolute inset-y-0 right-0 flex py-1.5 pr-1.5">
        <kbd className="inline-flex items-center rounded border border-gray-200 px-1 font-sans text-xs text-gray-400">
          <MagnifyingGlassIcon className="h-5 w-5" />
        </kbd>
      </div>
    </div>
  );
}
