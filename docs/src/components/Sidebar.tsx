'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface NavItem {
  title: string;
  href: string;
}

const navigation: NavItem[] = [
  { title: 'Introduction', href: '/' },
  { title: 'Getting Started', href: '/getting-started' },
  { title: 'Installation', href: '/installation' },
  { title: 'Configuration', href: '/configuration' },
  { title: 'Commands', href: '/commands' },
  { title: 'Clients', href: '/clients' },
  { title: 'Validator Setup', href: '/validator' },
  { title: 'Requirements', href: '/requirements' },
  { title: 'Contributing', href: '/contributing' },
];

export default function Sidebar() {
  const pathname = usePathname();

  const isActive = (href: string) => {
    if (href === '/') return pathname === href;
    return pathname.startsWith(href);
  };

  return (
    <aside className="fixed top-0 left-0 h-screen w-64 bg-gray-50 border-r border-gray-200 overflow-y-auto">
      <div className="p-6">
        <Link href="/">
          <h1 className="text-xl font-bold text-gray-900 mb-1">
            starknode-kit
          </h1>
          <p className="text-sm text-gray-600">Documentation</p>
        </Link>
      </div>

      <nav className="px-4 pb-8">
        {navigation.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={`block px-3 py-2 mb-1 rounded-md text-sm font-medium transition-colors ${
              isActive(item.href)
                ? 'bg-blue-50 text-blue-600'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
          >
            {item.title}
          </Link>
        ))}
      </nav>
    </aside>
  );
}

