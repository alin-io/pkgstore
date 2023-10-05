'use client'

import {Disclosure} from "@headlessui/react";
import Link from "next/link";
import {NavMenuProps} from "@/components/Navbar/NavMenu";
import {usePathname} from "next/navigation";
import {useMemo} from "react";

export function NavMenuMobile({ items, userMenu }: NavMenuProps) {
  const pathname = usePathname();

  const menuItems = useMemo(() => items.filter(item => !item.logo), [items]);

  return (
    <Disclosure.Panel className="md:hidden">
      <div className="space-y-1 pb-3 pt-2">
        {menuItems.map(({title, href}) => (
          <Disclosure.Button
            key={title}
            as={Link}
            href={href}
            className={`block ${href === pathname ? 'border-l-4 border-indigo-500 bg-indigo-50' : ''} py-2 pl-3 pr-4 text-base font-medium text-indigo-700 sm:pl-5 sm:pr-6`}
          >
            {title}
          </Disclosure.Button>
        ))}
      </div>
      {userMenu}
    </Disclosure.Panel>
  )
}
