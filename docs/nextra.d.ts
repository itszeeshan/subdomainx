declare module 'nextra' {
  const nextra: (config: Record<string, unknown>) => (config: Record<string, unknown>) => Record<string, unknown>
  export = nextra
}
