import CodeBlock from '@/components/CodeBlock';
import Link from 'next/link';

export default function Contributing() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Contributing</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        We welcome contributions to starknode-kit! This guide will help you get
        started with contributing to the project.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Ways to Contribute</h2>

      <p>There are many ways to contribute to starknode-kit:</p>

      <ul>
        <li>🐛 <strong>Report bugs</strong> - Help us identify and fix issues</li>
        <li>💡 <strong>Suggest features</strong> - Share your ideas for improvements</li>
        <li>📝 <strong>Improve documentation</strong> - Help others understand and use the tool</li>
        <li>💻 <strong>Write code</strong> - Implement new features or fix bugs</li>
        <li>🧪 <strong>Test</strong> - Test new releases and provide feedback</li>
        <li>🌍 <strong>Community support</strong> - Help other users in issues and discussions</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Getting Started</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">1. Fork the Repository</h3>

      <p>Start by forking the repository on GitHub:</p>

      <CodeBlock code="# Visit https://github.com/thebuidl-grid/starknode-kit and click 'Fork'" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">2. Clone Your Fork</h3>

      <CodeBlock code={`git clone https://github.com/YOUR_USERNAME/starknode-kit.git
cd starknode-kit`} />

      <h3 className="text-2xl font-semibold mt-10 mb-5">3. Set Up Development Environment</h3>

      <p>Make sure you have the required tools:</p>

      <ul>
        <li>Go 1.24 or later</li>
        <li>Make</li>
        <li>Git</li>
      </ul>

      <CodeBlock code={`# Check Go version
go version

# Install dependencies
go mod download

# Build the project
make build`} />

      <h3 className="text-2xl font-semibold mt-10 mb-5">4. Create a Branch</h3>

      <p>Create a branch for your changes:</p>

      <CodeBlock code={`git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Development Workflow</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Project Structure</h3>

      <CodeBlock code={`starknode-kit/
├── cli/              # CLI commands and root setup
│   ├── commands/     # Individual command implementations
│   └── options/      # Global options and config
├── pkg/              # Main package code
│   ├── clients/      # Client implementations (Geth, Reth, etc.)
│   ├── monitoring/   # Monitoring dashboard
│   ├── types/        # Type definitions
│   ├── utils/        # Utility functions
│   └── validator/    # Validator management
├── main.go           # Entry point
├── Makefile          # Build commands
└── README.md         # Project documentation`} />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Making Changes</h3>

      <ol>
        <li><strong>Write your code</strong> - Implement your feature or fix</li>
        <li><strong>Follow Go conventions</strong> - Use <code>gofmt</code> and follow Go best practices</li>
        <li><strong>Add tests</strong> - Write tests for your changes when applicable</li>
        <li><strong>Update documentation</strong> - Update README or add docs as needed</li>
      </ol>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Testing</h3>

      <p>Run tests before submitting:</p>

      <CodeBlock code={`# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./pkg/clients/`} />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Code Style</h3>

      <p>Format your code with <code>gofmt</code>:</p>

      <CodeBlock code={`# Format all Go files
gofmt -w .

# Check formatting
gofmt -l .`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Submitting Changes</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">1. Commit Your Changes</h3>

      <p>Write clear, descriptive commit messages:</p>

      <CodeBlock code={`git add .
git commit -m "feat: add support for new client"

# Follow conventional commits format:
# feat: new feature
# fix: bug fix
# docs: documentation changes
# refactor: code refactoring
# test: adding or updating tests
# chore: maintenance tasks`} />

      <h3 className="text-2xl font-semibold mt-10 mb-5">2. Push to Your Fork</h3>

      <CodeBlock code="git push origin feature/your-feature-name" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">3. Create a Pull Request</h3>

      <ol>
        <li>Go to the original repository on GitHub</li>
        <li>Click "New Pull Request"</li>
        <li>Select your fork and branch</li>
        <li>Fill out the PR template with:
          <ul>
            <li>Description of changes</li>
            <li>Related issues</li>
            <li>Testing performed</li>
            <li>Screenshots (if UI changes)</li>
          </ul>
        </li>
        <li>Submit the pull request</li>
      </ol>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Pull Request Guidelines</h2>

      <p>To increase the chances of your PR being accepted:</p>

      <ul>
        <li>✅ Make focused, single-purpose changes</li>
        <li>✅ Write clear commit messages</li>
        <li>✅ Include tests for new features</li>
        <li>✅ Update documentation as needed</li>
        <li>✅ Follow existing code style</li>
        <li>✅ Respond to review feedback promptly</li>
        <li>❌ Don't include unrelated changes</li>
        <li>❌ Don't submit untested code</li>
        <li>❌ Don't break existing functionality</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Reporting Bugs</h2>

      <p>When reporting bugs, include:</p>

      <ol>
        <li><strong>Description</strong> - Clear description of the bug</li>
        <li><strong>Steps to reproduce</strong> - How to trigger the bug</li>
        <li><strong>Expected behavior</strong> - What should happen</li>
        <li><strong>Actual behavior</strong> - What actually happens</li>
        <li><strong>Environment</strong> - OS, Go version, starknode-kit version</li>
        <li><strong>Logs</strong> - Relevant log output or error messages</li>
      </ol>

      <CodeBlock code={`# Example bug report template

**Description**
Brief description of the issue

**Steps to Reproduce**
1. Run \`starknode-kit start\`
2. ...
3. Error occurs

**Expected Behavior**
What should happen

**Actual Behavior**
What actually happens

**Environment**
- OS: Ubuntu 22.04
- Go: 1.24
- starknode-kit: v1.0.0

**Logs**
\`\`\`
Error log output here
\`\`\``} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Suggesting Features</h2>

      <p>When suggesting features, include:</p>

      <ul>
        <li>Clear description of the feature</li>
        <li>Use case and motivation</li>
        <li>Proposed implementation (if applicable)</li>
        <li>Potential impact on existing functionality</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Code of Conduct</h2>

      <p>We are committed to providing a welcoming and inclusive environment. Please:</p>

      <ul>
        <li>✅ Be respectful and considerate</li>
        <li>✅ Welcome newcomers and help them learn</li>
        <li>✅ Accept constructive criticism gracefully</li>
        <li>✅ Focus on what's best for the community</li>
        <li>❌ Don't harass or discriminate</li>
        <li>❌ Don't be disruptive or disrespectful</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Getting Help</h2>

      <p>If you need help with contributing:</p>

      <ul>
        <li>💬 <a href="https://t.me/+SCPbza9fk8dkYWI0" target="_blank" rel="noopener noreferrer">Join our Telegram</a></li>
        <li>🐛 <a href="https://github.com/thebuidl-grid/starknode-kit/issues" target="_blank" rel="noopener noreferrer">Open an issue</a></li>
        <li>💡 <a href="https://github.com/thebuidl-grid/starknode-kit/discussions" target="_blank" rel="noopener noreferrer">Start a discussion</a></li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Resources</h2>

      <ul>
        <li><a href="https://github.com/thebuidl-grid/starknode-kit" target="_blank" rel="noopener noreferrer">GitHub Repository</a></li>
        <li><a href="https://github.com/thebuidl-grid/starknode-kit/issues" target="_blank" rel="noopener noreferrer">Issues</a></li>
        <li><a href="https://github.com/thebuidl-grid/starknode-kit/pulls" target="_blank" rel="noopener noreferrer">Pull Requests</a></li>
        <li><a href="https://go.dev/doc/effective_go" target="_blank" rel="noopener noreferrer">Effective Go</a></li>
        <li><a href="https://www.conventionalcommits.org/" target="_blank" rel="noopener noreferrer">Conventional Commits</a></li>
      </ul>

      <div className="bg-green-50 border-l-4 border-green-500 p-4 my-6">
        <p className="font-semibold mb-2">🎉 Thank You!</p>
        <p className="mb-0">
          Thank you for considering contributing to starknode-kit! Your contributions help make this tool better for everyone.
        </p>
      </div>

    
    </div>
  );
}

