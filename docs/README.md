# starknode-kit Documentation

Official documentation for starknode-kit, built with Next.js.

## Quick Start

```bash
# Install dependencies
npm install

# Run development server
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the documentation.

## Features

- 🎨 Light/Dark mode toggle
- 📱 Responsive design
- 🔍 Clean, simple navigation
- 📝 Comprehensive guides
- 💻 Code examples with copy button

## Building

```bash
# Build for production
npm run build

# Start production server
npm start
```

## Documentation Pages

- **Introduction** - Overview and quick start
- **Getting Started** - Step-by-step setup guide  
- **Installation** - Installation methods
- **Configuration** - Node configuration
- **Commands** - CLI command reference
- **Clients** - Supported clients (Geth, Reth, Lighthouse, Prysm, Juno)
- **Validator Setup** - Validator node setup
- **Requirements** - Hardware/software requirements
- **Contributing** - Contribution guidelines

## Project Structure

```
docs/
├── src/
│   ├── app/              # Documentation pages
│   │   ├── page.tsx      # Homepage
│   │   ├── layout.tsx    # Root layout
│   │   ├── globals.css   # Global styles
│   │   └── [pages]/      # Documentation pages
│   └── components/       # Reusable components
│       ├── Sidebar.tsx   # Navigation sidebar
│       ├── Header.tsx    # Top header
│       ├── ThemeToggle.tsx  # Light/dark mode toggle
│       └── CodeBlock.tsx # Code block with copy
└── package.json
```

## Contributing

See the main [Contributing Guide](../README.md#contributing) for guidelines.

## License

MIT License - see LICENSE file for details.
