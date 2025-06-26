## `starknode-kit` CLI Documentation

**Starknode** is a command-line tool to help developers and node operators easily set up, manage, and maintain Ethereum and Starknet nodes.

---

### 📘 Available Commands

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `add`        | Add an Ethereum client to the configuration                |
| `remove`     | Removes a specified resource (client, config, etc.)        |
| `completion` | Generate the autocompletion script for the specified shell |
| `set`        | Set config values for execution or consensus clients       |
| `help`       | Display help about any command                             |

---

### 🧰 Global Flags

| Flag                       | Description                                            |
| -------------------------- | ------------------------------------------------------ |
| `-c`, `--consensus_client` | Specify the consensus client (e.g., Lighthouse, Prysm) |
| `-e`, `--execution_client` | Specify the execution client (e.g., Geth, Reth)        |
| `-s`, `--starknet_client`  | Specify the Starknet client (e.g., Juno)               |
| `-h`, `--help`             | Show help for the `starknode` command                  |

---

### 🧪 Example Usage

#### Generate Config file 

```bash
starknode init
```

#### Add a client pair (consensus + execution)

```bash
starknode add --consensus_client lighthouse --execution_client geth
```

#### Add a Starknet client

```bash
starknode add --starknet_client juno
```

#### Remove a configured client

```bash
starknode remove --consensus_client lighthouse
starknode remove --starknet_client juno
```

#### Set an execution client

```bash
starknodekit set el client=reth network=mainnet port=9000,9001
```

#### Run a Juno Starknet node

```bash
starknode run-juno --network mainnet --port 6060 --data-dir ./juno-data
```

#### Generate bash completion script

```bash
starknode completion bash > /etc/bash_completion.d/starknode
```

#### Show help for a subcommand

```bash
starknode help add
```
