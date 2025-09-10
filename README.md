# lazy-open-commity

A command-line tool that automatically generates high-quality, conventional commit messages using GitHub Copilot via the [opencode](https://github.com/sst/opencode) CLI, and commits your staged changes. Designed for seamless integration with [lazygit](https://github.com/jesseduffield/lazygit) for a streamlined, AI-powered Git workflow.

## Features
- Generates commit messages using Copilot, following the conventional commit specification
- Interactive selection of the best commit message
- Commits staged changes with the selected message
- Easy integration with lazygit via custom keybindings

---

## Prerequisites
- [Go](https://golang.org/dl/) (for building this tool)
- [opencode](https://github.com/sst/opencode) (for Copilot integration)
- [lazygit](https://github.com/jesseduffield/lazygit) (optional, for TUI Git)

---

## Installation

### 1. Build lazy-open-commity

```powershell
git clone <this-repo-url>
cd lazy-open-commity
go build -o lazy-open-commity.exe
```

Copy `lazy-open-commity.exe` to a directory in your `PATH` (or reference it directly in your lazygit config).

### 2. Install opencode

Download the latest release for your OS from the [opencode releases page](https://github.com/sst/opencode/releases) and place the binary somewhere in your `PATH`.

### 3. Authenticate opencode with Copilot

Run the following command to authenticate opencode with GitHub Copilot:

```powershell
opencode auth login
```

Follow the prompts to complete authentication. You must have an active Copilot subscription.

---

## Integrating with lazygit

To trigger `lazy-open-commity` directly from lazygit, add a custom command to your lazygit config file:

Edit (or create) `C:\Users\<YourUsername>\AppData\Local\lazygit\config.yml` and add:

```yaml
customCommands:
  - key: '<c-s>'
    context: 'files'
    command: 'lazy-open-commity'
    description: 'Automatically generate commit message and commit changes'
    output: terminal
```

- Press `Ctrl+S` in the files view to run the tool and commit staged changes with an AI-generated message.

---

## Usage

1. Stage your changes in lazygit (or via `git add`)
2. Press your configured keybinding (e.g., `Ctrl+S`) in lazygit, or run `lazy-open-commity` in your terminal
3. Select a commit message from the Copilot-generated suggestions
4. The tool will commit your changes with the selected message

---

## Troubleshooting
- Ensure both `opencode` and `lazy-open-commity` are in your `PATH`
- Make sure you have staged changes before running the tool
- If Copilot authentication fails, re-run `opencode auth login`

---

## License
MIT
