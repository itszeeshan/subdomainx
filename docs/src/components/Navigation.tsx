'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState } from 'react'
import { 
  Bars3Icon, 
  XMarkIcon,
  CommandLineIcon,
  BookOpenIcon,
  WrenchScrewdriverIcon,
  CogIcon,
  PlayIcon,
  InformationCircleIcon
} from '@heroicons/react/24/outline'
import { ThemeToggle } from './ThemeToggle'

const navigation = [
  { name: 'Home', href: '/', icon: CommandLineIcon },
  { name: 'Installation', href: '/installation', icon: BookOpenIcon },
  { name: 'CLI Reference', href: '/cli-reference', icon: WrenchScrewdriverIcon },
  { name: 'Supported Tools', href: '/supported-tools', icon: CogIcon },
  { name: 'Examples', href: '/examples', icon: PlayIcon },
  { name: 'Configuration', href: '/configuration', icon: InformationCircleIcon },
]

export default function Navigation() {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const pathname = usePathname()

  return (
    <>
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-gray-900/80" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 z-50 w-64 bg-white dark:bg-gray-900 shadow-xl">
          <div className="flex h-full flex-col">
            <div className="flex h-16 items-center justify-between px-6 border-b border-gray-200 dark:border-gray-700">
              <div className="flex items-center">
                <CommandLineIcon className="h-8 w-8 text-emerald-600" />
                <span className="ml-2 text-lg font-semibold text-gray-900 dark:text-white">SubdomainX</span>
              </div>
              <button
                onClick={() => setSidebarOpen(false)}
                className="p-2 rounded-lg text-gray-500 hover:text-gray-700 hover:bg-gray-100 dark:text-gray-400 dark:hover:text-gray-200 dark:hover:bg-gray-800"
              >
                <XMarkIcon className="h-6 w-6" />
              </button>
            </div>
            <nav className="flex-1 space-y-1 px-4 py-4">
              {navigation.map((item) => {
                const isActive = pathname === item.href
                return (
                  <Link
                    key={item.name}
                    href={item.href}
                    className={`group flex items-center px-3 py-2 text-sm font-medium rounded-lg transition-colors ${
                      isActive
                        ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800'
                        : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white'
                    }`}
                    onClick={() => setSidebarOpen(false)}
                  >
                    <item.icon className={`mr-3 h-5 w-5 ${
                      isActive ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 dark:text-gray-500'
                    }`} />
                    {item.name}
                  </Link>
                )
              })}
            </nav>
            <div className="border-t border-gray-200 dark:border-gray-700 p-4">
              <ThemeToggle />
            </div>
          </div>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-64 lg:flex-col">
        <div className="flex grow flex-col gap-y-5 overflow-y-auto bg-white dark:bg-gray-900 border-r border-gray-200 dark:border-gray-700">
          <div className="flex h-16 items-center px-6 border-b border-gray-200 dark:border-gray-700">
            <CommandLineIcon className="h-8 w-8 text-emerald-600" />
            <span className="ml-2 text-lg font-semibold text-gray-900 dark:text-white">SubdomainX</span>
          </div>
          <nav className="flex flex-1 flex-col px-4">
            <ul role="list" className="flex flex-1 flex-col gap-y-7">
              <li>
                <ul role="list" className="-mx-2 space-y-1">
                  {navigation.map((item) => {
                    const isActive = pathname === item.href
                    return (
                      <li key={item.name}>
                        <Link
                          href={item.href}
                          className={`group flex gap-x-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                            isActive
                              ? 'bg-emerald-50 dark:bg-emerald-900/20 text-emerald-700 dark:text-emerald-400 border border-emerald-200 dark:border-emerald-800'
                              : 'text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800 hover:text-gray-900 dark:hover:text-white'
                          }`}
                        >
                          <item.icon className={`h-5 w-5 shrink-0 ${
                            isActive ? 'text-emerald-600 dark:text-emerald-400' : 'text-gray-400 dark:text-gray-500'
                          }`} />
                          {item.name}
                        </Link>
                      </li>
                    )
                  })}
                </ul>
              </li>
              <li className="mt-auto">
                <div className="rounded-lg bg-gray-50 dark:bg-gray-800 p-4">
                  <div className="text-sm font-medium text-gray-900 dark:text-white">Need help?</div>
                  <div className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                    Check out our GitHub repository for issues and discussions.
                  </div>
                  <a
                    href="https://github.com/itszeeshan/subdomainx"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="mt-2 inline-flex items-center text-xs font-medium text-emerald-600 dark:text-emerald-400 hover:text-emerald-500 dark:hover:text-emerald-300"
                  >
                    View on GitHub
                  </a>
                </div>
              </li>
            </ul>
          </nav>
          <div className="border-t border-gray-200 dark:border-gray-700 p-4">
            <ThemeToggle />
          </div>
        </div>
      </div>

      {/* Mobile menu button */}
      <div className="sticky top-0 z-40 flex items-center gap-x-6 bg-white dark:bg-gray-950 px-4 py-4 shadow-sm sm:px-6 lg:hidden border-b border-gray-200 dark:border-gray-700">
        <button
          type="button"
          className="-m-2.5 p-2.5 text-gray-700 dark:text-gray-300 lg:hidden"
          onClick={() => setSidebarOpen(true)}
        >
          <span className="sr-only">Open sidebar</span>
          <Bars3Icon className="h-6 w-6" />
        </button>
        <div className="flex-1 text-sm font-semibold leading-6 text-gray-900 dark:text-white">
          SubdomainX Documentation
        </div>
        <ThemeToggle />
      </div>
    </>
  )
}
