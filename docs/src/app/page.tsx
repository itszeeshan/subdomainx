import Navigation from '@/components/Navigation'
import Link from 'next/link'
import { 
  CommandLineIcon, 
  RocketLaunchIcon,
  ShieldCheckIcon,
  BoltIcon,
  GlobeAltIcon,
  CogIcon,
  DocumentTextIcon,
  PlayIcon
} from '@heroicons/react/24/outline'

export default function Home() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            {/* Hero Section */}
            <div className="text-center mb-16">
              <div className="flex justify-center mb-6">
                <div className="bg-emerald-100 dark:bg-emerald-900/20 p-4 rounded-full">
                  <CommandLineIcon className="h-12 w-12 text-emerald-600 dark:text-emerald-400" />
                </div>
              </div>
              <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">
                SubdomainX
              </h1>
              <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto">
                All-in-one subdomain enumeration tool that combines the power of multiple tools 
                into a single, efficient command-line interface.
              </p>
              <div className="flex flex-col sm:flex-row gap-4 justify-center">
                <Link
                  href="/installation"
                  className="inline-flex items-center px-6 py-3 bg-emerald-600 hover:bg-emerald-700 text-white font-medium rounded-lg transition-colors"
                >
                  <RocketLaunchIcon className="h-5 w-5 mr-2" />
                  Get Started
                </Link>
                <Link
                  href="/cli-reference"
                  className="inline-flex items-center px-6 py-3 bg-gray-100 dark:bg-gray-800 hover:bg-gray-200 dark:hover:bg-gray-700 text-gray-900 dark:text-white font-medium rounded-lg transition-colors"
                >
                  <DocumentTextIcon className="h-5 w-5 mr-2" />
                  CLI Reference
                </Link>
              </div>
            </div>

            {/* Features */}
            <div className="mb-16">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-8 text-center">
                What&apos;s Included
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {[
                  {
                    icon: ShieldCheckIcon,
                    title: 'Multiple Tools',
                    description: 'Integrates subfinder, amass, findomain, and more for comprehensive enumeration.'
                  },
                  {
                    icon: BoltIcon,
                    title: 'High Performance',
                    description: 'Concurrent execution with configurable threads and rate limiting.'
                  },
                  {
                    icon: GlobeAltIcon,
                    title: 'HTTP Scanning',
                    description: 'Built-in httpx integration for web service discovery and analysis.'
                  },
                  {
                    icon: CogIcon,
                    title: 'Flexible Output',
                    description: 'Generate reports in JSON, TXT, or beautiful HTML formats.'
                  },
                  {
                    icon: PlayIcon,
                    title: 'Easy to Use',
                    description: 'Simple CLI interface with sensible defaults and optional configuration.'
                  },
                  {
                    icon: CommandLineIcon,
                    title: 'Tool Management',
                    description: 'Automatic tool detection and installation guidance for missing dependencies.'
                  }
                ].map((feature, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-shadow">
                    <div className="bg-emerald-100 dark:bg-emerald-900/20 p-3 rounded-lg w-fit mb-4">
                      <feature.icon className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
                    </div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-2">
                      {feature.title}
                    </h3>
                    <p className="text-gray-600 dark:text-gray-300">
                      {feature.description}
                    </p>
                  </div>
                ))}
              </div>
            </div>

            {/* Quick Start */}
            <div className="bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700 p-8 mb-16">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
                Quick Start
              </h2>
              <div className="space-y-4">
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
                  <h3 className="font-semibold text-gray-900 dark:text-white mb-2">1. Install</h3>
                  <code className="text-sm bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded">
                    go install github.com/itszeeshan/subdomainx@latest
                  </code>
                </div>
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
                  <h3 className="font-semibold text-gray-900 dark:text-white mb-2">2. Create domains file</h3>
                  <code className="text-sm bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded">
                    echo &quot;example.com&quot; &gt; domains.txt
                  </code>
                </div>
                <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
                  <h3 className="font-semibold text-gray-900 dark:text-white mb-2">3. Run scan</h3>
                  <code className="text-sm bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded">
                    subdomainx --wildcard domains.txt --format html
                  </code>
                </div>
              </div>
            </div>

            {/* Navigation Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <Link
                href="/installation"
                className="group bg-white dark:bg-gray-900 rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all hover:border-emerald-300 dark:hover:border-emerald-600"
              >
                <div className="flex items-center mb-4">
                  <div className="bg-emerald-100 dark:bg-emerald-900/20 p-2 rounded-lg mr-3">
                    <RocketLaunchIcon className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white group-hover:text-emerald-600 dark:group-hover:text-emerald-400">
                    Installation Guide
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Learn how to install SubdomainX and set up all required dependencies.
                </p>
              </Link>

              <Link
                href="/examples"
                className="group bg-white dark:bg-gray-900 rounded-lg p-6 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all hover:border-emerald-300 dark:hover:border-emerald-600"
              >
                <div className="flex items-center mb-4">
                  <div className="bg-emerald-100 dark:bg-emerald-900/20 p-2 rounded-lg mr-3">
                    <PlayIcon className="h-6 w-6 text-emerald-600 dark:text-emerald-400" />
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white group-hover:text-emerald-600 dark:group-hover:text-emerald-400">
                    Usage Examples
                  </h3>
                </div>
                <p className="text-gray-600 dark:text-gray-300">
                  Explore practical examples and use cases for different scanning scenarios.
                </p>
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
