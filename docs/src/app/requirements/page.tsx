import Link from 'next/link';
import CodeBlock from '@/components/CodeBlock';

export default function Requirements() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Requirements</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        Hardware and software requirements for running Ethereum and Starknet
        nodes with starknode-kit.
      </p>

      <h2 className="text-3xl font-semibold mb-6">Hardware Requirements</h2>

      <div className="bg-blue-50 border-l-4 border-blue-500 p-4 my-6">
        <p className="font-semibold mb-2">📚 Reference</p>
        <p className="mb-0">
          For a detailed breakdown of node hardware requirements, see the <a href="https://docs.rocketpool.net/guides/node/hardware.html" target="_blank" rel="noopener noreferrer">Rocket Pool Hardware Guide</a>.
        </p>
      </div>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Minimum Requirements</h3>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Component</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Requirement</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Notes</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">CPU</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">4+ cores</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Intel i3/i5 or AMD equivalent. Avoid Celeron.</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">RAM</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">32 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Minimum 16GB, 32GB recommended for comfort</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Storage</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">2+ TB NVMe SSD</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Must have DRAM cache, no QLC NAND</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Network</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">100+ Mbps</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Stable connection, unlimited data preferred</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Power</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">24/7 uptime</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">UPS recommended for validators</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Recommended Specifications</h3>

      <ul>
        <li><strong>CPU:</strong> Intel i5/i7 or AMD Ryzen 5/7 (6+ cores)</li>
        <li><strong>RAM:</strong> 64 GB DDR4</li>
        <li><strong>Storage:</strong> 4 TB NVMe SSD with DRAM cache</li>
        <li><strong>Network:</strong> 1 Gbps fiber connection</li>
        <li><strong>Backup Power:</strong> UPS with 30+ minutes runtime</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Storage Requirements</h2>

      <p>Storage is the most critical component for node operation.</p>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Storage Size</h3>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Client</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Current Size</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Growth Rate</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Ethereum (Geth)</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~1.2 TB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~150 GB/year</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Ethereum (Reth)</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~900 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~120 GB/year</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Lighthouse</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~200 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~50 GB/year</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Prysm</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~250 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~60 GB/year</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Juno (Starknet)</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~300 GB</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~100 GB/year</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h3 className="text-2xl font-semibold mt-10 mb-5">SSD Requirements</h3>

      <p>Your SSD <strong>must have</strong>:</p>

      <ul>
        <li>✅ <strong>DRAM cache</strong> - Essential for performance</li>
        <li>✅ <strong>TLC or better NAND</strong> - No QLC (Quad-Level Cell)</li>
        <li>✅ <strong>High endurance rating</strong> - 600+ TBW recommended</li>
        <li>✅ <strong>NVMe interface</strong> - SATA SSDs are too slow</li>
      </ul>

      <div className="bg-yellow-50 border-l-4 border-yellow-500 p-4 my-6">
        <p className="font-semibold mb-2">⚠️ Warning</p>
        <p className="mb-0">
          Using a QLC SSD or SSD without DRAM cache will result in poor performance and potential node failures. 
          See the <a href="https://gist.github.com/bkase/fab02c5b3c404e9ef8e5c2071ac1558c" target="_blank" rel="noopener noreferrer">tested SSD list</a> for recommendations.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Software Requirements</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Operating System</h3>

      <p>Supported operating systems:</p>

      <ul>
        <li><strong>Linux:</strong> Ubuntu 20.04+, Debian 11+, or other modern distributions</li>
        <li><strong>macOS:</strong> macOS 12 (Monterey) or later</li>
        <li><strong>Windows:</strong> Windows 10/11 with WSL2 (Ubuntu)</li>
      </ul>

      <p className="text-lg font-semibold">Linux is highly recommended for production use.</p>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Required Software</h3>

      <h4 className="text-xl font-semibold mt-6 mb-3">Go (for building from source)</h4>

      <p>Version 1.24 or later required:</p>

      <CodeBlock code={`# Check Go version
go version

# Should output: go version go1.24 or higher`} />

      <p>Install from: <a href="https://go.dev/dl/" target="_blank" rel="noopener noreferrer">https://go.dev/dl/</a></p>

      <h4 className="text-xl font-semibold mt-6 mb-3">Rust (for Starknet clients)</h4>

      <p>Recommended for building Juno and other Starknet clients:</p>

      <CodeBlock code="curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh" />

      <h4 className="text-xl font-semibold mt-6 mb-3">Make</h4>

      <p>Required for building certain clients:</p>

      <CodeBlock code={`# Ubuntu/Debian
sudo apt install make

# macOS (with Homebrew)
brew install make

# Check installation
make --version`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Network Requirements</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Bandwidth</h3>

      <ul>
        <li><strong>Download:</strong> 100+ Mbps</li>
        <li><strong>Upload:</strong> 25+ Mbps</li>
        <li><strong>Data Cap:</strong> Unlimited (or 2+ TB/month)</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Ports</h3>

      <p>Ensure these ports are accessible:</p>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Port</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Protocol</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Purpose</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code">30303</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">TCP/UDP</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Ethereum execution P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code">9000</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">TCP/UDP</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Lighthouse consensus P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code">13000</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">TCP</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Prysm consensus P2P</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-source-code">6060</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">TCP</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Juno RPC (localhost only)</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">For Validator Nodes</h2>

      <p>Additional requirements for running a validator:</p>

      <ul>
        <li><strong>Uptime:</strong> 99.9%+ availability required</li>
        <li><strong>Backup Power:</strong> UPS mandatory</li>
        <li><strong>Monitoring:</strong> 24/7 monitoring and alerting</li>
        <li><strong>Backup Internet:</strong> Secondary connection recommended</li>
        <li><strong>Dedicated Hardware:</strong> No shared resources</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Tested Hardware Configurations</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Budget Build (~$800)</h3>

      <ul>
        <li>Intel NUC 13 PRO (i5)</li>
        <li>32 GB DDR4 RAM</li>
        <li>2 TB NVMe SSD (Samsung 980 PRO)</li>
        <li>Ubuntu 22.04 LTS</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Recommended Build (~$1500)</h3>

      <ul>
        <li>Custom build: AMD Ryzen 7 or Intel i7</li>
        <li>64 GB DDR4 RAM</li>
        <li>4 TB NVMe SSD (Samsung 990 PRO)</li>
        <li>1 Gbps fiber connection</li>
        <li>UPS with 30min+ runtime</li>
        <li>Ubuntu 22.04 LTS</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Pro Build (~$3000+)</h3>

      <ul>
        <li>High-end workstation or server</li>
        <li>128 GB ECC RAM</li>
        <li>8 TB NVMe SSD (enterprise grade)</li>
        <li>Redundant power supplies</li>
        <li>Redundant network connections</li>
        <li>Professional monitoring and alerting</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Cloud Providers</h2>

      <p>If running in the cloud, recommended specifications:</p>

      <div className="not-prose my-6 overflow-x-auto -mx-4 sm:mx-0">
        <div className="inline-block min-w-full align-middle">
          <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 sm:rounded-lg">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Provider</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Instance Type</th>
                  <th className="px-4 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Est. Cost/Month</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">AWS</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">m5.2xlarge + 4TB gp3</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~$500-700</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Google Cloud</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">n2-standard-8 + 4TB SSD</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~$600-800</td>
                </tr>
                <tr>
                  <td className="px-4 sm:px-6 py-4 whitespace-nowrap text-sm font-medium">Azure</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">Standard_D8s_v3 + 4TB Premium SSD</td>
                  <td className="px-4 sm:px-6 py-4 text-sm">~$550-750</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>

      <div className="bg-blue-50 border-l-4 border-blue-500 p-4 my-6">
        <p className="font-semibold mb-2">💡 Cost Consideration</p>
        <p className="mb-0">
          Running on dedicated hardware is often more cost-effective long-term than cloud hosting, especially for validators.
        </p>
      </div>
    </div>
  );
}

