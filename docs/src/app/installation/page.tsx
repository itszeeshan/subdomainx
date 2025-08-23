import Navigation from '@/components/Navigation'
import { 
  BookOpenIcon, 
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon
} from '@heroicons/react/24/outline'

export default function Installation() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            {/* Header */}
            <div className="mb-8 flex items-center">
              <div className="bg-emerald-100 dark:bg-emerald-900/20 p-3 rounded-lg mr-4">
                <BookOpenIcon className="h-8 w-8 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Installation Guide</h1>
                <p className="mt-2 text-lg text-gray-600 dark:text-gray-300">
                  Get SubdomainX up and running on your system.
                </p>
              </div>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Prerequisites</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="space-y-4">
                  {[
                    {
                      title: 'Go 1.21 or later',
                      description: 'Required for building and running the tool',
                      check: 'go version'
                    },
                    {
                      title: 'Git',
                      description: 'For cloning the repository (if building from source)',
                      check: 'git --version'
                    },
                    {
                      title: 'External Tools',
                      description: 'Various enumeration and scanning tools (optional, can be installed later)',
                      check: 'See Supported Tools page'
                    }
                  ].map((prereq, index) => (
                    <div key={index} className="flex items-start space-x-3">
                      <CheckCircleIcon className="h-6 w-6 text-emerald-500 mt-0.5 flex-shrink-0" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">{prereq.title}</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300">{prereq.description}</p>
                        <code className="text-xs bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded mt-1 inline-block">
                          {prereq.check}
                        </code>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Installation Methods */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Installation Methods</h2>
              
              {/* Quick Install */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
                  <CheckCircleIcon className="h-5 w-5 text-emerald-500 mr-2" />
                  Quick Install (Recommended)
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Install directly from GitHub using Go's module system:
                </p>
                <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                  <code className="text-sm text-emerald-400 block">
                    go install github.com/itszeeshan/subdomainx@latest
                  </code>
                </div>
                <div className="mt-4 p-3 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg border border-emerald-200 dark:border-emerald-800">
                  <div className="flex items-start">
                    <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-2 flex-shrink-0" />
                    <div className="text-sm text-emerald-700 dark:text-emerald-300">
                      <strong>Note:</strong> This will install the binary to your Go bin directory. Make sure it's in your PATH.
                    </div>
                  </div>
                </div>
              </div>

              {/* Build from Source */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Build from Source</h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Clone the repository and build manually:
                </p>
                <div className="space-y-3">
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      git clone https://github.com/itszeeshan/subdomainx.git
                    </code>
                  </div>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      cd subdomainx
                    </code>
                  </div>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      go build -o subdomainx cmd/subdomainx/main.go
                    </code>
                  </div>
                </div>
              </div>

              {/* Using Makefile */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Using Makefile</h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  If you have the source code, you can use the provided Makefile:
                </p>
                <div className="space-y-3">
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      make build
                    </code>
                  </div>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      make install
                    </code>
                  </div>
                </div>
              </div>
            </div>

            {/* Verification */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Verification</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Verify that SubdomainX is installed correctly:
                </p>
                <div className="space-y-3">
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --version
                    </code>
                  </div>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --help
                    </code>
                  </div>
                </div>
                <div className="mt-4 p-3 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg border border-emerald-200 dark:border-emerald-800">
                  <div className="flex items-start">
                    <CheckCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-2 flex-shrink-0" />
                    <div className="text-sm text-emerald-700 dark:text-emerald-300">
                      If you see the help output, SubdomainX is installed correctly!
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Next Steps */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Next Steps</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white dark:bg-gray-900 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Check Tool Availability</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                    Verify which enumeration tools are available on your system:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --check-tools
                    </code>
                  </div>
                </div>

                <div className="bg-white dark:bg-gray-900 rounded-lg p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Get Installation Help</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                    Get instructions for installing missing tools:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --install-tools
                    </code>
                  </div>
                </div>
              </div>
            </div>

            {/* Troubleshooting */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Troubleshooting</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="space-y-6">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                      <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500 mr-2" />
                      Command Not Found
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      If you get a "command not found" error after installation:
                    </p>
                    <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
                      <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">
                        <strong>Solution:</strong> Add your Go bin directory to your PATH:
                      </p>
                      <code className="text-sm bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded">
                        export PATH=$PATH:$(go env GOPATH)/bin
                      </code>
                    </div>
                  </div>

                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                      <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500 mr-2" />
                      Permission Denied
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      If you encounter permission issues:
                    </p>
                    <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-4">
                      <p className="text-sm text-gray-700 dark:text-gray-300 mb-2">
                        <strong>Solution:</strong> Make the binary executable:
                      </p>
                      <code className="text-sm bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-200 px-2 py-1 rounded">
                        chmod +x $(go env GOPATH)/bin/subdomainx
                      </code>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
