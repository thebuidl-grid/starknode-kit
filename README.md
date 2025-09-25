# `starknode-kit` CLI Documentation

**starknode-kit** is a command-line tool to help developers and node operators easily set up, manage, and maintain Ethereum and Starknet nodes.

---

## 🚀 Installation

### Option 1: Install using the install script (Recommended)

1. Download and run the installation script:

   ```bash
   curl -sSL https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh | bash
   ```

2. Or download the script first and then run it:

   ```bash
   wget https://raw.githubusercontent.com/thebuidl-grid/starknode-kit/main/install.sh
   chmod +x install.sh
   ./install.sh
   ```

### Option 2: Install using Go

Make sure you have Go installed (version 1.24 or later). This method installs the latest version from the `main` branch.

```bash
go install -ldflags="-X 'github.com/thebuidl-grid/starknode-kit/pkg/versions.StarkNodeVersion=main'" github.com/thebuidl-grid/starknode-kit@latest
```

### Option 3: Manual Installation from Source

1. Clone the repository:

   ```bash
   git clone https://github.com/thebuidl-grid/starknode-kit.git
   cd starknode-kit
   ```

2. Build and install:

   ```bash
   make build
   sudo mv bin/starknode /usr/local/bin/
   ```

### Verify Installation

After installation, verify that `starknode-kit` is working:

```bash
starknode-kit --help
```
#### Generate Config file

```bash
starknode-kit config new
```

### 🧹 Uninstallation

To uninstall `starknode-kit`, remove the binary and the configuration directory:

```bash
sudo rm /usr/local/bin/starknode-kit
rm -rf ~/.starknode-kit
```

> **Note**: This will not remove any of the client data (e.g., blockchain data). The data is stored in the locations specified in your `~/.starknode-kit/starknode.yml` file.


---

## 📘 Available Commands

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `add`        | Add an Ethereum or Starknet client to the config           |
| `completion` | Generate the autocompletion script for the specified shell |
| `config`     | Create, show, and update your Starknet node configuration. |
| `help`       | Display help about any command                             |
| `monitor`    | Launch real-time monitoring dashboard                      |
| `remove`     | Remove a specified resource                                |
| `run`        | Run a specific local infrastructure service                |
| `status`     | Display status of running clients                          |
| `start`      | Run the configured Ethereum clients                        |
| `stop`       | Stop the configured Ethereum clients                       |
| `update`     | Check for and install client updates                       |
| `validator`  | Manage the Starknet validator client                       |
| `version`    | Show version of starknode-kit or a specific client         |

---

## 🧪 Example Usage

#### Add a client pair (consensus + execution)

```bash
starknode-kit add --consensus_client lighthouse --execution_client geth
```

#### Add a Starknet client

```bash
starknode-kit add --starknet_client juno
```

#### Remove a configured client

```bash
starknode-kit remove --consensus_client lighthouse
starknode-kit remove --starknet_client juno
```

#### Change network

```bash
starknode-kit config set network sepolia
```

#### Set an execution client

```bash
starknode-kit config set el client=reth port=9000,9001
```

#### Show configuration

```bash
starknode-kit config show --all
starknode-kit config show --el
```

#### Check version

```bash
starknode-kit version
starknode-kit version geth
```

#### Start Ethereum clients

```bash
starknode-kit start
```

> ⚠️ **Note**: The `start` command only launches the configured **execution (EL)** and **consensus (CL)** clients. It does **not** start any Starknet clients.

#### Run a specific client

To run a specific client using its configured settings:

```bash
starknode-kit run juno
starknode-kit run geth
starknode-kit run lighthouse
```

#### Validator Commands

Manage the Starknet validator client.

- **Get validator status:**
  ```bash
  starknode-kit validator status
  ```

- **Get validator version:**
  ```bash
  starknode-kit validator --version
  ```

- **Set Juno RPC endpoint:**
  ```bash
  starknode-kit validator --rpc <YOUR_RPC_URL>
  ```

#### Generate bash completion script

```bash
starknode-kit completion bash > /etc/bash_completion.d/starknode-kit
```

#### Show help for a subcommand

```bash
starknode-kit help add
```

---

## 📋 Requirements

### 🛠️ Software Dependencies

Make sure the following are installed on your system before using or building `starknode-kit`:

* **Go**: Version **1.24 or later**
  Install from: [https://go.dev/dl/](https://go.dev/dl/)

* **Rust**: Recommended for building Starknet clients (e.g., Juno)
  Install with:

  ```bash
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
  ```

* **Make**: Required to build certain clients and scripts
  Install via package manager:

  * Ubuntu/Debian: `sudo apt install make`
  * macOS (with Homebrew): `brew install make`
  * Windows (WSL): included or `sudo apt install make`

### 🖥️ Hardware Requirements

See this [Rocket Pool Hardware Guide](https://docs.rocketpool.net/guides/node/hardware.html) for a detailed breakdown of node hardware requirements.

* **CPU**: Node operation doesn't require heavy CPU power. The BG Client has run well on both i3 and i5 models of the ASUS NUC 13 PRO. Be cautious if using Celeron processors, as they may have limitations.
* **RAM**: At least **32 GB** is recommended for good performance with overhead.
* **Storage (SSD)**: The most critical component. Use a **2 TB+ NVMe SSD** with:

  * A **DRAM cache**
  * **No Quad-Level Cell (QLC)** NAND architecture
    See this [SSD List Gist](https://gist.github.com/bkase/fab02c5b3c404e9ef8e5c2071ac1558c) for tested options.

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🌐 Join the Community

Join the community to stay updated, ask questions, or contribute:

- Telegram: [https://t.g/+-SCPbza9fk8dkYWI0](https://t.me/+SCPbza9fk8dkYWI0)

Whether you're a seasoned validator, hobbyist, or first-time node runner, you're welcome!

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Uninstallation

```bash
sudo rm /usr/local/bin/starknode-kit
rm -rf ~/.starknode-kit
```