import Link from 'next/link';

export default function Clients() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Supported Clients</h1>

      <p className="text-xl text-white mb-4 leading-relaxed">
        starknode-kit supports multiple client implementations for both Ethereum
        and Starknet networks.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Why Multiple Clients?</h2>

      <p className="text-lg mb-10">
        Running diverse client implementations is crucial for network health and
        resilience. Client diversity prevents single points of failure and
        reduces the impact of bugs in any one implementation.
      </p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Ethereum Clients</h2>

      <p className="text-lg mb-8">
        To run an Ethereum node, you need both an{" "}
        <strong>execution client</strong> and a <strong>consensus client</strong>
        . They work together to validate and process Ethereum blocks post-merge.
      </p>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 not-prose my-8">
        <div className="p-6 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
          <h3 className="text-xl font-semibold mb-2 text-blue-900 dark:text-blue-400">Execution Clients</h3>
          <p className="text-gray-700 dark:text-gray-300 mb-3">Handle transaction execution and state management</p>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Geth (Go)</li>
            <li>• Reth (Rust)</li>
          </ul>
        </div>

        <div className="p-6 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800">
          <h3 className="text-xl font-semibold mb-2 text-purple-900 dark:text-purple-400">Consensus Clients</h3>
          <p className="text-gray-700 dark:text-gray-300 mb-3">Handle proof-of-stake consensus mechanism</p>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Lighthouse (Rust)</li>
            <li>• Prysm (Go)</li>
          </ul>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Starknet Clients</h2>

      <p>
        Starknet clients allow you to run a Starknet full node, enabling interaction with the Starknet Layer 2 network.
      </p>

      <div className="not-prose my-8">
        <div className="p-6 bg-orange-50 dark:bg-orange-900/20 rounded-lg border border-orange-200 dark:border-orange-800 max-w-2xl">
          <h3 className="text-xl font-semibold mb-2 text-orange-900 dark:text-orange-400">Starknet Clients</h3>
          <p className="text-gray-700 dark:text-gray-300 mb-3">Full node implementations for Starknet</p>
          <ul className="text-sm space-y-1 text-gray-700 dark:text-gray-300">
            <li>• Juno (Go) - Full node client</li>
            <li>• Starknet Validator - Validator client for staking</li>
          </ul>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Client Combinations</h2>

      <p>Popular client combinations for Ethereum nodes:</p>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Execution</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Consensus</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Characteristics</th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Geth</td>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">Lighthouse</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Most popular, well-tested</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Reth</td>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">Lighthouse</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">High performance, modern</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Geth</td>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">Prysm</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Stable, feature-rich</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Reth</td>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">Prysm</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">Performance-focused</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Choosing Clients</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Execution Clients</h3>

      <h4 className="text-xl font-semibold mt-6 mb-3">Geth</h4>
      <ul className="space-y-1 mb-8">
        <li className="text-base">✅ Most widely used and tested</li>
        <li className="text-base">✅ Excellent documentation</li>
        <li className="text-base">✅ Large community support</li>
        <li className="text-base">✅ Stable and reliable</li>
        <li className="text-base">⚠️ Higher resource usage</li>
        <li className="text-base">⚠️ Larger disk footprint</li>
      </ul>

      <h4 className="text-xl font-semibold mt-6 mb-3">Reth</h4>
      <ul className="space-y-1 mb-8">
        <li className="text-base">✅ Excellent performance</li>
        <li className="text-base">✅ Lower disk usage</li>
        <li className="text-base">✅ Modern codebase (Rust)</li>
        <li className="text-base">✅ Fast sync times</li>
        <li className="text-base">⚠️ Newer, less battle-tested</li>
        <li className="text-base">⚠️ Smaller community</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Consensus Clients</h3>

      <h4 className="text-xl font-semibold mt-6 mb-3">Lighthouse</h4>
      <ul className="space-y-1 mb-8">
        <li className="text-base">✅ Fast and efficient</li>
        <li className="text-base">✅ Low resource usage</li>
        <li className="text-base">✅ Great documentation</li>
        <li className="text-base">✅ Active development</li>
        <li className="text-base">✅ Written in Rust</li>
      </ul>

      <h4 className="text-xl font-semibold mt-6 mb-3">Prysm</h4>
      <ul className="space-y-1 mb-8">
        <li className="text-base">✅ Feature-rich</li>
        <li className="text-base">✅ Good performance</li>
        <li className="text-base">✅ Strong community</li>
        <li className="text-base">✅ Comprehensive tooling</li>
        <li className="text-base">✅ Written in Go</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Starknet Clients</h3>

      <h4 className="text-xl font-semibold mt-6 mb-3">Juno</h4>
      <ul className="space-y-1 mb-8">
        <li className="text-base">✅ Official full node client</li>
        <li className="text-base">✅ Well-maintained</li>
        <li className="text-base">✅ Fast sync</li>
        <li className="text-base">✅ Active community</li>
        <li className="text-base">✅ Required for Starknet validator</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Resource Requirements by Client</h2>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50 dark:bg-gray-800">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Client</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">RAM</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Disk</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">CPU</th>
                </tr>
              </thead>
              <tbody className="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Geth</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">16+ GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">~1.2 TB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">4+ cores</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Reth</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">16+ GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">~900 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">4+ cores</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Lighthouse</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">8+ GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">~200 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">2+ cores</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Prysm</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">8+ GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">~250 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">2+ cores</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 dark:text-white">Juno</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">8+ GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">~300 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm text-gray-900 dark:text-gray-300">2+ cores</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 dark:border-blue-600 p-4 my-6">
        <p className="font-semibold mb-2 text-gray-900 dark:text-blue-400">💡 Recommendation</p>
        <p className="mb-0 text-gray-700 dark:text-gray-300">
          For most users, we recommend <strong>Reth + Lighthouse</strong> for Ethereum (best performance) 
          and <strong>Juno</strong> for Starknet.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Client Diversity</h2>

      <p>
        Client diversity is critical for network health. If a single client has a bug and it's used by the majority of nodes, 
        it could cause network issues or even finality problems.
      </p>

      <p><strong>Current client distribution matters!</strong> Consider using minority clients to help decentralize the network.</p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Switching Clients</h2>

      <p>You can switch clients at any time:</p>

      <ol>
        <li>Stop your current clients: <code>starknode-kit stop</code></li>
        <li>Remove old client: <code>starknode-kit remove --execution_client geth</code></li>
        <li>Add new client: <code>starknode-kit add --execution_client reth</code></li>
        <li>Start nodes: <code>starknode-kit start</code></li>
      </ol>

      <div className="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-500 dark:border-yellow-600 p-4 my-6">
        <p className="font-semibold mb-2 text-gray-900 dark:text-yellow-400">⚠️ Note</p>
        <p className="mb-0 text-gray-700 dark:text-gray-300">
          Switching clients may require re-syncing from scratch, which can take several days. 
          Plan accordingly and ensure you have sufficient disk space.
        </p>
      </div>

      <div className="mt-12 p-6 bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg">
        <h3 className="text-lg font-semibold mb-2 text-gray-900 dark:text-yellow-400">📖 Next Steps</h3>
        <p className="text-gray-700 dark:text-gray-300 mb-4">
          Ready to dive deeper? Check out our validator guide:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/validator" className="text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 font-medium">Validator Guide</Link>
        </div>
      </div>
    </div>
  );
}

