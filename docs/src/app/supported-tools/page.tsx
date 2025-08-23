import Navigation from '@/components/Navigation'
import { 
  ShieldCheckIcon, 
  GlobeAltIcon,
  CommandLineIcon,
  InformationCircleIcon
} from '@heroicons/react/24/outline'

export default function SupportedTools() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Supported Tools</h1>
              <p className="mt-2 text-lg text-gray-600">
                SubdomainX integrates with popular subdomain enumeration and scanning tools.
              </p>
            </div>

            {/* Tool Management */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Tool Management</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">Check Tool Availability</h3>
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                      </div>
                      <code className="text-sm text-green-400 block">
                        subdomainx --check-tools
                      </code>
                    </div>
                  </div>
                  
                  <div>
                    <h3 className="text-lg font-semibold text-gray-900 mb-3">Get Installation Help</h3>
                    <div className="rounded-lg bg-gray-900 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
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

            {/* Subdomain Enumeration Tools */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Subdomain Enumeration Tools</h2>
              <div className="grid grid-cols-1 gap-6">
                {[
                  { name: 'subfinder', desc: 'Fast subdomain discovery tool with passive online sources' },
                  { name: 'amass', desc: 'Comprehensive subdomain enumeration and network mapping tool' },
                  { name: 'findomain', desc: 'Fast and cross-platform subdomain discovery tool' },
                  { name: 'assetfinder', desc: 'Find subdomains related to a domain' },
                  { name: 'sublist3r', desc: 'Subdomain enumeration using OSINT' },
                  { name: 'knockpy', desc: 'Subdomain enumeration tool' },
                  { name: 'dnsrecon', desc: 'DNS enumeration and reconnaissance tool' },
                  { name: 'fierce', desc: 'DNS reconnaissance tool' },
                  { name: 'massdns', desc: 'High-performance DNS stub resolver' },
                  { name: 'altdns', desc: 'Subdomain permutation and alteration' }
                ].map((tool) => (
                  <div key={tool.name} className="bg-white rounded-lg shadow-sm p-6">
                    <div className="flex items-center mb-2">
                      <ShieldCheckIcon className="h-5 w-5 text-indigo-600 mr-2" />
                      <h3 className="text-lg font-semibold text-gray-900">{tool.name}</h3>
                    </div>
                    <p className="text-gray-600">{tool.desc}</p>
                  </div>
                ))}
              </div>
            </div>

            {/* Scanning Tools */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Scanning Tools</h2>
              <div className="grid grid-cols-1 gap-6">
                {[
                  { name: 'httpx', desc: 'Fast and multi-purpose HTTP probe' },
                  { name: 'smap', desc: 'Port scanner and service discovery' }
                ].map((tool) => (
                  <div key={tool.name} className="bg-white rounded-lg shadow-sm p-6">
                    <div className="flex items-center mb-2">
                      <GlobeAltIcon className="h-5 w-5 text-green-600 mr-2" />
                      <h3 className="text-lg font-semibold text-gray-900">{tool.name}</h3>
                    </div>
                    <p className="text-gray-600">{tool.desc}</p>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
