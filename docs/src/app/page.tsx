import Link from "next/link";
import CodeBlock from "@/components/CodeBlock";

export default function Home() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Welcome to starknode-kit</h1>

      <p className="text-xl text-gray-600 mb-8">
        A powerful command-line tool to help developers and node operators
        easily set up, manage, and maintain Ethereum and Starknet nodes.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 my-8 not-prose">
        <Link
          href="/getting-started"
          className="block p-6 bg-blue-50 rounded-lg border border-blue-200 hover:border-blue-400 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-blue-900">
            🚀 Getting Started
          </h3>
          <p className="text-gray-700">
            Learn how to install and configure starknode-kit for your node
            setup.
          </p>
        </Link>

        <Link
          href="/commands"
          className="block p-6 bg-purple-50 rounded-lg border border-purple-200 hover:border-purple-400 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-purple-900">
            📘 Commands
          </h3>
          <p className="text-gray-700">
            Explore all available commands and their usage.
          </p>
        </Link>

        <Link
          href="/configuration"
          className="block p-6 bg-green-50 rounded-lg border border-green-200 hover:border-green-400 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-green-900">
            ⚙️ Configuration
          </h3>
          <p className="text-gray-700">
            Configure your Ethereum and Starknet clients.
          </p>
        </Link>

        <Link
          href="/validator"
          className="block p-6 bg-orange-50 rounded-lg border border-orange-200 hover:border-orange-400 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-orange-900">
            🔐 Validator Setup
          </h3>
          <p className="text-gray-700">
            Set up and manage your Starknet validator node.
          </p>
        </Link>
      </div>

      <h2 className="text-2xl font-bold mt-12 mb-4">Quick Start</h2>

      <p>Install starknode-kit with a single command:</p>

      <CodeBlock code='/bin/bash -c "$(curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh)"' />

      <p>Generate your configuration file:</p>

      <CodeBlock code="starknode-kit config new" />

      <p>Add your first client pair:</p>

      <CodeBlock code="starknode-kit add --consensus_client lighthouse --execution_client geth" />

      <h2 className="text-2xl font-bold mt-12 mb-4">Key Features</h2>

      <ul className="space-y-2">
        <li>
          ✅ <strong>Easy Setup</strong> - Get your node running in minutes
        </li>
        <li>
          ✅ <strong>Multi-Client Support</strong> - Works with Geth, Reth,
          Lighthouse, Prysm, and Juno
        </li>
        <li>
          ✅ <strong>Real-time Monitoring</strong> - Built-in dashboard to
          monitor your nodes
        </li>
        <li>
          ✅ <strong>Auto Updates</strong> - Keep your clients up to date
          automatically
        </li>
        <li>
          ✅ <strong>Validator Management</strong> - Simplified Starknet
          validator operations
        </li>
        <li>
          ✅ <strong>Network Flexibility</strong> - Support for mainnet,
          sepolia, and custom networks
        </li>
      </ul>

      <h2 className="text-2xl font-bold mt-12 mb-4">Supported Clients</h2>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 not-prose">
        <div className="p-4 bg-gray-50 rounded-lg">
          <h4 className="font-semibold mb-2">Execution Layer</h4>
          <ul className="text-sm space-y-1 text-gray-700">
            <li>• Geth</li>
            <li>• Reth</li>
          </ul>
        </div>
        <div className="p-4 bg-gray-50 rounded-lg">
          <h4 className="font-semibold mb-2">Consensus Layer</h4>
          <ul className="text-sm space-y-1 text-gray-700">
            <li>• Lighthouse</li>
            <li>• Prysm</li>
          </ul>
        </div>
        <div className="p-4 bg-gray-50 rounded-lg">
          <h4 className="font-semibold mb-2">Starknet</h4>
          <ul className="text-sm space-y-1 text-gray-700">
            <li>• Juno</li>
            <li>• Starknet Validator</li>
          </ul>
        </div>
      </div>

      <div className="mt-12 p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">📖 Next Steps</h3>
        <p className="text-gray-700 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
        <Link href="/installation" className="text-white">Installation Guide</Link>
        </div>
      </div>
    </div>
  );
}
