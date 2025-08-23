import Navigation from '@/components/Navigation'
import { 
  CommandLineIcon, 
  CheckCircleIcon,
  ExclamationTriangleIcon,
  ArrowRightIcon
} from '@heroicons/react/24/outline'

export default function Installation() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Installation</h1>
              <p className="mt-2 text-lg text-gray-600">
                Get SubdomainX up and running on your system with these simple installation methods.
              </p>
            </div>

            {/* Quick Install */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Quick Install</h2>
              <div className="rounded-lg bg-gradient-to-r from-indigo-50 to-purple-50 p-6 border border-indigo-100">
                <div className="flex items-start">
                  <CheckCircleIcon className="h-6 w-6 text-green-500 mt-1 mr-3" />
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900">Recommended Method</h3>
                    <p className="text-gray-600 mt-1">
                      Install directly using Go&apos;s package manager. This is the fastest and most reliable method.
                    </p>
                  </div>
                </div>
                
                <div className="mt-6 rounded-lg bg-gray-900 p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300">Install Command</span>
                    <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                  </div>
                  <code className="text-sm text-green-400 block">
                    go install github.com/itszeeshan/subdomainx/cmd/subdomainx@latest
                  </code>
                </div>

                <div className="mt-4 rounded-lg bg-gray-900 p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300">Verify Installation</span>
                    <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                  </div>
                  <code className="text-sm text-green-400 block">
                    subdomainx --help
                  </code>
                </div>
              </div>
            </div>

            {/* Prerequisites */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Prerequisites</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Required Software</h3>
                <div className="space-y-4">
                  <div className="flex items-start">
                    <CheckCircleIcon className="h-5 w-5 text-green-500 mt-0.5 mr-3" />
                    <div>
                      <span className="font-medium text-gray-900">Go 1.21 or later</span>
                      <p className="text-sm text-gray-600 mt-1">
                        Required for building and running SubdomainX. Download from{' '}
                        <a href="https://golang.org/dl/" className="text-indigo-600 hover:text-indigo-500">golang.org</a>
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <CheckCircleIcon className="h-5 w-5 text-green-500 mt-0.5 mr-3" />
                    <div>
                      <span className="font-medium text-gray-900">External Tools (Optional)</span>
                      <p className="text-sm text-gray-600 mt-1">
                        SubdomainX integrates with external tools for enhanced functionality. 
                        These are optional but recommended for full feature access.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Installation Methods */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Installation Methods</h2>
              
              {/* Method 1: Go Install */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">1. Go Install (Recommended)</h3>
                <p className="text-gray-600 mb-4">
                  Install directly from the Go module registry. This method automatically handles dependencies and updates.
                </p>
                
                <div className="space-y-3">
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Install</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      go install github.com/itszeeshan/subdomainx/cmd/subdomainx@latest
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Make Global (Optional)</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      sudo mv $(go env GOPATH)/bin/subdomainx /usr/local/bin/
                    </code>
                  </div>
                </div>
              </div>

              {/* Method 2: Build from Source */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">2. Build from Source</h3>
                <p className="text-gray-600 mb-4">
                  Clone the repository and build locally. Useful for development or custom modifications.
                </p>
                
                <div className="space-y-3">
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Clone Repository</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      git clone https://github.com/itszeeshan/subdomainx.git
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Build Binary</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      cd subdomainx && go build -o subdomainx ./cmd/subdomainx
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Install Globally</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      sudo mv subdomainx /usr/local/bin/
                    </code>
                  </div>
                </div>
              </div>

              {/* Method 3: Using Makefile */}
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-3">3. Using Makefile</h3>
                <p className="text-gray-600 mb-4">
                  Use the provided Makefile for common development tasks and installation.
                </p>
                
                <div className="space-y-3">
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Build and Install</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      make install
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Install from GitHub</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      make install-remote
                    </code>
                  </div>
                </div>
              </div>
            </div>

            {/* Verification */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Verify Installation</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <p className="text-gray-600 mb-4">
                  After installation, verify that SubdomainX is working correctly:
                </p>
                
                <div className="space-y-3">
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Check Version</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --version
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Show Help</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --help
                    </code>
                  </div>
                  
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Check Tools</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --check-tools
                    </code>
                  </div>
                </div>
              </div>
            </div>

            {/* Next Steps */}
            <div className="bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg p-6 border border-indigo-100">
              <h2 className="text-xl font-bold text-gray-900 mb-3">Next Steps</h2>
              <p className="text-gray-600 mb-4">
                Now that you have SubdomainX installed, you can start using it for subdomain enumeration.
              </p>
              <div className="flex flex-col sm:flex-row gap-3">
                <a
                  href="/cli-reference"
                  className="inline-flex items-center rounded-md bg-indigo-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
                >
                  CLI Reference
                  <ArrowRightIcon className="ml-2 h-4 w-4" />
                </a>
                <a
                  href="/examples"
                  className="inline-flex items-center rounded-md bg-white px-4 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
                >
                  View Examples
                </a>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
