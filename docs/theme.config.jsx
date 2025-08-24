export default {
  logo: (
    <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
      <img src="/logo.png" alt="SubdomainX" style={{ height: '32px', width: 'auto' }} />
      <span>SubdomainX</span>
      <span style={{ 
        fontSize: '12px', 
        padding: '2px 6px', 
        backgroundColor: '#16a34a', 
        color: 'white', 
        borderRadius: '12px',
        fontWeight: '500'
      }}>
        v1.4.1
      </span>
    </div>
  ),
  project: {
    link: 'https://github.com/itszeeshan/subdomainx',
  },
  docsRepositoryBase: 'https://github.com/itszeeshan/subdomainx/tree/main/docs',
  footer: {
    text: 'Made with ❤️ by Zeeshan • Star us on GitHub',
    component: () => (
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '10px', marginBottom: '10px', gap: '8px', fontSize: '14px', color: 'var(--tw-prose-body)' }}>
        <span>Made with ❤️ by Zeeshan</span>
        <span>•</span>
        <a 
          href="https://github.com/itszeeshan/subdomainx" 
          target="_blank" 
          rel="noopener noreferrer"
          style={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '4px',
            textDecoration: 'none',
            transition: 'color 0.2s ease'
          }}
          onMouseEnter={(e) => e.target.style.color = '#16a34a'}
          onMouseLeave={(e) => e.target.style.color = 'var(--tw-prose-body)'}
        >
          <svg style={{ width: '12px', height: '12px' }} fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
          Star us
        </a>
      </div>
    ),
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
    toggleButton: true,
  },
  navigation: {
    prev: false,
    next: true,
  },
  useNextSeoProps() {
    return {
      titleTemplate: '%s – SubdomainX'
    }
  },
  head: (
    <>
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta property="og:title" content="SubdomainX" />
      <meta property="og:description" content="Advanced Subdomain Discovery & Security Reconnaissance Tool" />
      <link rel="stylesheet" href="/styles/globals.css" />
      <link rel="icon" type="image/png" href="/logo.png" />
      <link rel="apple-touch-icon" href="/logo.png" />
    </>
  ),
  primaryHue: {
    dark: 160,
    light: 160
  },
  primarySaturation: {
    dark: 100,
    light: 100
  },
  nextThemes: {
    defaultTheme: 'system'
  },
  components: {
    callout: {
      className: 'nx-mt-6 nx-rounded-lg nx-border nx-p-4 [&>*:first-child]:nx-mt-0 [&>*:last-child]:nx-mb-0',
      titleClassName: 'nx-text-sm nx-font-medium',
      infoIcon: (
        <span className="nx-inline-flex nx-rounded-lg nx-bg-blue-500/10 nx-px-2 nx-py-1 nx-text-xs nx-font-medium nx-text-blue-600 nx-ring-1 nx-ring-inset nx-ring-blue-500/20 dark:nx-bg-blue-400/10 dark:nx-text-blue-400 dark:nx-ring-blue-400/20">
          Info
        </span>
      ),
      warningIcon: (
        <span className="nx-inline-flex nx-rounded-lg nx-bg-yellow-500/10 nx-px-2 nx-py-1 nx-text-xs nx-font-medium nx-text-yellow-600 nx-ring-1 nx-ring-inset nx-ring-yellow-500/20 dark:nx-bg-yellow-400/10 dark:nx-text-yellow-400 dark:nx-ring-yellow-400/20">
          Warning
        </span>
      ),
      errorIcon: (
        <span className="nx-inline-flex nx-rounded-lg nx-bg-red-500/10 nx-px-2 nx-py-1 nx-text-xs nx-font-medium nx-text-red-600 nx-ring-1 nx-ring-inset nx-ring-red-500/20 dark:nx-bg-red-400/10 dark:nx-text-red-400 dark:nx-ring-red-400/20">
          Error
        </span>
      ),
      tipIcon: (
        <span className="nx-inline-flex nx-rounded-lg nx-bg-green-500/10 nx-px-2 nx-py-1 nx-text-xs nx-font-medium nx-text-green-600 nx-ring-1 nx-ring-inset nx-ring-green-500/20 dark:nx-bg-green-400/10 dark:nx-text-green-400 dark:nx-ring-green-400/20">
          Tip
        </span>
      ),
      noteIcon: (
        <span className="nx-inline-flex nx-rounded-lg nx-bg-purple-500/10 nx-px-2 nx-py-1 nx-text-xs nx-font-medium nx-text-purple-600 nx-ring-1 nx-ring-inset nx-ring-purple-500/20 dark:nx-bg-purple-400/10 dark:nx-text-purple-400 dark:nx-ring-purple-400/20">
          Note
        </span>
      )
    }
  }
}
