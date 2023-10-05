import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import Link from "next/link";
import {Button} from "@/components/Button";
import {PlusIcon} from "@heroicons/react/20/solid";
import {NavMenu, NavMenuMobile, NavMenuProps} from "@/components/Navbar";
import {Suspense} from "react";
import {Container} from "@/components/Container";
import {LibConfig} from "@/components";

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Alin.io Package Store',
  description: 'Alin.io Package Store Open Source Version'
}

const MenuItems: NavMenuProps['items'] = [{
  title: 'Alin.io',
  href: LibConfig.Routes.Home,
  logo: true,
}, {
  title: 'Packages',
  href: LibConfig.Routes.Home,
}, {
  title: 'Documentation',
  href: 'https://docs.alin.io',
}];

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="h-full">
      <body className={inter.className}>
        <NavMenu
          ctaItems={(
            <Link href="/">
              <Button
                HeroIcon={PlusIcon}
              >
                New Package
              </Button>
            </Link>
          )}
          items={MenuItems}
        >
          <NavMenuMobile items={MenuItems} />
        </NavMenu>
        <div className="mt-16 mb-auto">
          <Container>
            <Suspense fallback={<div>Loading...</div>}>
              {children}
            </Suspense>
          </Container>
        </div>
      </body>
    </html>
  )
}
