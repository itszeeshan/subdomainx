export default {
  logo: <span>SubdomainX</span>,
  project: {
    link: 'https://github.com/itszeeshan/subdomainx',
  },
  docsRepositoryBase: 'https://github.com/itszeeshan/subdomainx/tree/main/docs',
  footer: {
    text: 'SubdomainX Documentation',
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
      <meta property="og:description" content="All-in-one subdomain enumeration tool" />
      <link rel="stylesheet" href="/styles/globals.css" />
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
