import Link from 'next/link';
import CodeBlock from '@/components/CodeBlock';

export default function Commands() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Commands Reference</h1>

      <p className="text-xl text-white mb-4 leading-relaxed">
        Complete reference for all starknode-kit commands. Each command helps
        you manage different aspects of your Ethereum and Starknet nodes.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Command Overview</h2>

      <div className="not-prose my-8 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Command</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Description</th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">add</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Add an Ethereum or Starknet client to the config</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">completion</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Generate the autocompletion script for the specified shell</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">config</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Create, show, and update your Starknet node configuration</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">monitor</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Launch real-time monitoring dashboard</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">remove</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Remove a specified resource</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">run</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Run a specific local infrastructure service</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">start</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Run the configured Ethereum clients</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">status</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Display status of running clients</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">stop</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Stop the configured Ethereum clients</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">update</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Check for and install client updates</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">validator</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Manage the Starknet validator client</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code text-gray-900 dark:text-gray-300">version</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Show version of starknode-kit or a specific client</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Quick Examples</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Add Clients</h3>
      <CodeBlock code="starknode-kit add --consensus_client lighthouse --execution_client geth" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Configure Network</h3>
      <CodeBlock code="starknode-kit config set network sepolia" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Start Ethereum Clients
      </h3>
      <CodeBlock code="starknode-kit start" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Monitor Nodes</h3>
      <CodeBlock code="starknode-kit monitor" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Check Status</h3>
      <CodeBlock code="starknode-kit status" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Run Individual Client
      </h3>
      <CodeBlock code="starknode-kit run juno" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Stop Ethereum Clients
      </h3>
      <CodeBlock code="starknode-kit stop" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Check Version</h3>
      <CodeBlock
        code={`# Check starknode-kit version
starknode-kit version

# Check specific client version
starknode-kit version geth`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Getting Help</h2>

      <p className="text-lg mb-6">
        For any command, you can use the <code>--help</code> flag to get
        detailed usage information:
      </p>

      <CodeBlock
        code={`starknode-kit --help
starknode-kit add --help
starknode-kit config --help`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Shell Completion</h2>

      <p className="text-lg mb-6">
        Generate autocompletion scripts for your shell:
      </p>

      <CodeBlock
        code={`# Bash
starknode-kit completion bash > /etc/bash_completion.d/starknode-kit

# Zsh
starknode-kit completion zsh > "\${fpath[1]}/_starknode-kit"

# Fish
starknode-kit completion fish > ~/.config/fish/completions/starknode-kit.fish`}
      />

      <div className="mt-12 p-6 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-900 dark:text-yellow-400">📖 Next Steps</h3>
        <p className="text-gray-700 dark:text-gray-300 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/clients" className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium">Supported Clients</Link>
        </div>
      </div>
    </div>
  );
}

