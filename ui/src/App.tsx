import { BrowserRouter, Link, Route, Routes } from 'react-router-dom';
import { LibConfig } from './components';
import { Container } from './components/Container.tsx';
import { PackagesPage } from './pages/Packages.tsx';
import { NavMenu, NavMenuMobile } from './components/Navbar';
import { Button } from './components/Button.tsx';
import PlusIcon from '@heroicons/react/20/solid/PlusIcon';
import { NavMenuProps } from './components/Navbar/NavMenu.tsx';

const MenuItems: NavMenuProps['items'] = [
  {
    title: 'Alin.io',
    href: LibConfig.Routes.Home,
    logo: true,
  },
  {
    title: 'Packages',
    href: LibConfig.Routes.Home,
  },
  {
    title: 'Documentation',
    href: 'https://docs.alin.io',
  },
];

export default function App() {
  return (
    <BrowserRouter basename={LibConfig.urlPrefix}>
      <NavMenu
        ctaItems={
          <Link to="/">
            <Button HeroIcon={PlusIcon}>New Package</Button>
          </Link>
        }
        items={MenuItems}
      >
        <NavMenuMobile items={MenuItems} />
      </NavMenu>
      <div className="mt-16 mb-auto">
        <Container>
          <Routes>
            <Route path="/" element={<PackagesPage />} />
          </Routes>
        </Container>
      </div>
    </BrowserRouter>
  );
}
