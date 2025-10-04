# Starknode-kit Documentation Guide

This guide explains how to run and develop the starknode-kit documentation site.

## Quick Start

### Installation

```bash
cd docs
npm install
```

### Development Server

Run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser to see the documentation.

The page auto-updates as you edit files.

### Production Build

Build for production:

```bash
npm run build
npm start
```

## Documentation Structure

```
docs/
├── src/
│   ├── app/                    # Next.js pages (App Router)
│   │   ├── page.tsx           # Homepage
│   │   ├── layout.tsx         # Root layout with sidebar/header
│   │   ├── globals.css        # Global styles
│   │   ├── getting-started/   # Getting started guide
│   │   ├── installation/      # Installation guide
│   │   ├── configuration/     # Configuration docs
│   │   ├── commands/          # Command reference
│   │   ├── clients/           # Client documentation
│   │   ├── validator/         # Validator setup guide
│   │   ├── requirements/      # Requirements page
│   │   └── contributing/      # Contributing guide
│   └── components/            # Reusable components
│       ├── Sidebar.tsx        # Navigation sidebar
│       ├── Header.tsx         # Top header with search
│       └── CodeBlock.tsx      # Code block component
├── public/                    # Static assets
├── package.json              # Dependencies
└── README.md                 # Documentation README

```

## Features

### GitBook-Style Layout

- **Fixed Sidebar** - Navigation always visible on the left
- **Header** - Search and links at the top
- **Content Area** - Main documentation content
- **Responsive** - Works on mobile and desktop

### Components

#### Sidebar (`components/Sidebar.tsx`)

- Hierarchical navigation
- Active page highlighting
- Collapsible sections
- Links to all documentation pages

#### Header (`components/Header.tsx`)

- Search bar (placeholder, can be enhanced)
- GitHub link
- Telegram link
- Dark mode support

#### CodeBlock (`components/CodeBlock.tsx`)

- Syntax-highlighted code blocks
- Copy-to-clipboard button
- Language support
- Dark theme optimized

### Styling

- **Tailwind CSS** - Utility-first CSS framework
- **Dark Mode** - Automatic based on system preference
- **Custom Scrollbars** - Styled for better UX
- **Prose** - Typography optimized for documentation

## Adding New Pages

### 1. Create Page File

Create a new `page.tsx` in the appropriate directory:

```bash
mkdir -p src/app/your-page
```

```tsx
// src/app/your-page/page.tsx
import CodeBlock from '@/components/CodeBlock';

export default function YourPage() {
  return (
    <div className="prose prose-lg dark:prose-invert max-w-none">
      <h1>Your Page Title</h1>
      <p>Your content here...</p>
      
      <CodeBlock code="your code here" />
    </div>
  );
}
```

### 2. Add to Navigation

Update `src/components/Sidebar.tsx`:

```tsx
const navigation: NavItem[] = [
  // ... existing items
  { title: 'Your Page', href: '/your-page' },
];
```

### 3. Test

```bash
npm run dev
```

Visit `http://localhost:3000/your-page`

## Deployment

### Vercel (Recommended)

1. Push to GitHub
2. Import project in Vercel
3. Deploy automatically

### Self-Hosted

```bash
npm run build
npm start
```

Or use a process manager like PM2:

```bash
pm2 start npm --name "starknode-docs" -- start
```

### Static Export

For static hosting (GitHub Pages, etc.):

```bash
# Add to next.config.ts:
# output: 'export'

npm run build
# Output will be in 'out/' directory
```

## Customization

### Colors

Edit `src/app/globals.css`:

```css
:root {
  --background: #ffffff;
  --foreground: #171717;
}
```

### Fonts

Edit `src/app/layout.tsx`:

```tsx
import { Inter } from "next/font/google";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});
```

### Navigation

Edit `src/components/Sidebar.tsx`:

```tsx
const navigation: NavItem[] = [
  // Add/remove/reorder items
];
```

## Tips

### Use CodeBlock Component

```tsx
import CodeBlock from '@/components/CodeBlock';

<CodeBlock code="your code here" language="bash" />
```

### Use Prose Styling

Always wrap content in prose div:

```tsx
<div className="prose prose-lg dark:prose-invert max-w-none">
  {/* Your content */}
</div>
```

### Add Info Boxes

```tsx
<div className="bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 p-4 my-6">
  <p className="font-semibold mb-2">💡 Tip</p>
  <p className="mb-0">Your tip here</p>
</div>
```

### Link Between Pages

```tsx
import Link from 'next/link';

<Link href="/other-page">Other Page</Link>
```

## Troubleshooting

### Port Already in Use

```bash
# Kill process on port 3000
lsof -i :3000
kill -9 [PID]

# Or use different port
PORT=3001 npm run dev
```

### Build Errors

```bash
# Clean and rebuild
rm -rf .next node_modules
npm install
npm run build
```

### Styling Not Updating

```bash
# Clear Next.js cache
rm -rf .next
npm run dev
```

## Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)
- [React Documentation](https://react.dev/)

## Contributing

See the main [Contributing Guide](../CONTRIBUTING.md) for guidelines on contributing to the documentation.

