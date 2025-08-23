import Navigation from '@/components/Navigation'
import { 
  InformationCircleIcon, 
  CheckCircleIcon,
  ExclamationTriangleIcon,
  CogIcon
} from '@heroicons/react/24/outline'

export default function Configuration() {
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
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Configuration</h1>
                <p className="mt-2 text-lg text-gray-600 dark:text-gray-300">
                  Learn how to configure SubdomainX using YAML files and CLI options.
                </p>
              </div>
            </div>

            {/* Overview */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Configuration Overview</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="space-y-4">
                  <div className="flex items-start space-x-3">
                    <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 flex-shrink-0" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-white">CLI-First Approach</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-300">
                        SubdomainX prioritizes command-line arguments over configuration files. All options can be passed directly via CLI flags.
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <CheckCircleIcon className="h-5 w-5 text-emerald-500 mt-0.5 flex-shrink-0" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-white">Optional YAML Config</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-300">
                        YAML configuration files are optional and provide default values that can be overridden by CLI arguments.
                      </p>
                    </div>
                  </div>
                  <div className="flex items-start space-x-3">
                    <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500 mt-0.5 flex-shrink-0" />
                    <div>
                      <h3 className="font-semibold text-gray-900 dark:text-white">CLI Overrides YAML</h3>
                      <p className="text-sm text-gray-600 dark:text-gray-300">
                        Command-line arguments always take precedence over configuration file settings.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* YAML Configuration */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">YAML Configuration</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Default Configuration File</h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Create a configuration file at <code className="bg-gray-100 dark:bg-gray-700 px-1.5 py-0.5 rounded text-sm text-gray-800 dark:text-gray-200">configs/default.yaml</code>:
                </p>
                <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                  <code className="text-sm text-emerald-400 block overflow-x-auto">
{`# SubdomainX Configuration File
# All options are optional and can be overridden via CLI

# Input configuration
wildcard_file: ""  # Path to domains file (required via CLI)

# Output configuration
unique_name: "scan"
output_format: "json"  # json, txt, html
output_dir: "output"

# Performance settings
threads: 10
retries: 3
timeout: 30
rate_limit: 100

# Tool selection (all false by default - use CLI flags)
tools:
  subfinder: false
  amass: false
  findomain: false
  assetfinder: false
  sublist3r: false
  knockpy: false
  dnsrecon: false
  fierce: false
  massdns: false
  altdns: false

# Scanning tools
scanners:
  httpx: false
  smap: false`}
                  </code>
                </div>
              </div>
            </div>

            {/* Configuration Parameters */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Configuration Parameters</h2>
              
              {/* Input Configuration */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Input Configuration</h3>
                <div className="space-y-4">
                  {[
                    {
                      name: 'wildcard_file',
                      type: 'string',
                      default: '""',
                      description: 'Path to file containing target domains (one per line)',
                      required: true,
                      cli: '--wildcard'
                    }
                  ].map((param, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between mb-2">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{param.name}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">{param.type}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">{param.description}</p>
                      <div className="flex items-center space-x-4 text-xs">
                        <span className="text-gray-500 dark:text-gray-400">Default: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.default}</code></span>
                        <span className="text-gray-500 dark:text-gray-400">CLI: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.cli}</code></span>
                        {param.required && <span className="text-red-600 dark:text-red-400 font-medium">Required</span>}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Output Configuration */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Output Configuration</h3>
                <div className="space-y-4">
                  {[
                    {
                      name: 'unique_name',
                      type: 'string',
                      default: '"scan"',
                      description: 'Unique name for output files',
                      cli: '--name'
                    },
                    {
                      name: 'output_format',
                      type: 'string',
                      default: '"json"',
                      description: 'Output format: json, txt, html',
                      cli: '--format'
                    },
                    {
                      name: 'output_dir',
                      type: 'string',
                      default: '"output"',
                      description: 'Output directory for generated files',
                      cli: '--output'
                    }
                  ].map((param, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between mb-2">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{param.name}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">{param.type}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">{param.description}</p>
                      <div className="flex items-center space-x-4 text-xs">
                        <span className="text-gray-500 dark:text-gray-400">Default: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.default}</code></span>
                        <span className="text-gray-500 dark:text-gray-400">CLI: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.cli}</code></span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Performance Configuration */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Performance Configuration</h3>
                <div className="space-y-4">
                  {[
                    {
                      name: 'threads',
                      type: 'integer',
                      default: '10',
                      description: 'Number of concurrent threads',
                      cli: '--threads'
                    },
                    {
                      name: 'retries',
                      type: 'integer',
                      default: '3',
                      description: 'Number of retry attempts',
                      cli: '--retries'
                    },
                    {
                      name: 'timeout',
                      type: 'integer',
                      default: '30',
                      description: 'Timeout in seconds',
                      cli: '--timeout'
                    },
                    {
                      name: 'rate_limit',
                      type: 'integer',
                      default: '100',
                      description: 'Rate limit per second',
                      cli: '--rate-limit'
                    }
                  ].map((param, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between mb-2">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{param.name}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">{param.type}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mb-2">{param.description}</p>
                      <div className="flex items-center space-x-4 text-xs">
                        <span className="text-gray-500 dark:text-gray-400">Default: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.default}</code></span>
                        <span className="text-gray-500 dark:text-gray-400">CLI: <code className="bg-gray-100 dark:bg-gray-700 px-1 py-0.5 rounded">{param.cli}</code></span>
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Tool Configuration */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Tool Configuration</h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Tool selection is primarily controlled via CLI flags. YAML configuration provides default states:
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <h4 className="font-medium text-gray-900 dark:text-white mb-2">Enumeration Tools</h4>
                    <div className="space-y-2">
                      {[
                        'subfinder', 'amass', 'findomain', 'assetfinder', 'sublist3r',
                        'knockpy', 'dnsrecon', 'fierce', 'massdns', 'altdns'
                      ].map((tool) => (
                        <div key={tool} className="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded">
                          <code className="text-xs font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{tool}</code>
                          <span className="text-xs text-gray-500 dark:text-gray-400">--{tool}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                  <div>
                    <h4 className="font-medium text-gray-900 dark:text-white mb-2">Scanning Tools</h4>
                    <div className="space-y-2">
                      {[
                        'httpx', 'smap'
                      ].map((tool) => (
                        <div key={tool} className="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded">
                          <code className="text-xs font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{tool}</code>
                          <span className="text-xs text-gray-500 dark:text-gray-400">--{tool}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Usage Examples */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Usage Examples</h2>
              <div className="space-y-6">
                {[
                  {
                    title: 'Use Default Configuration',
                    description: 'Run with default settings from config file',
                    command: 'subdomainx --wildcard domains.txt'
                  },
                  {
                    title: 'Override Configuration',
                    description: 'Override specific settings via CLI',
                    command: 'subdomainx --wildcard domains.txt --threads 20 --format html'
                  },
                  {
                    title: 'Custom Config File',
                    description: 'Use a custom configuration file',
                    command: 'subdomainx --wildcard domains.txt --config my-config.yaml'
                  },
                  {
                    title: 'CLI Only',
                    description: 'Ignore config file and use only CLI arguments',
                    command: 'subdomainx --wildcard domains.txt --subfinder --httpx --format json'
                  }
                ].map((example, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">{example.title}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      {example.description}
                    </p>
                    <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                      <code className="text-sm text-emerald-400 block overflow-x-auto">
                        {example.command}
                      </code>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Tips */}
            <div className="bg-emerald-50 dark:bg-emerald-900/20 rounded-lg p-6 border border-emerald-200 dark:border-emerald-800">
              <div className="flex items-start">
                <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-3 flex-shrink-0" />
                <div>
                  <h3 className="text-lg font-semibold text-emerald-900 dark:text-emerald-100 mb-2">
                    Configuration Tips
                  </h3>
                  <ul className="text-sm text-emerald-800 dark:text-emerald-200 space-y-1">
                    <li>• CLI arguments always override YAML configuration</li>
                    <li>• Use YAML for default settings and CLI for specific overrides</li>
                    <li>• Configuration files are optional - you can use CLI only</li>
                    <li>• Keep configuration files in version control for team consistency</li>
                    <li>• Use environment variables for sensitive configuration</li>
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
