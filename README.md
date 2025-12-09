# Capturl MCP Server

[![Go Report Card](https://goreportcard.com/badge/github.com/capturl/capturl-mcp)](https://goreportcard.com/report/github.com/capturl/capturl-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Latest Release](https://img.shields.io/github/v/release/capturl/capturl-mcp)](https://github.com/capturl/capturl-mcp/releases)

**The easiest way to pass screenshots to LLMs.**

The Official [Model Context Protocol (MCP)](https://modelcontextprotocol.io/) server for **capturl.com**. This tool allows AI Agents (Cursor, Claude, or any MCP client) to directly fetch screenshots using capturl URLs.

> Note: This MCP server only works with unencrypted Capturls. Encrypted Capturls are not supported.

---

## 🚀 Quick Start (Cursor)

Setting up Capturl in [Cursor](https://cursor.com) takes two steps:

**Step 1: Install the binary**
You need the `capturl-mcp` tool installed on your machine first.
```bash
# MacOS / Linux
brew tap capturl/tap
brew install --cask capturl-mcp
```
*(See below for other installation methods)*

**Step 2: Add to Cursor**
Click the link below to automatically add the server to your Cursor MCP settings:

[![Install MCP Server](https://cursor.com/deeplink/mcp-install-dark.svg)](https://cursor.com/en-US/install-mcp?name=capturl-mcp&config=eyJjb21tYW5kIjoiY2FwdHVybC1tY3AiLCJlbnYiOnsiQ0FQVFVSTF9BVVRIX1RPS0VOIjoiZ2V0X3RoaXNfZnJvbV9odHRwczovL2NhcHR1cmwuY29tL3Byb2ZpbGUvc2V0dGluZ3MifX0%3D)

> **Important:** This link pre-fills the configuration, but you must replace the placeholder text with your actual **Auth Token** (found in your [Capturl Settings](https://capturl.com/profile/settings)).

---

## 📦 Installation Options

If you haven't installed the binary yet, choose the method that works best for your workflow.

### Option 1: Homebrew (Recommended for macOS/Linux)
The easiest way to install and keep the CLI updated.

```bash
brew tap capturl/tap
brew install --cask capturl-mcp
```

### Option 2: Go Install (For Go Developers)
Install directly using the Go toolchain (requires Go 1.25+).

```bash
go install [github.com/capturl/capturl-mcp@latest](https://github.com/capturl/capturl-mcp@latest)
```

### Option 3: Build from Source
If you prefer to build the binary manually or are developing on the project:

```bash
# Clone the repository
git clone [https://github.com/capturl/capturl-mcp.git](https://github.com/capturl/capturl-mcp.git)
cd capturl-mcp

# Install to /usr/local/bin (requires sudo)
sudo make install

# OR install to a local directory (no sudo required)
make install BIN_DIR=$HOME/.local/bin
```

---

## 🔑 Configuration

To use the MCP server, you must authenticate with your Capturl account.

1.  **Get your Token:** Navigate to [capturl.com/profile/settings](https://capturl.com/profile/settings) and copy your **Refresh Token**.

2.  **Set the Environment Variable:** The server requires the `CAPTURL_AUTH_TOKEN` variable to be set.

    *If running manually:*
    ```bash
    export CAPTURL_AUTH_TOKEN="your_token_here"
    ```

    *If using with an MCP Client:* Add the environment variable to your client configuration file (e.g., `claude_desktop_config.json` or Cursor settings).

---

## 🛠️ Usage

Once installed, the `capturl-mcp` tool allows AI agents to interact with https://capturl.com. Things to try:

> "What does this mean? https://capturl.com/o/pro/c/some_id"

> "I'm seeing an error in my app, can you fix it? Here's the capturl: https://capturl.com/o/pro/c/some_id"

> "Can you fix the padding here: https://capturl.com/o/pro/c/some_id"

---

## 💻 Development

### Requirements
* **Go:** v1.25.1 or later
* **Make:** For running build scripts
* **GoReleaser** If you're messing with the goreleaser config.

### Common Commands
* `make build`
* `make test`
* `make clean`

### Making a new release
1. Tag the release: `git tag v0.0.1 && git push origin v0.0.1`
2. [Testing locally]: `goreleaser release --snapshot --clean`
3. `export GITHUB_TOKEN="get_your_token_from_https://github.com/settings/tokens/new" && goreleaser release --clean`