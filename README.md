# Gosh

> A lightweight, cross-platform shell written in Go — started as a learning project, now a growing bash/zsh alternative.

Gosh began as part of the [CodeCrafters](https://codecrafters.io) shell challenge, but it quickly became much more. While it doesn’t yet rival bash or zsh in every feature, it implements a strong foundation and several modern conveniences — all in pure Go.

---

## ✨ Features

### ✅ Implemented

- **Command Execution** – Run external binaries or use built-in commands.
- **Autocompletion** – Intelligent suggestions for commands and paths.
- **Aliases** – Define your own short commands via config.
- **History** – Navigate and recall previously run commands.
- **Configuration File** – Customize Gosh behavior with a simple config.
- **Cross-Platform Support** – Runs on Linux, macOS, and Windows.
- **Clean Prompt UI** – Simple, readable, and minimalistic prompt.
- **Logging** – Built-in logger for debugging and development.
- **Environment Variable Management** – `export`, `$FOO`, etc.

### 🚧 Not Yet Implemented (but planned)

- **Piping & Redirection** – `ls | grep foo` and friends.
- **Background Jobs** – Support for `&` and process control.
- **Temporary Variable Assignments** – Support for `VAR=test echo $VAR`.
- **Autocompletion for Environment Variables** – Better support for `$VAR` suggestions.
- **Theme/Color Configurations** – Customizable appearance for the prompt.

---

## 🚀 Getting Started

### Requirements

- Go 1.20+

### Install

Clone and build it yourself:

```bash
git clone https://github.com/SebastianRichiteanu/Gosh.git
cd gosh
make build
```

### Run

```bash
./gosh
```

---

## 💻 Example Usage

```bash
$ echo Hello, Gosh!
Hello, Gosh!

$ ls
# lists directory contents

$ alias gs='git status'
$ gs > git_status.txt
# shows git status using alias, outputing the std in the git_status file
```

Autocompletion and history navigation with arrow keys work out of the box.

---

## ⚙️ Configuration

Gosh reads its settings from the ~/.gosh/goshrc file. This file is created automatically on first run if it doesn’t exist.

You can customize behavior by setting environment variables in this file like so:

```
# Change the shell prompt symbol
export GOSH_SHELL_SYMBOL=">"

# Set the logging level (e.g., DEBUG, INFO, WARN, ERROR)
export GOSH_LOG_LEVEL="INFO"

# Enable or disable autocompletion (true or false)
export GOSH_ENABLE_AUTOCOMPLETE=true

# Set custom log file
export GOSH_LOG_FILE="gosh.log"

# Set custom history file
export GOSH_HISTORY_FILE="history"

# Set custom alias file
export GOSH_ALIAS_FILE="aliases"

# Limit the number of saved history entries
export GOSH_MAX_HISTORY_SIZE=1337
```

## 🗂️ Project Structure

```
Gosh/
├── cmd/gosh/           # Entry point
├── internal/
│   ├── autocompleter/  # Suggests commands and paths
│   ├── builtins/       # Built-in shell commands
│   ├── config/         # Config file parsing and defaults
│   ├── executor/       # Command execution logic
│   ├── prompt/         # REPL/prompt UI and input
│   ├── logger/, utils/ # Utilities
├── tests/              # Unit tests
```

---

## 🤝 Contributing

This project started as a way to learn — contributions and learning together are welcome! Open an issue or PR any time.

---

## 📄 License

MIT © Sebastian Richiteanu 2025
