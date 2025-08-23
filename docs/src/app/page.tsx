import Navigation from '@/components/Navigation'
import { 
  CommandLineIcon, 
  RocketLaunchIcon, 
  ShieldCheckIcon, 
  CogIcon,
  ArrowRightIcon,
  CheckCircleIcon
} from '@heroicons/react/24/outline'

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      {/* Main content */}
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          {/* Hero section */}
          <div className="mx-auto max-w-4xl">
            <div className="text-center">
              <div className="flex justify-center">
                <div className="flex h-16 w-16 items-center justify-center rounded-full bg-indigo-100">
                  <CommandLineIcon className="h-8 w-8 text-indigo-600" />
                </div>
              </div>
              <h1 className="mt-6 text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl">
                SubdomainX Documentation
              </h1>
              <p className="mt-4 text-lg text-gray-600">
                All-in-one subdomain enumeration tool with comprehensive scanning capabilities
              </p>
              <div className="mt-8 flex justify-center gap-4">
                <a
                  href="/installation"
                  className="inline-flex items-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
                >
                  Get Started
                  <ArrowRightIcon className="ml-2 h-4 w-4" />
                </a>
                <a
                  href="https://github.com/itszeeshan/subdomainx"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                >
                  View on GitHub
                </a>
              </div>
            </div>

            {/* Features */}
            <div className="mt-16 grid grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
              <div className="rounded-lg bg-white p-6 shadow-sm">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-indigo-100">
                  <ShieldCheckIcon className="h-6 w-6 text-indigo-600" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-gray-900">Multiple Tools</h3>
                <p className="mt-2 text-sm text-gray-600">
                  Integrates with popular subdomain enumeration tools like subfinder, amass, findomain, and more.
                </p>
              </div>

              <div className="rounded-lg bg-white p-6 shadow-sm">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-green-100">
                  <CogIcon className="h-6 w-6 text-green-600" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-gray-900">Easy Installation</h3>
                <p className="mt-2 text-sm text-gray-600">
                  Simple installation with go install. No complex setup required.
                </p>
              </div>

              <div className="rounded-lg bg-white p-6 shadow-sm">
                <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-purple-100">
                  <RocketLaunchIcon className="h-6 w-6 text-purple-600" />
                </div>
                <h3 className="mt-4 text-lg font-semibold text-gray-900">Fast & Efficient</h3>
                <p className="mt-2 text-sm text-gray-600">
                  Multi-threaded execution with rate limiting and DNS caching for optimal performance.
                </p>
              </div>
            </div>

            {/* Quick start */}
            <div className="mt-16">
              <h2 className="text-2xl font-bold text-gray-900">Quick Start</h2>
              <div className="mt-6 rounded-lg bg-gray-900 p-6">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-300">Installation</span>
                  <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                </div>
                <div className="mt-2">
                  <code className="text-sm text-green-400">
                    go install github.com/itszeeshan/subdomainx/cmd/subdomainx@latest
                  </code>
                </div>
              </div>

              <div className="mt-4 rounded-lg bg-gray-900 p-6">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-gray-300">Basic Usage</span>
                  <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                </div>
                <div className="mt-2">
                  <code className="text-sm text-green-400">
                    subdomainx --wildcard domains.txt --format html
                  </code>
                </div>
              </div>
            </div>

            {/* What's included */}
            <div className="mt-16">
              <h2 className="text-2xl font-bold text-gray-900">What&apos;s Included</h2>
              <div className="mt-6 grid grid-cols-1 gap-4 sm:grid-cols-2">
                <div className="rounded-lg bg-white p-6 shadow-sm">
                  <h3 className="text-lg font-semibold text-gray-900">Subdomain Enumeration</h3>
                  <ul className="mt-4 space-y-2">
                    {[
                      'subfinder - Fast subdomain discovery',
                      'amass - Comprehensive enumeration',
                      'findomain - Cross-platform discovery',
                      'assetfinder - Domain-related subdomains',
                      'sublist3r - OSINT-based enumeration',
                      'knockpy - Subdomain enumeration tool',
                      'dnsrecon - DNS reconnaissance',
                      'fierce - DNS reconnaissance tool',
                      'massdns - High-performance DNS resolver',
                      'altdns - Subdomain permutation'
                    ].map((tool) => (
                      <li key={tool} className="flex items-start">
                        <CheckCircleIcon className="mr-2 h-4 w-4 text-green-500 mt-0.5" />
                        <span className="text-sm text-gray-600">{tool}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="rounded-lg bg-white p-6 shadow-sm">
                  <h3 className="text-lg font-semibold text-gray-900">Scanning Tools</h3>
                  <ul className="mt-4 space-y-2">
                    {[
                      'httpx - Fast HTTP probe',
                      'smap - Port scanner and service discovery'
                    ].map((tool) => (
                      <li key={tool} className="flex items-start">
                        <CheckCircleIcon className="mr-2 h-4 w-4 text-green-500 mt-0.5" />
                        <span className="text-sm text-gray-600">{tool}</span>
                      </li>
                    ))}
                  </ul>

                  <h3 className="mt-6 text-lg font-semibold text-gray-900">Output Formats</h3>
                  <ul className="mt-4 space-y-2">
                    {[
                      'JSON - Structured data output',
                      'TXT - Plain text format',
                      'HTML - Beautiful web reports'
                    ].map((format) => (
                      <li key={format} className="flex items-start">
                        <CheckCircleIcon className="mr-2 h-4 w-4 text-green-500 mt-0.5" />
                        <span className="text-sm text-gray-600">{format}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
