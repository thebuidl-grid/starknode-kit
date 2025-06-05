## `starknode-kit` CLI Documentation

**Starknode** is a command-line tool to help developers and node operators easily set up, manage, and maintain Ethereum and Starknet nodes.

---

### 📘 Available Commands

| Command      | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `add`        | Add an Ethereum client to the configuration                |
| `remove`     | Remove a specified resource (client, config, etc.)         |
| `set`        | Set config values for execution or consensus clients       |
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

### 🔧 `set` Command

The `set` command updates the configuration for execution or consensus clients.

#### Usage:

```bash
starknodekit set el client=reth network=mainnet port=9000,9001
starknodekit set cl client=lighthouse network=mainnet port=9000
```

#### Available keys:

| Key       | Description                                                   |
| --------- | ------------------------------------------------------------- |
| `client`  | Sets the client (e.g., `client=reth`, `client=prysm`)         |
| `network` | Sets the client network (e.g., `network=mainnet`)             |
| `port`    | Comma-separated list of client ports (e.g., `port=9000,9001`) |

#### Aliases:

* `el` – Target the execution client
* `cl` – Target the consensus client

---

### 🧪 Example Usage

#### Add a client pair (consensus + execution)

```bash
starknode add --consensus_client lighthouse --execution_client geth
```

#### Set config for execution client

```bash
starknode set el client=geth network=holesky port=8550,8551
```

#### Set config for consensus client

```bash
starknode set cl client=prysm network=mainnet
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

