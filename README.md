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

Make sure you have Go installed (version 1.19 or later), then run:

```bash
go install github.com/thebuidl-grid/starknode-kit@latest
```

### Option 3: Manual Installation from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/thebuidl-grid/starknode-kit.git
   cd starknode-kit
   ```

2. Build and install:
   ```bash
   go build -o starknode-kit .
   sudo mv starknode-kit /usr/local/bin/
   ```

### Verify Installation

After installation, verify that `starknode-kit` is working:
```bash
starknode-kit --help
```

---

## 📘 Available Commands

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `add`        | Add an Ethereum or Starknet client to the config           |
| `completion` | Generate the autocompletion script for the specified shell |
| `config`     | Show the configured Ethereum clients                       |
| `help`       | Display help about any command                             |
| `init`       | Create a default configuration file                        |
| `monitor`    | Launch real-time monitoring dashboard                      |
| `remove`     | Remove a specified resource                                |
| `run`        | Run local Starknet infrastructure services                 |
| `set`        | Set config values for execution or consensus clients       |
| `start`      | Run the configured Ethereum clients                        |
| `stop`       | Stop the configured Ethereum clients                       |
| `update`     | Check for and install client updates                       |

---

## 🧰 Global Flags

| Flag                       | Description                                            |
| -------------------------- | ------------------------------------------------------ |
| `-c`, `--consensus_client` | Specify the consensus client (e.g., Lighthouse, Prysm) |
| `-e`, `--execution_client` | Specify the execution client (e.g., Geth, Reth)        |
| `-s`, `--starknet_client`  | Specify the Starknet client (e.g., Juno)               |
| `-h`, `--help`             | Show help for the `starknode-kit` command             |

---

## 🧪 Example Usage

#### Generate Config file 
```bash
starknode-kit init
```

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

#### Set an execution client
```bash
starknode-kit set el client=reth network=mainnet port=9000,9001
```

#### Run a Juno Starknet node
```bash
starknode-kit run juno --network mainnet --port 6060 --data-dir ./juno-data
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

- **Go**: Version 1.24^ 
- **Operating System**: Linux, macOS, or Windows (WSL recommended for Windows)

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
