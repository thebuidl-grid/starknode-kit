import Link from 'next/link';
import CodeBlock from '@/components/CodeBlock';

export default function Installation() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Installation</h1>

      <p className="text-xl text-white mb-4 leading-relaxed">
        There are multiple ways to install starknode-kit. Choose the method that
        best suits your needs.
      </p>

      <h2 className="text-3xl font-semibold mb-6">
        Option 1: Install Script (Recommended)
      </h2>

      <p className="text-lg mb-6">
        The easiest and recommended way to install starknode-kit:
      </p>

      <CodeBlock code='/bin/bash -c "$(curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh)"' />

      <p className="text-base mt-6 mb-6">
        Or download the script first and then run it:
      </p>

      <CodeBlock
        code={`wget https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh
chmod +x install.sh
./install.sh`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Option 2: Install using Go
      </h2>

      <p className="text-lg mb-6">
        If you have Go installed (version 1.24 or later), you can install
        starknode-kit directly:
      </p>

      <CodeBlock
        code={`go install -ldflags="-X 'github.com/thebuidl-grid/starknode-kit/pkg/versions.StarkNodeVersion=main'" github.com/thebuidl-grid/starknode-kit@latest`}
      />

      <p className="text-base mt-6 mb-10">
        This installs the latest version from the <code>main</code> branch.
      </p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Option 3: Manual Installation from Source
      </h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        1. Clone the Repository
      </h3>

      <CodeBlock
        code={`git clone https://github.com/thebuidl-grid/starknode-kit.git
cd starknode-kit`}
      />

      <h3 className="text-2xl font-semibold mt-10 mb-5">2. Build and Install</h3>

      <CodeBlock
        code={`make build
sudo mv bin/starknode /usr/local/bin/`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Verify Installation</h2>

      <p className="text-lg mb-6">
        After installation, verify that starknode-kit is working correctly:
      </p>

      <CodeBlock code="starknode-kit --help" />

      <p className="text-base mt-6 mb-6">You should see output similar to:</p>

      <CodeBlock
        code={`starknode-kit is a CLI tool designed to simplify the setup and management 
of Ethereum and Starknet nodes. It helps developers quickly configure, 
launch, monitor, and maintain full nodes or validator setups for both networks.

Usage:
  starknode [command]

Available Commands:
  add         Add an Ethereum or Starknet client to the config
  completion  Generate the autocompletion script for the specified shell
  config      Create, show, and update your Starknet node configuration.
  help        Help about any command
  monitor     Launch real-time monitoring dashboard
  remove      Remove a specified resource
  run         Run a specific local infrastructure service
  start       Run the configured Ethereum clients
  status      Display status of running clients
  stop        Stop the configured Ethereum clients
  update      Check for and install client updates
  validator   Manage the Starknet validator client
  version     Show version of starknode-kit or a specific client

Flags:
  -h, --help   help for starknode

Use "starknode [command] --help" for more information about a command.`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Initial Setup</h2>

      <p className="text-lg mb-6">
        After successful installation, generate your configuration file:
      </p>

      <CodeBlock code="starknode-kit config new" />

      <p className="text-base mt-6 mb-10">
        This creates a configuration file at{" "}
        <code>~/.starknode-kit/starknode.yml</code>.
      </p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Uninstallation</h2>

      <p className="text-lg mb-6">
        To uninstall starknode-kit, remove the binary and configuration
        directory:
      </p>

      <CodeBlock
        code={`sudo rm /usr/local/bin/starknode-kit
rm -rf ~/.config/starknode-kit`}
      />

      <div className="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-500 dark:border-yellow-600 p-6 my-10 rounded-r-lg">
        <p className="font-semibold mb-3 text-lg text-gray-900 dark:text-yellow-400">⚠️ Note</p>
        <p className="mb-0 text-base text-gray-700 dark:text-gray-300">
          This will not remove any client data (e.g., blockchain data). The data
          is stored in the locations specified in your{" "}
          <code>~/.starknode-kit/starknode.yml</code> file.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Troubleshooting</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Command not found</h3>

      <p className="text-lg mb-6">
        If you get a "command not found" error, make sure{" "}
        <code>/usr/local/bin</code> is in your PATH:
      </p>

      <CodeBlock code="export PATH=$PATH:/usr/local/bin" />

      <p className="text-base mt-6 mb-10">
        Add this to your <code>~/.bashrc</code> or <code>~/.zshrc</code> to make
        it permanent.
      </p>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Permission denied</h3>

      <p className="text-lg mb-6">
        If you encounter permission issues during installation, make sure you
        have sudo access or contact your system administrator.
      </p>

      <div className="mt-12 p-6 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-900 dark:text-yellow-400">📖 Next Steps</h3>
        <p className="text-gray-700 dark:text-gray-300 mb-4">
          Ready to dive deeper? Check out our configuration guide:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/configuration" className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium">Configuration Guide</Link>
        </div>
      </div>
    </div>
  );
}

