import Navigation from '@/components/Navigation'
import { 
  CogIcon, 
  DocumentTextIcon,
  InformationCircleIcon
} from '@heroicons/react/24/outline'

export default function Configuration() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Configuration</h1>
              <p className="mt-2 text-lg text-gray-600">
                Learn how to configure SubdomainX using YAML files and command-line options.
              </p>
            </div>

            {/* Configuration Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Configuration Overview</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">YAML Configuration</h3>
                    <p className="text-gray-600 mb-3">
                      Use a YAML file for default settings and complex configurations.
                    </p>
                    <ul className="space-y-2 text-sm text-gray-600">
                      <li>• Default location: <code className="bg-gray-100 px-1 rounded">configs/default.yaml</code></li>
                      <li>• Optional - CLI flags take precedence</li>
                      <li>• Good for team settings and complex setups</li>
                    </ul>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">CLI Options</h3>
                    <p className="text-gray-600 mb-3">
                      Use command-line flags for quick configuration and automation.
                    </p>
                    <ul className="space-y-2 text-sm text-gray-600">
                      <li>• Override YAML settings</li>
                      <li>• Perfect for scripts and automation</li>
                      <li>• Immediate effect without file editing</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>

            {/* YAML Configuration */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">YAML Configuration File</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="mb-4">
                  <h3 className="text-lg font-semibold text-gray-900 mb-2">Default Configuration</h3>
                  <p className="text-sm text-gray-600 mb-4">
                    Create a <code className="bg-gray-100 px-1 rounded">configs/default.yaml</code> file with your preferred settings:
                  </p>
                </div>
                
                <div className="rounded-lg bg-gray-900 p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300">configs/default.yaml</span>
                    <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                  </div>
                  <pre className="text-sm text-green-400 overflow-x-auto">
{`# SubdomainX Configuration File
wildcard_file: domains.txt          # File containing target domains
unique_name: scan                   # Unique name for output files
output_dir: output                  # Output directory
output_format: json                 # Output format: json, txt, html
threads: 10                         # Number of concurrent threads
retries: 3                          # Number of retry attempts
timeout: 30                         # Timeout in seconds
rate_limit: 100                     # Rate limit per second
wordlist: ""                        # Custom wordlist (optional)

# Filters for results
filters:
  status_code: "200,301,302"        # Filter by HTTP status codes
  ports: "80,443,8080"              # Filter by ports

# Tool selection
tools:
  subfinder: true                   # Enable/disable specific tools
  findomain: true
  assetfinder: true
  amass: true
  sublist3r: true
  knockpy: true
  dnsrecon: true
  fierce: true
  massdns: true
  altdns: true
  httpx: true                       # HTTP scanning
  smap: true                        # Port scanning`}
                  </pre>
                </div>
              </div>
            </div>

            {/* Configuration Options */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Configuration Options</h2>
              
              {/* Basic Options */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Basic Options</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">wildcard_file</code>
                      <span className="text-xs text-gray-500">Required</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Path to file containing target domains (one per line)
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">unique_name</code>
                      <span className="text-xs text-gray-500">Default: scan</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Unique name for output files (e.g., scan_results.json)
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">output_dir</code>
                      <span className="text-xs text-gray-500">Default: output</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Directory where output files will be saved
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">output_format</code>
                      <span className="text-xs text-gray-500">Default: json</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Output format: json, txt, html
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
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">threads</code>
                      <span className="text-xs text-gray-500">Default: 10</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Number of concurrent threads for processing
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">retries</code>
                      <span className="text-xs text-gray-500">Default: 3</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Number of retry attempts for failed operations
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">timeout</code>
                      <span className="text-xs text-gray-500">Default: 30</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Timeout in seconds for individual operations
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">rate_limit</code>
                      <span className="text-xs text-gray-500">Default: 100</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Rate limit per second to avoid overwhelming targets
                    </p>
                  </div>
                </div>
              </div>

              {/* Filters */}
              <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Filters</h3>
                <div className="space-y-4">
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">filters.status_code</code>
                      <span className="text-xs text-gray-500">Optional</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Comma-separated list of HTTP status codes to include
                    </p>
                  </div>
                  
                  <div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">filters.ports</code>
                      <span className="text-xs text-gray-500">Optional</span>
                    </div>
                    <p className="text-sm text-gray-600 mt-1">
                      Comma-separated list of ports to include in results
                    </p>
                  </div>
                </div>
              </div>

              {/* Tool Selection */}
              <div className="bg-white rounded-lg shadow-sm p-6">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Tool Selection</h3>
                <p className="text-sm text-gray-600 mb-4">
                  Enable or disable specific tools by setting them to <code className="bg-gray-100 px-1 rounded">true</code> or <code className="bg-gray-100 px-1 rounded">false</code>.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.subfinder</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.amass</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.findomain</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.assetfinder</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.sublist3r</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.knockpy</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.dnsrecon</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.fierce</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.massdns</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-gray-100 px-2 py-1 rounded">tools.altdns</code>
                      <span className="text-xs text-gray-500">Default: true</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Using Custom Configuration */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Using Custom Configuration</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="space-y-4">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">Custom Config File</h3>
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                      </div>
                      <code className="text-sm text-green-400 block">
                        subdomainx --wildcard domains.txt --config custom.yaml
                      </code>
                    </div>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">CLI Override</h3>
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                      </div>
                      <code className="text-sm text-green-400 block">
                        subdomainx --wildcard domains.txt --threads 20 --format html --name my_scan
                      </code>
                    </div>
                    <p className="text-sm text-gray-600 mt-2">
                      CLI flags always override YAML configuration values
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Configuration Tips */}
            <div className="bg-gradient-to-r from-indigo-50 to-purple-50 rounded-lg p-6 border border-indigo-100">
              <h2 className="text-xl font-bold text-gray-900 mb-3">Configuration Tips</h2>
              <div className="space-y-3">
                <div className="flex items-start">
                  <InformationCircleIcon className="h-5 w-5 text-indigo-500 mt-0.5 mr-2" />
                  <div className="text-sm text-gray-700">
                    <strong>Start Simple:</strong> Begin with basic configuration and add complexity as needed.
                  </div>
                </div>
                <div className="flex items-start">
                  <InformationCircleIcon className="h-5 w-5 text-indigo-500 mt-0.5 mr-2" />
                  <div className="text-sm text-gray-700">
                    <strong>Team Settings:</strong> Use YAML files for consistent team configurations.
                  </div>
                </div>
                <div className="flex items-start">
                  <InformationCircleIcon className="h-5 w-5 text-indigo-500 mt-0.5 mr-2" />
                  <div className="text-sm text-gray-700">
                    <strong>Automation:</strong> Use CLI flags in scripts for automated scanning.
                  </div>
                </div>
                <div className="flex items-start">
                  <InformationCircleIcon className="h-5 w-5 text-indigo-500 mt-0.5 mr-2" />
                  <div className="text-sm text-gray-700">
                    <strong>Performance:</strong> Adjust threads and rate limits based on your network and target.
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
