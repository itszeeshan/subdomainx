import Navigation from '@/components/Navigation'
import { 
  CommandLineIcon, 
  DocumentTextIcon,
  CodeBracketIcon
} from '@heroicons/react/24/outline'

const examples = [
  {
    title: 'Basic Subdomain Enumeration',
    description: 'Run a basic scan with all available tools',
    command: 'subdomainx --wildcard domains.txt',
    explanation: 'This will use all available enumeration tools to discover subdomains for the domains listed in domains.txt'
  },
  {
    title: 'Specific Tools Only',
    description: 'Use only specific enumeration tools',
    command: 'subdomainx --wildcard domains.txt --amass --subfinder --findomain',
    explanation: 'Limit the scan to only amass, subfinder, and findomain for faster execution'
  },
  {
    title: 'Generate HTML Report',
    description: 'Create a beautiful HTML report',
    command: 'subdomainx --wildcard domains.txt --format html --name my_scan',
    explanation: 'Generate an HTML report with the name "my_scan" for easy viewing and sharing'
  },
  {
    title: 'High Performance Scan',
    description: 'Optimize for speed with more threads',
    command: 'subdomainx --wildcard domains.txt --threads 20 --timeout 60',
    explanation: 'Use 20 concurrent threads and 60-second timeout for faster scanning'
  },
  {
    title: 'HTTP Scanning Only',
    description: 'Perform HTTP scanning on discovered subdomains',
    command: 'subdomainx --wildcard domains.txt --httpx',
    explanation: 'Use httpx to discover web services and extract information from discovered subdomains'
  },
  {
    title: 'Port Scanning',
    description: 'Perform port scanning on discovered hosts',
    command: 'subdomainx --wildcard domains.txt --smap',
    explanation: 'Use smap to identify open ports and services on discovered hosts'
  },
  {
    title: 'Complete Reconnaissance',
    description: 'Full enumeration with HTTP and port scanning',
    command: 'subdomainx --wildcard domains.txt --amass --subfinder --httpx --smap --format html',
    explanation: 'Comprehensive reconnaissance with enumeration, HTTP scanning, port scanning, and HTML report'
  },
  {
    title: 'Custom Output Directory',
    description: 'Save results to a custom directory',
    command: 'subdomainx --wildcard domains.txt --output /path/to/results --name penetration_test',
    explanation: 'Save all results to /path/to/results with the name "penetration_test"'
  }
]

export default function Examples() {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            <div className="mb-8">
              <h1 className="text-3xl font-bold text-gray-900">Examples</h1>
              <p className="mt-2 text-lg text-gray-600">
                Practical examples of how to use SubdomainX for different scenarios.
              </p>
            </div>

            {/* Examples Grid */}
            <div className="space-y-8">
              {examples.map((example, index) => (
                <div key={index} className="bg-white rounded-lg shadow-sm p-6">
                  <div className="flex items-start">
                    <div className="flex-shrink-0">
                      <div className="flex h-8 w-8 items-center justify-center rounded-full bg-indigo-100">
                        <CodeBracketIcon className="h-5 w-5 text-indigo-600" />
                      </div>
                    </div>
                    <div className="ml-4 flex-1">
                      <h3 className="text-lg font-semibold text-gray-900">{example.title}</h3>
                      <p className="text-sm text-gray-600 mt-1">{example.description}</p>
                      
                      <div className="mt-4 rounded-lg bg-gray-900 p-4">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium text-gray-300">Command</span>
                          <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                        </div>
                        <code className="text-sm text-green-400 block">{example.command}</code>
                      </div>
                      
                      <div className="mt-3 p-3 bg-blue-50 rounded-lg">
                        <p className="text-sm text-blue-700">
                          <strong>Explanation:</strong> {example.explanation}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Sample Files */}
            <div className="mt-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Sample Files</h2>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">Sample domains.txt</h3>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">domains.txt</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <pre className="text-sm text-green-400">
{`example.com
test.org
demo.net
sample.co.uk`}
                    </pre>
                  </div>
                </div>
                
                <div className="bg-white rounded-lg shadow-sm p-6">
                  <h3 className="text-lg font-semibold text-gray-900 mb-3">Sample Configuration</h3>
                  <div className="rounded-lg bg-gray-900 p-4">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-gray-300">config.yaml</span>
                      <button className="text-xs text-gray-400 hover:text-gray-300">Copy</button>
                    </div>
                    <pre className="text-sm text-green-400">
{`wildcard_file: domains.txt
unique_name: my_scan
output_format: html
threads: 15
retries: 3
timeout: 45
rate_limit: 150
tools:
  subfinder: true
  amass: true
  findomain: true
  httpx: true`}
                    </pre>
                  </div>
                </div>
              </div>
            </div>

            {/* Best Practices */}
            <div className="mt-12">
              <h2 className="text-2xl font-bold text-gray-900 mb-4">Best Practices</h2>
              <div className="bg-white rounded-lg shadow-sm p-6">
                <div className="space-y-4">
                  <div className="flex items-start">
                    <div className="w-2 h-2 bg-indigo-600 rounded-full mt-2 mr-3"></div>
                    <div>
                      <h4 className="font-medium text-gray-900">Start with Basic Scans</h4>
                      <p className="text-sm text-gray-600 mt-1">
                        Begin with basic enumeration before adding HTTP and port scanning for faster initial results.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <div className="w-2 h-2 bg-indigo-600 rounded-full mt-2 mr-3"></div>
                    <div>
                      <h4 className="font-medium text-gray-900">Use Appropriate Thread Counts</h4>
                      <p className="text-sm text-gray-600 mt-1">
                        Higher thread counts speed up scanning but may trigger rate limits. Start with 10-15 threads.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <div className="w-2 h-2 bg-indigo-600 rounded-full mt-2 mr-3"></div>
                    <div>
                      <h4 className="font-medium text-gray-900">Organize Output Files</h4>
                      <p className="text-sm text-gray-600 mt-1">
                        Use descriptive names and organize output directories for better result management.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <div className="w-2 h-2 bg-indigo-600 rounded-full mt-2 mr-3"></div>
                    <div>
                      <h4 className="font-medium text-gray-900">Check Tool Availability</h4>
                      <p className="text-sm text-gray-600 mt-1">
                        Always run <code className="bg-gray-100 px-1 rounded text-xs">--check-tools</code> to see which tools are available.
                      </p>
                    </div>
                  </div>
                  
                  <div className="flex items-start">
                    <div className="w-2 h-2 bg-indigo-600 rounded-full mt-2 mr-3"></div>
                    <div>
                      <h4 className="font-medium text-gray-900">Use HTML Reports for Sharing</h4>
                      <p className="text-sm text-gray-600 mt-1">
                        HTML reports are perfect for sharing results with stakeholders and team members.
                      </p>
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
