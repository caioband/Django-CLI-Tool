# 🚀 pushy

**pushy** is a modern CLI tool written in Go to simplify project deployment via SSH.  
It automates packaging, secure copy (`scp`), remote command execution, and configuration — making deployment effortless and consistent across environments.

---

## ✨ Features

- 📦 Automatically compress the current project into `.tar.gz`
- 🔐 Supports custom SSH key configuration
- 📤 Deploy to remote servers using `scp`
- 🧩 Execute post-deploy commands remotely via `ssh`
- ⚙️ Interactive `init` command to generate configuration file
- 🛠 `config` command to manage local settings
- 💻 Cross-platform: Windows, Linux, macOS

---

## 📥 Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-username/pushy.git
cd pushy
```

### 2. Build the binary

#### 💻 Linux/macOS

```bash
go build -o pushy
```

#### 🪟 Windows

```bash
go build -o pushy.exe
```

To cross-compile:

```bash
GOOS=linux GOARCH=amd64 go build -o pushy-linux
GOOS=windows GOARCH=amd64 go build -o pushy.exe
```

---

## ⚙️ Usage

### 1. Configure your SSH key (optional but recommended)

```bash
pushy config ssh-key ~/.ssh/id_rsa
```

### 2. Initialize the project configuration file

```bash
pushy init
```

### 3. Deploy your project

```bash
pushy deploy
```

---

## 🧪 Available Commands

```bash
pushy init                     # Generate pushy.json interactively
pushy deploy                   # Compress, send and run remote post-deploy commands
pushy config ssh-key <path>    # Set SSH private key path
pushy ssh                      # Open direct SSH session using saved config
```

---

## 📁 Example pushy.json

```json
{
  "host": "user@your-server.com",
  "remote_path": "/var/www/myapp",
  "archive_name": "pushy_deploy.tar.gz",
  "exclude": [".git", "node_modules", "pushy.json"],
  "post_deploy": [
    "cd /var/www/myapp",
    "tar -xzf pushy_deploy.tar.gz",
    "rm pushy_deploy.tar.gz"
  ]
}
```

---

## 🔐 Security

- The SSH key path is stored in `~/.pushy/config.json`
- Ensure your private key has restricted permissions:

```bash
chmod 400 ~/.ssh/your-key.pem
```


---

## 📦 Requirements

- Go 1.18+ installed and available in your system's PATH
- Access to a remote server with:
  - OpenSSH server running (port 22 by default)
  - Writable access to the specified remote_path
- (Optional) SSH key properly configured and whitelisted on the server
- SCP (secure copy) installed locally (automatically suggested if missing)


---

## 📄 License

This project is licensed under the **MIT License**.  
Feel free to use, modify, and contribute.

---

## 🤝 Contributing

Contributions are welcome!  
Open a [pull request](https://github.com/your-username/pushy/pulls) or [issue](https://github.com/your-username/pushy/issues) to collaborate.

---

## ✉️ Contact

Made with 💻 by **Your Name**  
📧 Email: you@example.com  
🌐 GitHub: [https://github.com/your-username](https://github.com/your-username)
