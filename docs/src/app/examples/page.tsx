'use client'
import Navigation from '@/components/Navigation'
import { 
  PlayIcon, 
  CheckCircleIcon,
  InformationCircleIcon,
  ExclamationTriangleIcon
} from '@heroicons/react/24/outline'
import { useState } from 'react'

export default function Examples() {
  const [copiedCommand, setCopiedCommand] = useState<string | null>(null);

  const copyToClipboard = (text: string, commandName: string) => {
    navigator.clipboard.writeText(text);
    setCopiedCommand(commandName);
    setTimeout(() => setCopiedCommand(null), 2000);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950">
      <Navigation />
      
      <div className="lg:pl-64">
        <div className="px-4 py-10 sm:px-6 lg:px-8">
          <div className="mx-auto max-w-4xl">
            {/* Header */}
            <div className="mb-8 flex items-center">
              <div className="bg-emerald-100 dark:bg-emerald-900/20 p-3 rounded-lg mr-4">
                <PlayIcon className="h-8 w-8 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Usage Examples</h1>
                <p className="mt-2 text-lg text-gray-600 dark:text-gray-300">
                  Practical examples and use cases for SubdomainX.
                </p>
              </div>
            </div>

            {/* Basic Examples */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Basic Examples</h2>
              <div className="space-y-6">
                {[
                  {
                    title: 'Simple Scan',
                    description: 'Run a basic scan with all available tools on a single domain',
                    command: 'subdomainx --wildcard domains.txt',
                    notes: 'Uses all available enumeration tools and generates JSON output by default'
                  },
                  {
                    title: 'Specific Tools',
                    description: 'Use only specific enumeration tools for targeted scanning',
                    command: 'subdomainx --wildcard domains.txt --subfinder --amass --findomain',
                    notes: 'Limits the scan to only the specified tools for faster execution'
                  },
                  {
                    title: 'With HTTP Scanning',
                    description: 'Include HTTP scanning to discover web services',
                    command: 'subdomainx --wildcard domains.txt --httpx',
                    notes: 'Discovers web services, extracts titles, status codes, and technologies'
                  },
                  {
                    title: 'Complete Scan',
                    description: 'Full enumeration with both HTTP and port scanning',
                    command: 'subdomainx --wildcard domains.txt --httpx --smap',
                    notes: 'Comprehensive scan including port discovery and service identification'
                  }
                ].map((example, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">{example.title}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                      {example.description}
                    </p>
                    <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4 mb-3">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button 
                          className="text-xs bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-gray-300 px-2 py-1 rounded transition-colors"
                          onClick={() => copyToClipboard(example.command, `basic-${index}`)}
                        >
                          {copiedCommand === `basic-${index}` ? 'Copied!' : 'Copy'}
                        </button>
                      </div>
                      <code className="text-sm text-emerald-400 block overflow-x-auto">
                        {example.command}
                      </code>
                    </div>
                    <div className="flex items-start">
                      <InformationCircleIcon className="h-4 w-4 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-2 flex-shrink-0" />
                      <p className="text-xs text-gray-600 dark:text-gray-300">{example.notes}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Advanced Examples */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Advanced Examples</h2>
              <div className="space-y-6">
                {[
                  {
                    title: 'Custom Output Format',
                    description: 'Generate a beautiful HTML report with custom naming',
                    command: 'subdomainx --wildcard domains.txt --format html --name my_scan --output reports/',
                    notes: 'Creates an HTML report in the reports/ directory with custom filename'
                  },
                  {
                    title: 'High Performance',
                    description: 'Optimize for speed with increased threads and timeout',
                    command: 'subdomainx --wildcard domains.txt --threads 20 --timeout 60 --rate-limit 200',
                    notes: 'Uses more threads and higher rate limits for faster scanning'
                  },
                  {
                    title: 'Verbose Output',
                    description: 'Get detailed information about the scanning process',
                    command: 'subdomainx --wildcard domains.txt --verbose',
                    notes: 'Shows detailed progress and debugging information'
                  },
                  {
                    title: 'Custom Configuration',
                    description: 'Use a custom configuration file',
                    command: 'subdomainx --wildcard domains.txt --config my-config.yaml',
                    notes: 'Loads settings from a custom YAML configuration file'
                  }
                ].map((example, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">{example.title}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                      {example.description}
                    </p>
                    <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4 mb-3">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button 
                          className="text-xs bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-gray-300 px-2 py-1 rounded transition-colors"
                          onClick={() => copyToClipboard(example.command, `advanced-${index}`)}
                        >
                          {copiedCommand === `advanced-${index}` ? 'Copied!' : 'Copy'}
                        </button>
                      </div>
                      <code className="text-sm text-emerald-400 block overflow-x-auto">
                        {example.command}
                      </code>
                    </div>
                    <div className="flex items-start">
                      <InformationCircleIcon className="h-4 w-4 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-2 flex-shrink-0" />
                      <p className="text-xs text-gray-600 dark:text-gray-300">{example.notes}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Sample Files */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Sample Files</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">domains.txt</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    A simple text file containing target domains:
                  </p>
                  <div className="bg-gray-50 dark:bg-gray-800 rounded p-3">
                    <code className="text-xs text-gray-800 dark:text-gray-200 block">
                      example.com<br/>
                      test.com<br/>
                      demo.org<br/>
                      sample.net
                    </code>
                  </div>
                </div>

                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">config.yaml</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Custom configuration file:
                  </p>
                  <div className="bg-gray-50 dark:bg-gray-800 rounded p-3">
                    <code className="text-xs text-gray-800 dark:text-gray-200 block">
                      threads: 15<br/>
                      timeout: 45<br/>
                      retries: 5<br/>
                      rate_limit: 150<br/>
                      output_format: html<br/>
                      output_dir: scans/
                    </code>
                  </div>
                </div>
              </div>
            </div>

            {/* Best Practices */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Best Practices</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="space-y-4">
                  {[
                    {
                      title: 'Start Small',
                      description: 'Begin with a few domains and specific tools to test your setup',
                      icon: CheckCircleIcon
                    },
                    {
                      title: 'Use Appropriate Rate Limits',
                      description: 'Set reasonable rate limits to avoid being blocked by services',
                      icon: ExclamationTriangleIcon
                    },
                    {
                      title: 'Monitor Resources',
                      description: 'Keep an eye on system resources when using high thread counts',
                      icon: InformationCircleIcon
                    },
                    {
                      title: 'Organize Output',
                      description: 'Use meaningful names and organize output directories',
                      icon: CheckCircleIcon
                    },
                    {
                      title: 'Verify Results',
                      description: 'Always verify discovered subdomains and validate findings',
                      icon: ExclamationTriangleIcon
                    }
                  ].map((practice, index) => (
                    <div key={index} className="flex items-start space-x-3">
                      <practice.icon className="h-5 w-5 text-emerald-500 mt-0.5 flex-shrink-0" />
                      <div>
                        <h3 className="font-semibold text-gray-900 dark:text-white">{practice.title}</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300">{practice.description}</p>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Common Use Cases */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Common Use Cases</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Bug Bounty</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Comprehensive subdomain discovery for bug bounty programs:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --wildcard targets.txt --httpx --smap --format html --name bugbounty_scan
                    </code>
                  </div>
                </div>

                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Security Assessment</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Quick reconnaissance for security assessments:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --wildcard domains.txt --subfinder --amass --httpx --threads 10
                    </code>
                  </div>
                </div>

                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Asset Discovery</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Discover all assets belonging to an organization:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --wildcard org_domains.txt --format json --output assets/
                    </code>
                  </div>
                </div>

                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Research</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Academic or research purposes with detailed output:
                  </p>
                  <div className="bg-gray-900 dark:bg-gray-800 rounded-lg p-4">
                    <code className="text-sm text-emerald-400 block">
                      subdomainx --wildcard research.txt --verbose --format json --name research_data
                    </code>
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
