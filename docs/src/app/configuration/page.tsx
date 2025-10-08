import CodeBlock from '@/components/CodeBlock';
import Link from 'next/link';

export default function Configuration() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Configuration</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        Learn how to configure starknode-kit for your Ethereum and Starknet
        nodes.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Configuration File</h2>

      <p className="text-lg mb-6">
        starknode-kit stores its configuration in a YAML file located at:
      </p>

      <CodeBlock code="~/.starknode-kit/starknode.yml" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Creating a Configuration
      </h2>

      <p className="text-lg mb-6">
        Generate a new configuration file with default settings:
      </p>

      <CodeBlock code="starknode-kit config new" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Viewing Configuration
      </h2>

      <p className="text-lg mb-6">View your entire configuration:</p>

      <CodeBlock code="starknode-kit config show --all" />

      <p className="text-base mt-6 mb-6">View specific sections:</p>

      <CodeBlock
        code={`# View execution client config
starknode-kit config show --el

# View consensus client config
starknode-kit config show --cl

# View Juno (Starknet) config
starknode-kit config show --juno`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Configuration Structure
      </h2>

      <p className="text-lg mb-6">
        The configuration file has the following structure:
      </p>

      <CodeBlock code={`network: mainnet

execution_client:
  name: geth
  ports:
    - 8545  # HTTP RPC
    - 8546  # WebSocket RPC
    - 30303 # P2P

consensus_client:
  name: lighthouse
  ports:
    - 5052  # HTTP API
    - 9000  # P2P
  consensus_checkpoint: ""

juno_client:
  port: 6060
  eth_node: "http://localhost:8545"
  environment: []

is_validator_node: false

wallet:
  name: ""
  reward_address: ""
  commision: ""

validator_config:
  provider_config:
    juno_rpc_http: "http://localhost:6060"
    juno_rpc_ws: "ws://localhost:6060"
  signer:
    operational_address: ""
    privateKey: ""`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Modifying Configuration
      </h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Change Network</h3>

      <p className="text-lg mb-6">
        Switch between mainnet, sepolia, or custom networks:
      </p>

      <CodeBlock
        code={`# Switch to sepolia testnet
starknode-kit config set network sepolia

# Switch to mainnet
starknode-kit config set network mainnet`}
      />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Configure Execution Client
      </h3>

      <p className="text-lg mb-6">Set execution client and ports:</p>

      <CodeBlock
        code={`# Set client type
starknode-kit config set el client=reth

# Set ports
starknode-kit config set el port=9000,9001

# Set both
starknode-kit config set el client=geth port=8545,8546,30303`}
      />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Configure Consensus Client
      </h3>

      <p className="text-lg mb-6">Set consensus client and checkpoint:</p>

      <CodeBlock
        code={`# Set client type
starknode-kit config set cl client=prysm

# Set checkpoint URL
starknode-kit config set cl checkpoint=https://checkpoint.example.com`}
      />

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Configure Juno (Starknet)
      </h3>

      <p className="text-lg mb-6">Configure your Juno Starknet client:</p>

      <CodeBlock
        code={`# Set Juno port
starknode-kit config set juno port=6060

# Set Ethereum node connection
starknode-kit config set juno eth_node=http://localhost:8545`}
      />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Network Settings</h2>

      <p className="text-lg mb-6">starknode-kit supports multiple networks:</p>

      <ul className="space-y-2 mb-10">
        <li className="text-base">
          <strong>mainnet</strong> - Ethereum and Starknet mainnet
        </li>
        <li className="text-base">
          <strong>sepolia</strong> - Ethereum Sepolia and Starknet Sepolia
          testnet
        </li>
      </ul>

      <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 my-10 rounded-r-lg">
        <p className="font-semibold mb-3 text-lg">⚠️ Important</p>
        <p className="mb-0 text-base">
          Changing the network will affect all clients. Make sure to stop your
          nodes before changing networks.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Port Configuration
      </h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Default Ports</h3>

      <p className="text-lg mb-6">Default ports for each client:</p>

      <div className="not-prose my-8 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Client</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Ports</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Purpose</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Geth</td>
                  <td className="px-4 sm:px-6 py-4 text-sm font-source-code">8545, 8546, 30303</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">HTTP RPC, WS RPC, P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Reth</td>
                  <td className="px-4 sm:px-6 py-4 text-sm font-source-code">8545, 8546, 30303</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">HTTP RPC, WS RPC, P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Lighthouse</td>
                  <td className="px-4 sm:px-6 py-4 text-sm font-source-code">5052, 9000</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">HTTP API, P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Prysm</td>
                  <td className="px-4 sm:px-6 py-4 text-sm font-source-code">4000, 13000</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">HTTP API, P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Juno</td>
                  <td className="px-4 sm:px-6 py-4 text-sm font-source-code">6060</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">RPC</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Validator Configuration
      </h2>

      <p className="text-lg mb-10">
        For validator nodes, additional configuration is required. See the
        Validator Setup page for details.
      </p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Environment Variables
      </h2>

      <p className="text-lg mb-6">
        Some sensitive data can be stored as environment variables:
      </p>

      <ul className="space-y-2 mb-10">
        <li className="text-base">
          <code>STARKNET_WALLET</code> - Wallet address
        </li>
        <li className="text-base">
          <code>STARKNET_PRIVATE_KEY</code> - Private key
        </li>
        <li className="text-base">
          <code>STARKNET_PUBLIC_KEY</code> - Public key
        </li>
        <li className="text-base">
          <code>STARKNET_CLASS_HASH</code> - Class hash
        </li>
        <li className="text-base">
          <code>STARKNET_SALT</code> - Salt value
        </li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">
        Configuration Best Practices
      </h2>

      <ol className="space-y-2 mb-10">
        <li className="text-base">
          <strong>Backup your config</strong> - Keep a backup of your
          configuration file
        </li>
        <li className="text-base">
          <strong>Use environment variables</strong> - Store sensitive data in
          environment variables
        </li>
        <li className="text-base">
          <strong>Document changes</strong> - Keep notes of any custom
          configurations
        </li>
        <li className="text-base">
          <strong>Test on testnet</strong> - Always test configuration changes
          on a testnet first
        </li>
      </ol>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Troubleshooting</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">
        Configuration not loading
      </h3>

      <p className="text-lg mb-6">
        If your configuration isn't loading, check:
      </p>

      <ul className="space-y-2 mb-10">
        <li className="text-base">
          File exists at <code>~/.starknode-kit/starknode.yml</code>
        </li>
        <li className="text-base">File has correct YAML syntax</li>
        <li className="text-base">
          File has correct permissions (readable by your user)
        </li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Port conflicts</h3>

      <p className="text-lg mb-6">If you get port conflicts:</p>

      <ul className="space-y-2 mb-10">
        <li className="text-base">
          Check if ports are already in use: <code>lsof -i :[port]</code>
        </li>
        <li className="text-base">Configure different ports in your config</li>
        <li className="text-base">Stop conflicting services</li>
      </ul>

      <div className="mt-12 p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">📖 Next Steps</h3>
        <p className="text-gray-700 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/commands" className="text-blue-600 hover:text-blue-800 font-medium">Commands Reference</Link>
        </div>
      </div>
    </div>
  );
}

