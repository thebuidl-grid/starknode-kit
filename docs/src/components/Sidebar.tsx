'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';

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

interface SidebarProps {
  isOpen: boolean;
  onClose: () => void;
}

export default function Sidebar({ isOpen, onClose }: SidebarProps) {
  const pathname = usePathname();

  const isActive = (href: string) => {
    if (href === '/') return pathname === href;
    return pathname.startsWith(href);
  };

  // Close sidebar when route changes on mobile
  useEffect(() => {
    onClose();
  }, [pathname, onClose]);

  // Prevent body scroll when mobile menu is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = 'unset';
    }
    return () => {
      document.body.style.overflow = 'unset';
    };
  }, [isOpen]);

  return (
    <>
      {/* Backdrop for mobile */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black/30 bg-opacity-50 z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Sidebar */}
      <aside
        className={`fixed top-0 left-0 h-screen w-[20rem] md:w-64 bg-gray-50 dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 overflow-y-auto z-50 transition-all duration-300 ease-in-out ${
          isOpen ? 'translate-x-0' : '-translate-x-full'
        } lg:translate-x-0`}
      >
        <div className="p-6">
          <div className="flex items-center justify-between">
            <Link href="/" className="flex-1">
              <h1 className="text-xl font-bold text-gray-900 dark:text-white mb-1">
                starknode-kit
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">Documentation</p>
            </Link>
            {/* Close button for mobile */}
            <button
              onClick={onClose}
              className="lg:hidden p-2 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white"
              aria-label="Close menu"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>
        </div>

        <nav className="px-4 pb-8">
          {navigation.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className={`block px-3 py-2 mb-1 rounded-md text-sm font-medium transition-colors ${
                isActive(item.href)
                  ? 'bg-blue-50 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
              }`}
            >
              {item.title}
            </Link>
          ))}
        </nav>
      </aside>
    </>
  );
}

