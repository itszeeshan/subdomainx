# SubdomainX Documentation

This is the documentation site for SubdomainX, built with Next.js and deployed on Vercel.

## Features

- 📚 Comprehensive documentation
- 🎨 Beautiful, modern UI inspired by Supabase
- 📱 Fully responsive design
- ⚡ Fast and optimized
- 🔍 Search functionality (coming soon)
- 📄 Multiple output formats

## Development

### Prerequisites

- Node.js 18+ 
- npm or yarn

### Getting Started

1. Install dependencies:
   ```bash
   npm install
   ```

2. Run the development server:
   ```bash
   npm run dev
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser.

### Building for Production

```bash
npm run build
```

### Deployment

This site is configured for deployment on Vercel. Simply connect your GitHub repository to Vercel and it will automatically deploy.

## Project Structure

```
docs/
├── src/
│   ├── app/                 # Next.js app directory
│   │   ├── page.tsx         # Home page
│   │   ├── installation/    # Installation guide
│   │   ├── cli-reference/   # CLI documentation
│   │   ├── supported-tools/ # Supported tools
│   │   ├── configuration/   # Configuration guide
│   │   ├── examples/        # Usage examples
│   │   └── api-reference/   # API documentation
│   ├── components/          # Reusable components
│   └── styles/             # Global styles
├── public/                 # Static assets
├── vercel.json            # Vercel configuration
└── package.json           # Dependencies
```

## Technologies Used

- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Heroicons** - Icons
- **Framer Motion** - Animations (coming soon)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test locally
5. Submit a pull request

## License

MIT License - see the main repository for details.
