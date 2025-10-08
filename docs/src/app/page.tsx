import Link from "next/link";
import CodeBlock from "@/components/CodeBlock";

export default function Home() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Welcome to starknode-kit</h1>

      <p className="text-xl text-gray-400 dark:text-gray-400 mb-8">
        A powerful command-line tool to help developers and node operators
        easily set up, manage, and maintain Ethereum and Starknet nodes.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 my-8 not-prose">
        <Link
          href="/getting-started"
          className="block p-6 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800 hover:border-blue-400 dark:hover:border-blue-600 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-blue-900 dark:text-blue-400">
            🚀 Getting Started
          </h3>
          <p className="text-gray-700 dark:text-gray-300">
            Learn how to install and configure starknode-kit for your node
            setup.
          </p>
        </Link>

        <Link
          href="/commands"
          className="block p-6 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800 hover:border-purple-400 dark:hover:border-purple-600 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-purple-900 dark:text-purple-400">
            📘 Commands
          </h3>
          <p className="text-gray-700 dark:text-gray-300">
            Explore all available commands and their usage.
          </p>
        </Link>

        <Link
          href="/configuration"
          className="block p-6 bg-green-50 dark:bg-green-900/20 rounded-lg border border-green-200 dark:border-green-800 hover:border-green-400 dark:hover:border-green-600 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-green-900 dark:text-green-400">
            ⚙️ Configuration
          </h3>
          <p className="text-gray-700 dark:text-gray-300">
            Configure your Ethereum and Starknet clients.
          </p>
        </Link>

        <Link
          href="/validator"
          className="block p-6 bg-orange-50 dark:bg-orange-900/20 rounded-lg border border-orange-200 dark:border-orange-800 hover:border-orange-400 dark:hover:border-orange-600 transition-colors no-underline"
        >
          <h3 className="text-xl font-semibold mb-2 text-orange-900 dark:text-orange-400">
            🔐 Validator Setup
          </h3>
          <p className="text-gray-700 dark:text-gray-300">
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
        <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold mb-2 text-gray-900 dark:text-white">Execution Layer</h4>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Geth</li>
            <li>• Reth</li>
          </ul>
        </div>
        <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold mb-2 text-gray-900 dark:text-white">Consensus Layer</h4>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Lighthouse</li>
            <li>• Prysm</li>
          </ul>
        </div>
        <div className="p-4 bg-gray-50 dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700">
          <h4 className="font-semibold mb-2 text-gray-900 dark:text-white">Starknet</h4>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Juno</li>
            <li>• Starknet Validator</li>
          </ul>
        </div>
      </div>

      <div className="mt-12 p-6 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-900 dark:text-yellow-400">📖 Next Steps</h3>
        <p className="text-gray-700 dark:text-gray-300 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
        <Link href="/installation" className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium">Installation Guide</Link>
        </div>
      </div>
    </div>
  );
}
