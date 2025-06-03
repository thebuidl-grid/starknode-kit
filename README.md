## `starknode-kit` CLI Documentation

**Starknode** is a command-line tool to help developers and node operators easily set up, manage, and maintain Ethereum and StarkNet nodes.

---

### 📘 Available Commands

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `add`        | Add an Ethereum client to the configuration                |
| `remove`     | Removes a specified resource (client, config, etc.)        |
| `completion` | Generate the autocompletion script for the specified shell |
| `help`       | Display help about any command                             |

---

### 🧰 Global Flags

| Flag                       | Description                                            |
| -------------------------- | ------------------------------------------------------ |
| `-c`, `--consensus_client` | Specify the consensus client (e.g., Lighthouse, Prysm) |
| `-e`, `--execution_client` | Specify the execution client (e.g., Geth, Reth)        |
| `-h`, `--help`             | Show help for the `starknode` command                  |

---

### 🧪 Example Usage

#### Add a client pair (consensus + execution)

```bash
starknode add --consensus_client lighthouse --execution_client geth
```

#### Remove a configured client

```bash
starknode remove --consensus_client lighthouse
```

#### Generate bash completion script

```bash
starknode completion bash > /etc/bash_completion.d/starknode
```

#### Show help for a subcommand

```bash
starknode help add
```
