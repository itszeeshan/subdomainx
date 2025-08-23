import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'SubdomainX Documentation',
  description: 'Comprehensive documentation for SubdomainX - All-in-one subdomain enumeration tool',
  keywords: 'subdomain, enumeration, security, pentesting, cli, tool',
  authors: [{ name: 'Muhammad Zeeshan' }],
  openGraph: {
    title: 'SubdomainX Documentation',
    description: 'All-in-one subdomain enumeration tool documentation',
    type: 'website',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" className="h-full">
      <body className={`${inter.className} h-full bg-gray-50`}>
        {children}
      </body>
    </html>
  )
}
