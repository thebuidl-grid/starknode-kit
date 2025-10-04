import CodeBlock from '@/components/CodeBlock';
import Link from 'next/link';
export default function Validator() {
  return (
    <div className="prose prose-lg max-w-none">
      <h1 className="text-4xl font-bold mb-4">Validator Setup</h1>

      <p className="text-xl text-gray-600 mb-4 leading-relaxed">
        Set up and manage your Starknet validator node using starknode-kit.
      </p>

      <div className="bg-yellow-50 border-l-4 border-yellow-500 p-6 my-10 rounded-r-lg">
        <p className="font-semibold mb-3 text-lg">⚠️ Important</p>
        <p className="mb-0 text-base">
          Running a validator requires significant responsibility. Make sure you
          understand the requirements and risks before proceeding.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mb-6">Prerequisites</h2>

      <p className="text-lg mb-6">Before setting up a validator, ensure you have:</p>

      <ul className="space-y-3 mb-10">
        <li className="text-base">✅ A fully synced Juno (Starknet) node</li>
        <li className="text-base">✅ A Starknet wallet with sufficient funds for staking</li>
        <li className="text-base">✅ Stable internet connection with 99.9%+ uptime</li>
        <li className="text-base">✅ Understanding of validator responsibilities</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Installation</h2>

      <p>The validator client is managed through starknode-kit. First, ensure you have starknode-kit installed and configured.</p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Validator Commands</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Check Validator Status</h3>

      <p>Check the status of your validator:</p>

      <CodeBlock code="starknode-kit validator status" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Get Validator Version</h3>

      <p>Check the installed version of the validator client:</p>

      <CodeBlock code="starknode-kit validator --version" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Set Juno RPC Endpoint</h3>

      <p>Configure the Juno RPC endpoint for your validator:</p>

      <CodeBlock code="starknode-kit validator --rpc http://localhost:6060" />

      <p>Or use a remote Juno node:</p>

      <CodeBlock code="starknode-kit validator --rpc https://your-juno-node.example.com" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Configuration</h2>

      <p>Validator configuration is stored in your starknode-kit config file. Key settings include:</p>

      <CodeBlock code={`is_validator_node: true

wallet:
  name: "my-validator"
  reward_address: "0x..."
  commision: "10"  # 10% commission
  wallet:
    address: "\${STARKNET_WALLET}"
    private_key: "\${STARKNET_PRIVATE_KEY}"
    public_key: "\${STARKNET_PUBLIC_KEY}"
    class_hash: "\${STARKNET_CLASS_HASH}"
    salt: "\${STARKNET_SALT}"
    deployed: false
    legacy: false

validator_config:
  provider_config:
    juno_rpc_http: "http://localhost:6060"
    juno_rpc_ws: "ws://localhost:6060"
  signer:
    operational_address: "0x..."
    privateKey: "\${STARKNET_PRIVATE_KEY}"`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Setting Up Environment Variables</h2>

      <p>Store sensitive validator data in environment variables:</p>

      <CodeBlock code={`export STARKNET_WALLET="0x..."
export STARKNET_PRIVATE_KEY="0x..."
export STARKNET_PUBLIC_KEY="0x..."
export STARKNET_CLASS_HASH="0x..."
export STARKNET_SALT="0x..."`} />

      <p>Add these to your <code>~/.bashrc</code> or <code>~/.zshrc</code> to make them persistent.</p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Starting Your Validator</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Step 1: Ensure Juno is Running</h3>

      <p>Your Juno node must be fully synced and running:</p>

      <CodeBlock code="starknode-kit run juno" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Step 2: Verify Configuration</h3>

      <p>Check your validator configuration:</p>

      <CodeBlock code="starknode-kit config show --all" />

      <h3 className="text-2xl font-semibold mt-10 mb-5">Step 3: Start the Validator</h3>

      <p>Start your validator client:</p>

      <CodeBlock code="starknode-kit run starknet-staking-v2" />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Monitoring Your Validator</h2>

      <p>Monitor your validator status:</p>

      <CodeBlock code={`# Check validator status
starknode-kit validator status

# Monitor all services
starknode-kit monitor`} />

      <h2 className="text-3xl font-semibold mt-16 mb-6">Validator Responsibilities</h2>

      <ol>
        <li><strong>Uptime</strong> - Maintain high availability (99.9%+)</li>
        <li><strong>Security</strong> - Keep your keys secure and never share them</li>
        <li><strong>Updates</strong> - Keep your validator software up to date</li>
        <li><strong>Monitoring</strong> - Actively monitor your validator performance</li>
        <li><strong>Backup</strong> - Maintain secure backups of your keys</li>
      </ol>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Security Best Practices</h2>

      <ul>
        <li>✅ Use hardware wallets when possible</li>
        <li>✅ Store keys in environment variables, not in config files</li>
        <li>✅ Use firewall to restrict access to validator ports</li>
        <li>✅ Enable SSH key-based authentication</li>
        <li>✅ Keep your server updated with security patches</li>
        <li>✅ Monitor for unusual activity</li>
        <li>✅ Have a disaster recovery plan</li>
        <li>❌ Never share your private keys</li>
        <li>❌ Don't run validators on shared hosting</li>
        <li>❌ Avoid storing keys in version control</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Staking and Commission</h2>

      <p>Configure your validator's commission rate:</p>

      <CodeBlock code='starknode-kit config set wallet commision="10"' />

      <p>This sets a 10% commission on staking rewards. Validators typically charge 5-15% commission.</p>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Troubleshooting</h2>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Validator Not Connecting to Juno</h3>

      <p>Check that:</p>
      <ul>
        <li>Juno is running: <code>starknode-kit status</code></li>
        <li>RPC endpoint is correct in config</li>
        <li>Firewall allows connections to Juno port</li>
      </ul>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Keys Not Loading</h3>

      <p>Verify environment variables are set:</p>

      <CodeBlock code="echo $STARKNET_PRIVATE_KEY" />

      <p>If empty, add to your shell profile and reload.</p>

      <h3 className="text-2xl font-semibold mt-10 mb-5">Validator Offline</h3>

      <p className="text-lg mb-6">If your validator goes offline:</p>
      <ol className="space-y-2 mb-10">
        <li className="text-base">✅ Check system resources</li>
        <li className="text-base">✅ Check network connectivity</li>
        <li className="text-base">✅ Review validator logs</li>
        <li className="text-base">✅ Restart validator if needed</li>
      </ol>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Performance Metrics</h2>

      <p>Monitor these key metrics:</p>

      <ul>
        <li><strong>Attestation rate</strong> - Percentage of successful attestations</li>
        <li><strong>Block proposals</strong> - Number of blocks proposed</li>
        <li><strong>Uptime</strong> - Validator availability percentage</li>
        <li><strong>Rewards</strong> - Staking rewards earned</li>
      </ul>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Validator Economics</h2>

      <p>Understand the economics:</p>

      <ul>
        <li><strong>Minimum Stake</strong> - Required amount to become a validator</li>
        <li><strong>Rewards</strong> - Earned from successful validation</li>
        <li><strong>Commission</strong> - Your fee for running the validator</li>
        <li><strong>Penalties</strong> - For downtime or malicious behavior</li>
      </ul>

      <div className="bg-blue-50 border-l-4 border-blue-500 p-4 my-6">
        <p className="font-semibold mb-2">💡 Tip</p>
        <p className="mb-0">
          Start on the testnet (Sepolia) to familiarize yourself with validator operations before running on mainnet.
        </p>
      </div>

      <h2 className="text-3xl font-semibold mt-16 mb-6">Resources</h2>

      <ul>
        <li><a href="https://docs.starknet.io/" target="_blank" rel="noopener noreferrer">Starknet Documentation</a></li>
        <li><a href="https://t.me/+SCPbza9fk8dkYWI0" target="_blank" rel="noopener noreferrer">Community Telegram</a></li>
        <li><a href="https://github.com/thebuidl-grid/starknode-kit" target="_blank" rel="noopener noreferrer">GitHub Repository</a></li>
      </ul>

      <div className="mt-12 p-6 bg-yellow-50 border border-yellow-200 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">📖 Next Steps</h3>
        <p className="text-gray-700 mb-4">
          Ready to dive deeper? Check out our comprehensive guides:
        </p>
        <div className="flex flex-wrap gap-3">
          <Link href="/requirements" className="text-blue-600 hover:text-blue-800 font-medium">System Requirements</Link>
        </div>
      </div>
    </div>
  );
}

