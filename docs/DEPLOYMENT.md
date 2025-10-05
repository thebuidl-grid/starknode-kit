# Documentation Deployment Guide

This guide covers how to deploy the Starknode Kit documentation using GitBook and other platforms.

## GitBook Deployment

### Option 1: GitBook.com (Recommended)

1. **Create GitBook Account**
   - Go to [GitBook.com](https://www.gitbook.com)
   - Sign up with your GitHub account

2. **Import Repository**
   - Click "New Space"
   - Select "Import from GitHub"
   - Choose `thebuidl-grid/starknode-kit`
   - Select the `docs` folder as the root

3. **Configure Settings**
   - Set title: "Starknode Kit Documentation"
   - Set description: "Complete guide for Ethereum and Starknet node management"
   - Enable public access

4. **Custom Domain (Optional)**
   - Go to Space Settings > Domains
   - Add your custom domain (e.g., `docs.starknode-kit.com`)
   - Update DNS records as instructed

### Option 2: Self-Hosted GitBook

1. **Install GitBook CLI**
   ```bash
   npm install -g gitbook-cli
   ```

2. **Install Dependencies**
   ```bash
   cd docs
   npm install
   gitbook install
   ```

3. **Build Documentation**
   ```bash
   gitbook build
   ```

4. **Serve Locally**
   ```bash
   gitbook serve
   ```

5. **Deploy to Web Server**
   ```bash
   # Copy built files to web server
   rsync -av _book/ user@server:/var/www/docs/
   ```

## GitHub Pages Deployment

### Automatic Deployment with GitHub Actions

1. **Create GitHub Action**
   ```yaml
   # .github/workflows/docs.yml
   name: Deploy Documentation
   
   on:
     push:
       branches: [ main ]
       paths: [ 'docs/**' ]
   
   jobs:
     deploy:
       runs-on: ubuntu-latest
       steps:
         - uses: actions/checkout@v3
         
         - name: Setup Node.js
           uses: actions/setup-node@v3
           with:
             node-version: '16'
             
         - name: Install GitBook
           run: npm install -g gitbook-cli
           
         - name: Build Documentation
           run: |
             cd docs
             gitbook install
             gitbook build
             
         - name: Deploy to GitHub Pages
           uses: peaceiris/actions-gh-pages@v3
           with:
             github_token: ${{ secrets.GITHUB_TOKEN }}
             publish_dir: ./docs/_book
   ```

2. **Enable GitHub Pages**
   - Go to repository Settings > Pages
   - Select "GitHub Actions" as source
   - The documentation will be available at `https://thebuidl-grid.github.io/starknode-kit`

## Netlify Deployment

1. **Connect Repository**
   - Go to [Netlify](https://netlify.com)
   - Connect your GitHub repository

2. **Build Settings**
   - Build command: `cd docs && gitbook build`
   - Publish directory: `docs/_book`

3. **Deploy**
   - Netlify will automatically build and deploy
   - Custom domain can be configured in site settings

## Vercel Deployment

1. **Install Vercel CLI**
   ```bash
   npm install -g vercel
   ```

2. **Configure Project**
   ```json
   // vercel.json
   {
     "builds": [
       {
         "src": "docs/package.json",
         "use": "@vercel/static-build",
         "config": {
           "distDir": "_book"
         }
       }
     ]
   }
   ```

3. **Deploy**
   ```bash
   cd docs
   vercel --prod
   ```

## Docker Deployment

1. **Create Dockerfile**
   ```dockerfile
   FROM node:16-alpine
   
   WORKDIR /app
   
   # Install GitBook
   RUN npm install -g gitbook-cli
   
   # Copy documentation
   COPY docs/ ./docs/
   
   # Build documentation
   WORKDIR /app/docs
   RUN gitbook install && gitbook build
   
   # Serve with nginx
   FROM nginx:alpine
   COPY --from=0 /app/docs/_book /usr/share/nginx/html
   EXPOSE 80
   ```

2. **Build and Run**
   ```bash
   docker build -t starknode-kit-docs .
   docker run -p 8080:80 starknode-kit-docs
   ```

## Continuous Integration

### GitHub Actions Workflow

```yaml
name: Documentation CI/CD

on:
  push:
    branches: [ main ]
    paths: [ 'docs/**' ]
  pull_request:
    branches: [ main ]
    paths: [ 'docs/**' ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '16'
          
      - name: Install dependencies
        run: |
          cd docs
          npm install
          npm install -g gitbook-cli
          
      - name: Build documentation
        run: |
          cd docs
          gitbook install
          gitbook build
          
      - name: Test build
        run: |
          cd docs
          ls -la _book/
          
  deploy:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v3
      
      - name: Deploy to production
        run: |
          # Add deployment commands here
          echo "Deploying to production..."
```

## Monitoring and Analytics

### Google Analytics

Add to `book.json`:
```json
{
  "pluginsConfig": {
    "google-analytics": {
      "token": "GA_TRACKING_ID"
    }
  }
}
```

### Search Engine Optimization

1. **Add meta tags** to `README.md`
2. **Create sitemap** for better indexing
3. **Use descriptive URLs** and headings
4. **Add structured data** for rich snippets

## Maintenance

### Regular Updates

1. **Content Updates**
   - Review and update documentation monthly
   - Keep examples current with latest versions
   - Update screenshots and diagrams

2. **Technical Updates**
   - Update GitBook plugins regularly
   - Test builds on different platforms
   - Monitor for broken links

3. **Performance Monitoring**
   - Monitor page load times
   - Check for broken links
   - Review user feedback

### Backup Strategy

1. **Repository Backup**
   - Documentation is stored in Git repository
   - Regular backups to multiple locations

2. **Build Artifacts**
   - Store built documentation as releases
   - Keep multiple versions for rollback

## Troubleshooting

### Common Issues

**Build Failures**
```bash
# Clear cache and rebuild
rm -rf node_modules
npm install
gitbook install
gitbook build
```

**Plugin Issues**
```bash
# Update plugins
gitbook update
```

**Deployment Issues**
- Check build logs
- Verify file permissions
- Test locally first

## Getting Help

For deployment issues:

- Check GitBook documentation
- Review platform-specific guides
- Join our [Telegram community](https://t.me/+SCPbza9fk8dkYWI0)
- Open an issue on [GitHub](https://github.com/thebuidl-grid/starknode-kit/issues)
