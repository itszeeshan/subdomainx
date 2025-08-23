'use client'
import Navigation from '@/components/Navigation'
import { 
  CommandLineIcon, 
  InformationCircleIcon,
  ExclamationTriangleIcon
} from '@heroicons/react/24/outline'
import { useState } from 'react'

export default function CLIReference() {
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
                <CommandLineIcon className="h-8 w-8 text-emerald-600 dark:text-emerald-400" />
              </div>
              <div>
                <h1 className="text-3xl font-bold text-gray-900 dark:text-white">CLI Reference</h1>
                <p className="mt-2 text-lg text-gray-600 dark:text-gray-300">
                  Complete command-line interface reference for SubdomainX.
                </p>
              </div>
            </div>

            {/* Basic Usage */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-4">Basic Usage</h2>
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <div className="rounded-lg bg-gray-900 dark:bg-gray-800 p-4 relative">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-gray-300">Basic Command</span>
                    <button 
                      className="text-xs bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-gray-300 px-2 py-1 rounded transition-colors"
                      onClick={() => copyToClipboard('subdomainx --wildcard <domains_file> [OPTIONS]', 'basic')}
                    >
                      {copiedCommand === 'basic' ? 'Copied!' : 'Copy'}
                    </button>
                  </div>
                  <code className="text-sm text-emerald-400 block overflow-x-auto">
                    subdomainx --wildcard &lt;domains_file&gt; [OPTIONS]
                  </code>
                </div>
                
                <div className="mt-4 text-sm text-gray-600 dark:text-gray-300">
                  <p>
                    <strong>Required:</strong> You must provide a file containing target domains using the <code className="bg-gray-100 dark:bg-gray-700 px-1.5 py-0.5 rounded text-sm text-gray-800 dark:text-gray-200">--wildcard</code> flag.
                  </p>
                </div>
              </div>
            </div>

            {/* Command Options */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Command Options</h2>
              
              {/* Required Options */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center">
                  <ExclamationTriangleIcon className="h-5 w-5 text-red-500 mr-2" />
                  Required Options
                </h3>
                <div className="space-y-4">
                  <div className="border-l-4 border-red-500 pl-4 py-1">
                    <div className="flex items-center justify-between">
                      <code className="text-sm font-mono bg-red-50 dark:bg-red-900/20 px-2 py-1 rounded text-red-700 dark:text-red-400">--wildcard FILE</code>
                      <span className="text-xs text-red-600 dark:text-red-400 font-medium bg-red-100 dark:bg-red-900/30 px-2 py-0.5 rounded">Required</span>
                    </div>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mt-1 ml-1">
                      Path to file containing target domains (one per line)
                    </p>
                  </div>
                </div>
              </div>

              {/* Tool Selection */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Tool Selection</h3>
                <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
                  Use specific tools, otherwise all available tools will be used.
                </p>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="space-y-3">
                    {[
                      { command: '--subfinder', desc: 'Use subfinder tool' },
                      { command: '--amass', desc: 'Use amass tool' },
                      { command: '--findomain', desc: 'Use findomain tool' },
                      { command: '--assetfinder', desc: 'Use assetfinder tool' },
                      { command: '--sublist3r', desc: 'Use sublist3r tool' }
                    ].map((tool, index) => (
                      <div key={index} className="p-2 bg-gray-50 dark:bg-gray-800 rounded">
                        <div className="flex items-center justify-between">
                          <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{tool.command}</code>
                          <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">Flag</span>
                        </div>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-1 ml-1">{tool.desc}</p>
                      </div>
                    ))}
                  </div>
                  
                  <div className="space-y-3">
                    {[
                      { command: '--knockpy', desc: 'Use knockpy tool' },
                      { command: '--dnsrecon', desc: 'Use dnsrecon tool' },
                      { command: '--fierce', desc: 'Use fierce tool' },
                      { command: '--massdns', desc: 'Use massdns tool' },
                      { command: '--altdns', desc: 'Use altdns tool' }
                    ].map((tool, index) => (
                      <div key={index} className="p-2 bg-gray-50 dark:bg-gray-800 rounded">
                        <div className="flex items-center justify-between">
                          <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{tool.command}</code>
                          <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">Flag</span>
                        </div>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-1 ml-1">{tool.desc}</p>
                      </div>
                    ))}
                  </div>
                </div>
                
                <div className="mt-4 p-3 bg-emerald-50 dark:bg-emerald-900/20 rounded-lg border border-emerald-200 dark:border-emerald-800">
                  <div className="flex items-start">
                    <InformationCircleIcon className="h-5 w-5 text-emerald-600 dark:text-emerald-400 mt-0.5 mr-2 flex-shrink-0" />
                    <div className="text-sm text-emerald-700 dark:text-emerald-300">
                      <strong>Note:</strong> If no specific tools are specified, SubdomainX will use all available tools on your system.
                    </div>
                  </div>
                </div>
              </div>

              {/* Scanning Tools */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Scanning Tools</h3>
                <div className="space-y-4">
                  {[
                    { 
                      command: '--httpx', 
                      type: 'Flag',
                      desc: 'Use httpx for HTTP scanning (discovers web services, extracts titles, status codes, and technologies)' 
                    },
                    { 
                      command: '--smap', 
                      type: 'Flag',
                      desc: 'Use smap for port scanning (identifies open ports and services on discovered hosts)' 
                    }
                  ].map((tool, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{tool.command}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">{tool.type}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mt-2">
                        {tool.desc}
                      </p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Output Options */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Output Options</h3>
                <div className="space-y-4">
                  {[
                    { 
                      command: '--name NAME', 
                      default: 'scan',
                      desc: 'Unique name for output files' 
                    },
                    { 
                      command: '--format FORMAT', 
                      default: 'json',
                      desc: 'Output format: json, txt, html' 
                    },
                    { 
                      command: '--output DIR', 
                      default: 'output',
                      desc: 'Output directory for generated files' 
                    }
                  ].map((option, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{option.command}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">Default: {option.default}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mt-2">
                        {option.desc}
                      </p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Performance Options */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 mb-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Performance Options</h3>
                <div className="space-y-4">
                  {[
                    { 
                      command: '--threads N', 
                      default: '10',
                      desc: 'Number of concurrent threads' 
                    },
                    { 
                      command: '--retries N', 
                      default: '3',
                      desc: 'Number of retry attempts' 
                    },
                    { 
                      command: '--timeout N', 
                      default: '30',
                      desc: 'Timeout in seconds' 
                    },
                    { 
                      command: '--rate-limit N', 
                      default: '100',
                      desc: 'Rate limit per second' 
                    }
                  ].map((option, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{option.command}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">Default: {option.default}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mt-2">
                        {option.desc}
                      </p>
                    </div>
                  ))}
                </div>
              </div>

              {/* Utility Options */}
              <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4">Utility Options</h3>
                <div className="space-y-4">
                  {[
                    { 
                      command: '--help', 
                      type: 'Flag',
                      desc: 'Show help message' 
                    },
                    { 
                      command: '--version', 
                      type: 'Flag',
                      desc: 'Show version information' 
                    },
                    { 
                      command: '--check-tools', 
                      type: 'Flag',
                      desc: 'Check availability of enumeration tools' 
                    },
                    { 
                      command: '--install-tools', 
                      type: 'Flag',
                      desc: 'Show installation instructions for missing tools' 
                    },
                    { 
                      command: '--config FILE', 
                      type: 'Optional',
                      desc: 'Use custom configuration file' 
                    },
                    { 
                      command: '-v, --verbose', 
                      type: 'Flag',
                      desc: 'Enable verbose output' 
                    }
                  ].map((option, index) => (
                    <div key={index} className="p-3 bg-gray-50 dark:bg-gray-800 rounded-lg">
                      <div className="flex items-center justify-between">
                        <code className="text-sm font-mono bg-white dark:bg-gray-700 px-2 py-1 rounded border border-gray-200 dark:border-gray-600 text-gray-800 dark:text-gray-200">{option.command}</code>
                        <span className="text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">{option.type}</span>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 mt-2">
                        {option.desc}
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Examples */}
            <div className="mb-12">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">Examples</h2>
              
              <div className="space-y-6">
                {[
                  {
                    title: 'Basic Scan',
                    desc: 'Run a basic scan with all available tools:',
                    command: 'subdomainx --wildcard domains.txt'
                  },
                  {
                    title: 'Specific Tools',
                    desc: 'Use only specific enumeration and scanning tools:',
                    command: 'subdomainx --wildcard domains.txt --amass --subfinder --httpx'
                  },
                  {
                    title: 'HTML Report',
                    desc: 'Generate a beautiful HTML report with custom name:',
                    command: 'subdomainx --wildcard domains.txt --format html --name my_scan'
                  },
                  {
                    title: 'High Performance',
                    desc: 'Optimize for speed with more threads and longer timeout:',
                    command: 'subdomainx --wildcard domains.txt --threads 20 --timeout 60'
                  }
                ].map((example, index) => (
                  <div key={index} className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">{example.title}</h3>
                    <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                      {example.desc}
                    </p>
                    <div className="rounded-lg bg-gray-900 dark:bg-gray-800 p-4">
                      <div className="flex items-center justify-between mb-2">
                        <span className="text-sm font-medium text-gray-300">Command</span>
                        <button 
                          className="text-xs bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-gray-300 px-2 py-1 rounded transition-colors"
                          onClick={() => copyToClipboard(example.command, `example-${index}`)}
                        >
                          {copiedCommand === `example-${index}` ? 'Copied!' : 'Copy'}
                        </button>
                      </div>
                      <code className="text-sm text-emerald-400 block overflow-x-auto">
                        {example.command}
                      </code>
                    </div>
                  </div>
                ))}

                <div className="bg-white dark:bg-gray-900 rounded-lg shadow-sm p-6 border border-gray-200 dark:border-gray-700">
                  <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-3">Tool Management</h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                    Check which tools are available and get installation help:
                  </p>
                  <div className="space-y-3">
                    {[
                      {
                        title: 'Check Tools',
                        command: 'subdomainx --check-tools'
                      },
                      {
                        title: 'Install Help',
                        command: 'subdomainx --install-tools'
                      }
                    ].map((toolCmd, index) => (
                      <div key={index} className="rounded-lg bg-gray-900 dark:bg-gray-800 p-4">
                        <div className="flex items-center justify-between mb-2">
                          <span className="text-sm font-medium text-gray-300">{toolCmd.title}</span>
                          <button 
                            className="text-xs bg-gray-800 dark:bg-gray-700 hover:bg-gray-700 dark:hover:bg-gray-600 text-gray-300 px-2 py-1 rounded transition-colors"
                            onClick={() => copyToClipboard(toolCmd.command, `tool-${index}`)}
                          >
                            {copiedCommand === `tool-${index}` ? 'Copied!' : 'Copy'}
                          </button>
                        </div>
                        <code className="text-sm text-emerald-400 block overflow-x-auto">
                          {toolCmd.command}
                        </code>
                      </div>
                    ))}
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