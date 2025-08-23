import Navigation from '@/components/Navigation'
import { 
  CogIcon, 
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon
} from '@heroicons/react/24/outline'

export default function SupportedTools() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            {/* Header */}
            <div className="mb-8 flex items-center">
              <div className="bg-emerald-100 dark:bg-emerald-900/20 p-3 rounded-lg mr-4">
                <CogIcon className="h-8 w-8 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Supported Tools</h1>
                <p className="mt-2 text-lg text-gray-600 dark:text-gray-300">
                  Comprehensive list of enumeration and scanning tools integrated with SubdomainX.
                </p>
              </div>
            </div>

            {/* Enumeration Tools */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Enumeration Tools</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {[
                  {
                    name: 'subfinder',
                    description: 'Fast subdomain discovery tool with passive online sources',
                    website: 'https://github.com/projectdiscovery/subfinder',
                    install: 'go install -v github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest'
                  },
                  {
                    name: 'amass',
                    description: 'Comprehensive network reconnaissance and attack surface mapping',
                    website: 'https://github.com/owasp-amass/amass',
                    install: 'go install -v github.com/owasp-amass/amass/v4/...@master'
                  },
                  {
                    name: 'findomain',
                    description: 'Cross-platform subdomain enumeration tool',
                    website: 'https://github.com/findomain/findomain',
                    install: 'curl -LO https://github.com/findomain/findomain/releases/latest/download/findomain-linux'
                  },
                  {
                    name: 'assetfinder',
                    description: 'Find domains and subdomains potentially related to a given domain',
                    website: 'https://github.com/tomnomnom/assetfinder',
                    install: 'go install github.com/tomnomnom/assetfinder@latest'
                  },
                  {
                    name: 'sublist3r',
                    description: 'Fast subdomain enumeration tool for penetration testers',
                    website: 'https://github.com/aboul3la/Sublist3r',
                    install: 'pip install sublist3r'
                  },
                  {
                    name: 'knockpy',
                    description: 'Python3 tool designed to enumerate subdomains on a target domain',
                    website: 'https://github.com/guelfoweb/knock',
                    install: 'pip install knockpy'
                  },
                  {
                    name: 'dnsrecon',
                    description: 'DNS enumeration and reconnaissance tool',
                    website: 'https://github.com/darkoperator/dnsrecon',
                    install: 'pip install dnsrecon'
                  },
                  {
                    name: 'fierce',
                    description: 'DNS reconnaissance tool for locating non-contiguous IP space',
                    website: 'https://github.com/davidpepper/fierce',
                    install: 'pip install fierce'
                  },
                  {
                    name: 'massdns',
                    description: 'High-performance DNS stub resolver',
                    website: 'https://github.com/blechschmidt/massdns',
                    install: 'git clone https://github.com/blechschmidt/massdns.git && cd massdns && make'
                  },
                  {
                    name: 'altdns',
                    description: 'Generates permutations, alterations and mutations of subdomains',
                    website: 'https://github.com/infosec-au/altdns',
                    install: 'pip install py-altdns'
                  }
                ].map((tool, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <div className="flex items-start justify-between mb-3">
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{tool.name}</h3>
                      <CheckCircleIcon className="h-5 w-5 text-emerald-500 flex-shrink-0" />
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                      {tool.description}
                    </p>
                    <div className="space-y-2">
                      <a
                        href={tool.website}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-emerald-600 dark:text-emerald-400 hover:text-emerald-500 dark:hover:text-emerald-300 font-medium"
                      >
                        View on GitHub →
                      </a>
                      <div className="bg-gray-50 dark:bg-gray-800 rounded p-3">
                        <code className="text-xs text-gray-800 dark:text-gray-200 block overflow-x-auto">
                          {tool.install}
                        </code>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Scanning Tools */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Scanning Tools</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {[
                  {
                    name: 'httpx',
                    description: 'Fast and multi-purpose HTTP probe for web services',
                    website: 'https://github.com/projectdiscovery/httpx',
                    install: 'go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest'
                  },
                  {
                    name: 'smap',
                    description: 'Port scanner and service discovery tool',
                    website: 'https://github.com/s0md3v/Smap',
                    install: 'pip install smap'
                  }
                ].map((tool, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <div className="flex items-start justify-between mb-3">
                      <h3 className="text-lg font-semibold text-gray-900 dark:text-white">{tool.name}</h3>
                      <CheckCircleIcon className="h-5 w-5 text-emerald-500 flex-shrink-0" />
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                      {tool.description}
                    </p>
                    <div className="space-y-2">
                      <a
                        href={tool.website}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-xs text-emerald-600 dark:text-emerald-400 hover:text-emerald-500 dark:hover:text-emerald-300 font-medium"
                      >
                        View on GitHub →
                      </a>
                      <div className="bg-gray-50 dark:bg-gray-800 rounded p-3">
                        <code className="text-xs text-gray-800 dark:text-gray-200 block overflow-x-auto">
                          {tool.install}
                        </code>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Tool Management */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Tool Management</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="space-y-6">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                      <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mr-2" />
                      Check Tool Availability
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      Verify which tools are available on your system:
                    </p>
                    <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                      <code className="text-sm text-emerald-400 block">
                        subdomainx --check-tools
                      </code>
                    </div>
                  </div>

                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3 flex items-center">
                      <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500 mr-2" />
                      Get Installation Help
                    </h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      Get detailed installation instructions for missing tools:
                    </p>
                    <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                      <code className="text-sm text-emerald-400 block">
                        subdomainx --install-tools
                      </code>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Tips */}
            <div className="bg-emerald-50 dark:bg-emerald-900/20 rounded-lg p-6 border border-emerald-200 dark:border-emerald-800">
              <div className="flex items-start">
                <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-3 flex-shrink-0" />
                <div>
                  <h3 className="text-lg font-semibold text-emerald-900 dark:text-emerald-100 mb-2">
                    Pro Tips
                  </h3>
                  <ul className="text-sm text-emerald-800 dark:text-emerald-200 space-y-1">
                    <li>• Install tools as needed - SubdomainX will work with any combination</li>
                    <li>• More tools = better coverage and results</li>
                    <li>• Some tools may require API keys for optimal performance</li>
                    <li>• Check tool documentation for specific configuration requirements</li>
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
