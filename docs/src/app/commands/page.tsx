import Link from 'next/link';
import CodeBlock from '@/components/CodeBlock';

export default function Commands() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Commands Reference</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        Complete reference for all starknode-kit commands. Each command helps
        you manage different aspects of your Ethereum and Starknet nodes.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Command Overview</h2>

      <div className="not-prose my-8">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Command</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">add</td>
              <td className="px-6 py-4 text-sm">Add an Ethereum or Starknet client to the config</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">completion</td>
              <td className="px-6 py-4 text-sm">Generate the autocompletion script for the specified shell</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">config</td>
              <td className="px-6 py-4 text-sm">Create, show, and update your Starknet node configuration</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">monitor</td>
              <td className="px-6 py-4 text-sm">Launch real-time monitoring dashboard</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">remove</td>
              <td className="px-6 py-4 text-sm">Remove a specified resource</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">run</td>
              <td className="px-6 py-4 text-sm">Run a specific local infrastructure service</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">start</td>
              <td className="px-6 py-4 text-sm">Run the configured Ethereum clients</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">status</td>
              <td className="px-6 py-4 text-sm">Display status of running clients</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">stop</td>
              <td className="px-6 py-4 text-sm">Stop the configured Ethereum clients</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">update</td>
              <td className="px-6 py-4 text-sm">Check for and install client updates</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">validator</td>
              <td className="px-6 py-4 text-sm">Manage the Starknet validator client</td>
            </tr>
            <tr>
              <td className="px-6 py-4 whitespace-nowrap text-sm font-mono">version</td>
              <td className="px-6 py-4 text-sm">Show version of starknode-kit or a specific client</td>
            </tr>
          </tbody>
        </table>
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

      <div className="mt-12 p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">📖 Next Steps</h3>
        <p className="text-gray-700 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/clients" className="text-blue-600 hover:text-blue-800 font-medium">Supported Clients</Link>
        </div>
      </div>
    </div>
  );
}

