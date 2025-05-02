# Pushy

**Pushy** is a lightweight CLI tool to simplify project deployment via SSH.

## ✨ Features

* Interactive configuration setup
* Automatic SSH file transfers using `scp`
* Automatic unpacking and optional post-deploy commands
* Environment management via `~/.pushy/environments`

## 🛠️ Requirements

* Go 1.18 or higher
* SSH access and private key
* `scp` and `ssh` installed (Pushy can auto-install them on most systems)

## 🔧 Installation

Install via Go:

```bash
go install github.com/caioband/pushy@latest
```

## ⚡ Usage

### 1. Set SSH Key Path

```bash
pushy config ssh-key ~/.ssh/deploy_key
```

### 2. Initialize a New Environment

```bash
pushy init
```

This command will prompt for deployment settings and create a config file at:

```
~/.pushy/environments/default.json
```

### 3. Deploy the Project

```bash
pushy deploy
```

This will:

* Compress the current directory into a `.tar.gz`
* Transfer the archive to the configured server
* Optionally run post-deploy commands like `tar -xzf` or `cd` commands

### 4. Connect via SSH

```bash
pushy ssh
```

Connects to the remote server using your saved SSH key and hostname.

### 5. Check Current Configuration

```bash
pushy show <enviroment>
```

Displays the active deployment settings stored in your environment.

### 6. Check Version

```bash
pushy version
```

Shows the installed Pushy version.

## 💡 Example

```bash
pushy config ssh-key ~/.ssh/deploy_key
pushy init
pushy deploy
```

## 📂 Configuration File Format

```json
{
  "host": "user@yourserver.com",
  "remote_path": "pushy",
  "archive_name": "project.tar.gz",
  "exclude": [".git", "node_modules", "venv", "__pycache__"],
  "post_deploy": [
    "tar -xzf project.tar.gz",
    "cd project"
  ]
}
```

## 🤝 Contributing

We welcome contributions to Pushy!

To contribute:

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Make your changes and commit them
4. Open a pull request with a clear description

Please ensure your changes are consistent with the style of the project and include tests or usage examples when applicable.

## ✉ License

This project is licensed under the MIT License.

## 📬 Contact

* Email: [caioobsantos@gmail.com](mailto:caioobsantos@gmail.com)
* GitHub: [github.com/caioband](https://github.com/caioband)
