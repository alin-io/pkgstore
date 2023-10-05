import { Disclosure } from '@headlessui/react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/24/outline';
import { ReactNode, useMemo } from 'react';
import { Link } from 'react-router-dom';
import { useLocation } from 'react-router-dom';

export interface NavMenuProps {
  items: {
    logo?: string | boolean;
    title: string;
    href: string;
  }[];
  userMenu?: ReactNode | ReactNode[];
  ctaItems?: ReactNode | ReactNode[];
  children?: ReactNode | ReactNode[];
}

export const NavMenu = ({ items, userMenu, children, ctaItems }: NavMenuProps) => {
  const { pathname } = useLocation();

  const logoItem = useMemo(() => items.find((item) => !!item.logo), [items]);
  const menuItems = useMemo(() => items.filter((item) => !item.logo), [items]);

  return (
    <Disclosure as="nav" className="bg-white shadow fixed top-0 w-full z-40">
      {({ open }) => (
        <>
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div className="flex h-16 justify-between">
              <div className="flex">
                <div className="-ml-2 mr-2 flex items-center md:hidden">
                  <Disclosure.Button className="inline-flex items-center justify-center rounded-md p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-indigo-500">
                    <span className="sr-only">Open main menu</span>
                    {open ? <XMarkIcon className="block h-6 w-6" aria-hidden="true" /> : <Bars3Icon className="block h-6 w-6" aria-hidden="true" />}
                  </Disclosure.Button>
                </div>
                {logoItem && (
                  <div className="flex flex-shrink-0 items-center">
                    <Link to={logoItem.href}>
                      {typeof logoItem.logo === 'string' ? (
                        <img className="h-10 w-auto" width={144} height={40} src={logoItem.logo} alt={logoItem.title} />
                      ) : (
                        <span className="text-2xl">Alin.io</span>
                      )}
                    </Link>
                  </div>
                )}
                <div className="hidden md:ml-6 md:flex md:space-x-8">
                  {menuItems.map(({ title, href }) => (
                    <Link
                      key={title}
                      to={href}
                      className={`inline-flex items-center ${
                        href === pathname ? 'border-b-2 border-indigo-500' : ''
                      } px-1 pt-1 text-sm font-medium text-gray-900`}
                    >
                      {title}
                    </Link>
                  ))}
                </div>
              </div>
              <div className="flex items-center">
                <div className="flex-shrink-0">{ctaItems}</div>

                {userMenu}
              </div>
            </div>
          </div>

          {children}
        </>
      )}
    </Disclosure>
  );
};
