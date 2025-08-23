import Navigation from '@/components/Navigation'
import { 
  CommandLineIcon, 
  InformationCircleIcon,
  ExclamationTriangleIcon
} from '@heroicons/react/24/outline'

export default function CLIReference() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">CLI Reference</h1>
              <p className="mt-2 text-lg text-gray-600">
                Complete command-line interface reference for SubdomainX.
              </p>
            </div>

            {/* Basic Usage */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Basic Usage</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="rounded-lg bg-gray-900 p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300">Basic Command</span>
                    <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                  </div>
                  <code className="text-sm text-green-400 block">
                    subdomainx --wildcard &lt;domains_file&gt; [OPTIONS]
                  </code>
                </div>
                
                <div className="mt-4 text-sm text-gray-600">
                  <p>
                    <strong>Required:</strong> You must provide a file containing target domains using the <code className="bg-gray-100 px-1 rounded">--wildcard</code> flag.
                  </p>
                </div>
              </div>
            </div>

            {/* Command Options */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Command Options</h2>
              
              {/* Required Options */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4 flex items-center">
                  <ExclamationTriangleIcon className="h-5 w-5 text-red-500 mr-2" />
                  Required Options
                </h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--wildcard FILE</code>
                      <span className="text-xs text-gray-500">Required</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Path to file containing target domains (one per line)
                    </p>
                  </div>
                </div>
              </div>

              {/* Tool Selection */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Tool Selection</h3>
                <p className="text-sm text-gray-600 mb-4">
                  Use specific tools, otherwise all available tools will be used.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--subfinder</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use subfinder tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--amass</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use amass tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--findomain</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use findomain tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--assetfinder</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use assetfinder tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--sublist3r</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use sublist3r tool</p>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--knockpy</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use knockpy tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--dnsrecon</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use dnsrecon tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--fierce</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use fierce tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--massdns</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use massdns tool</p>
                    
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--altdns</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-xs text-gray-600">Use altdns tool</p>
                  </div>
                </div>
                
                <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                  <div className="flex items-start">
                    <InformationCircleIcon className="h-5 w-5 text-blue-500 mt-0.5 mr-2" />
                    <div className="text-sm text-blue-700">
                      <strong>Note:</strong> If no specific tools are specified, SubdomainX will use all available tools on your system.
                    </div>
                  </div>
                </div>
              </div>

              {/* Scanning Tools */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Scanning Tools</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--httpx</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Use httpx for HTTP scanning (discovers web services, extracts titles, status codes, and technologies)
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--smap</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Use smap for port scanning (identifies open ports and services on discovered hosts)
                    </p>
                  </div>
                </div>
              </div>

              {/* Output Options */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Output Options</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--name NAME</code>
                      <span className="text-xs text-gray-500">Default: scan</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Unique name for output files
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--format FORMAT</code>
                      <span className="text-xs text-gray-500">Default: json</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Output format: json, txt, html
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--output DIR</code>
                      <span className="text-xs text-gray-500">Default: output</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Output directory for generated files
                    </p>
                  </div>
                </div>
              </div>

              {/* Performance Options */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Performance Options</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--threads N</code>
                      <span className="text-xs text-gray-500">Default: 10</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Number of concurrent threads
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--retries N</code>
                      <span className="text-xs text-gray-500">Default: 3</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Number of retry attempts
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--timeout N</code>
                      <span className="text-xs text-gray-500">Default: 30</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Timeout in seconds
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--rate-limit N</code>
                      <span className="text-xs text-gray-500">Default: 100</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Rate limit per second
                    </p>
                  </div>
                </div>
              </div>

              {/* Utility Options */}
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Utility Options</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--help</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Show help message
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--version</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Show version information
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--check-tools</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Check availability of enumeration tools
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--install-tools</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Show installation instructions for missing tools
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">--config FILE</code>
                      <span className="text-xs text-gray-500">Optional</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Use custom configuration file
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">-v, --verbose</code>
                      <span className="text-xs text-gray-500">Flag</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Enable verbose output
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Examples */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Examples</h2>
              
              <div className="space-y-6">
                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">Basic Scan</h3>
                  <p className="text-sm text-gray-600 mb-3">
                    Run a basic scan with all available tools:
                  </p>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Command</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --wildcard domains.txt
                    </code>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">Specific Tools</h3>
                  <p className="text-sm text-gray-600 mb-3">
                    Use only specific enumeration and scanning tools:
                  </p>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Command</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --wildcard domains.txt --amass --subfinder --httpx
                    </code>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">HTML Report</h3>
                  <p className="text-sm text-gray-600 mb-3">
                    Generate a beautiful HTML report with custom name:
                  </p>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Command</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --wildcard domains.txt --format html --name my_scan
                    </code>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">High Performance</h3>
                  <p className="text-sm text-gray-600 mb-3">
                    Optimize for speed with more threads and longer timeout:
                  </p>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">Command</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <code className="text-sm text-green-400 block">
                      subdomainx --wildcard domains.txt --threads 20 --timeout 60
                    </code>
                  </div>
                </div>

                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">Tool Management</h3>
                  <p className="text-sm text-gray-600 mb-3">
                    Check which tools are available and get installation help:
                  </p>
                  <div className="space-y-3">
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Check Tools</span>
                        <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                      </div>
                      <code className="text-sm text-green-400 block">
                        subdomainx --check-tools
                      </code>
                    </div>
                    
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Install Help</span>
                        <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                      </div>
                      <code className="text-sm text-green-400 block">
                        subdomainx --install-tools
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
