import Link from "next/link";
import CodeBlock from "@/components/CodeBlock";

export default function GettingStarted() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Getting Started</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        Welcome to starknode-kit! This guide will help you get up and running
        with your Ethereum and Starknet nodes in just a few minutes.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Prerequisites</h2>

      <p className="text-lg mb-6">Before you begin, make sure you have:</p>

      <ul className="space-y-3 mb-4">
        <li className="text-lg">
          <strong>Operating System:</strong> Linux or macOS (Windows via WSL)
        </li>
        <li className="text-lg">
          <strong>Go:</strong> Version 1.24 or later (if building from source)
        </li>
        <li className="text-lg">
          <strong>Storage:</strong> At least 2TB of free SSD space
        </li>
        <li className="text-lg">
          <strong>RAM:</strong> Minimum 32GB recommended
        </li>
        <li className="text-lg">
          <strong>Network:</strong> Stable internet connection
        </li>
      </ul>

      <div className="bg-blue-50 border-l-4 border-blue-500 p-6 my-10 rounded-r-lg">
        <p className="font-semibold mb-3 text-lg">📝 Note</p>
        <p className="mb-0 text-base">
          For detailed hardware requirements, check out our{" "}
          <Link href="/requirements" className="font-medium">
            Requirements page
          </Link>
          .
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Installation</h2>

      <p className="text-lg mb-6">
        The quickest way to install starknode-kit is using the installation
        script:
      </p>

      <CodeBlock code='/bin/bash -c "$(curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh)"' />

      <p className="text-lg mt-6 mb-4">This script will:</p>
      <ul className="space-y-2 mb-10">
        <li className="text-base">
          Download the latest version of starknode-kit
        </li>
        <li className="text-base">
          Install it to <code>/usr/local/bin/</code>
        </li>
        <li className="text-base">
          Create necessary configuration directories
        </li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Initial Configuration
      </h2>

      <p className="text-lg mb-6">
        After installation, generate your initial configuration file:
      </p>

      <CodeBlock code="starknode-kit config new" />

      <p className="text-base mt-6 mb-10">
        This creates a configuration file at{" "}
        <code>~/.starknode-kit/starknode.yml</code> with default settings.
      </p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Add Your First Clients
      </h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Ethereum Clients (Execution + Consensus)
      </h3>

      <p className="text-lg mb-6">
        To run an Ethereum node, you need both an execution client and a
        consensus client:
      </p>

      <CodeBlock code="starknode-kit add --consensus_client lighthouse --execution_client geth" />

      <p className="text-base mt-6 mb-6">Or with Reth and Prysm:</p>

      <CodeBlock code="starknode-kit add --consensus_client prysm --execution_client reth" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Starknet Client</h3>

      <p className="text-lg mb-6">To add a Starknet client (Juno):</p>

      <CodeBlock code="starknode-kit add --starknet_client juno" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Configure Network</h2>

      <p className="text-lg mb-6">
        By default, starknode-kit is configured for mainnet. To change to a
        test network:
      </p>

      <CodeBlock code="starknode-kit config set network sepolia" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Start Your Nodes</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Start Ethereum Clients
      </h3>

      <p className="text-lg mb-6">
        To start your configured Ethereum execution and consensus clients:
      </p>

      <CodeBlock code="starknode-kit start" />

      <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 my-10 rounded-r-lg">
        <p className="font-semibold mb-3 text-lg">⚠️ Important</p>
        <p className="mb-0 text-base">
          The <code>start</code> command only launches Ethereum clients
          (execution + consensus). It does not start Starknet clients.
        </p>
      </div>

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Run Individual Clients
      </h3>

      <p className="text-lg mb-6">To run a specific client:</p>

      <CodeBlock
        code={`# Run Juno (Starknet)
starknode-kit run juno

# Run Geth (Ethereum execution)
starknode-kit run geth

# Run Lighthouse (Ethereum consensus)
starknode-kit run lighthouse`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Monitor Your Nodes</h2>

      <p className="text-lg mb-6">
        Launch the real-time monitoring dashboard to see the status of your
        nodes:
      </p>

      <CodeBlock code="starknode-kit monitor" />

      <p className="text-lg mt-8 mb-6">The monitoring dashboard provides real-time insights:</p>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-10 not-prose">
        <div className="p-5 bg-blue-50 rounded-lg border border-blue-200">
          <div className="flex items-start gap-3">
            <span className="text-2xl">🔄</span>
            <div>
              <h4 className="font-semibold text-blue-900 mb-1">Node Sync Status</h4>
              <p className="text-sm text-gray-700">Real-time synchronization progress and health</p>
            </div>
          </div>
        </div>

        <div className="p-5 bg-green-50 rounded-lg border border-green-200">
          <div className="flex items-start gap-3">
            <span className="text-2xl">📊</span>
            <div>
              <h4 className="font-semibold text-green-900 mb-1">Current Block Height</h4>
              <p className="text-sm text-gray-700">Latest block number and sync progress</p>
            </div>
          </div>
        </div>

        <div className="p-5 bg-purple-50 rounded-lg border border-purple-200">
          <div className="flex items-start gap-3">
            <span className="text-2xl">🌐</span>
            <div>
              <h4 className="font-semibold text-purple-900 mb-1">Network Statistics</h4>
              <p className="text-sm text-gray-700">Peer connections and network performance</p>
            </div>
          </div>
        </div>

        <div className="p-5 bg-orange-50 rounded-lg border border-orange-200">
          <div className="flex items-start gap-3">
            <span className="text-2xl">💻</span>
            <div>
              <h4 className="font-semibold text-orange-900 mb-1">System Resources</h4>
              <p className="text-sm text-gray-700">CPU, RAM, and disk usage metrics</p>
            </div>
          </div>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Check Status</h2>

      <p className="text-lg mb-6">
        For a quick status check of all running clients:
      </p>

      <CodeBlock code="starknode-kit status" />

      <div className="mt-12 p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">📖 Next Steps</h3>
        <p className="text-gray-700 mb-4">
          Ready to dive deeper? Check out our installation guide:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/installation" className="text-blue-600 hover:text-blue-800 font-medium">Installation Guide</Link>
        </div>
      </div>
    </div>
  );
}

