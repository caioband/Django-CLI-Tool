package commands

import (
    "archive/tar"
    "compress/gzip"
    "io"
    "bufio"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "path/filepath"
    "encoding/json"
    "runtime"
)

type PushyProjectConfig struct {
    Host        string   `json:"host"`
    RemotePath  string   `json:"remote_path"`
    ArchiveName string   `json:"archive_name"`
    Exclude     []string `json:"exclude"`
    PostDeploy  []string `json:"post_deploy"`
}

func showInstallInstructions() {
    switch runtime.GOOS {
    case "windows":
        fmt.Println("🔹 Windows:")
        fmt.Println("- Enable OpenSSH Client in Windows Features")
        fmt.Println("- Or use: winget install OpenSSH.Client")
    case "linux":
        fmt.Println("🔹 Linux:")
        fmt.Println("- sudo apt install openssh-client")
    case "darwin":
        fmt.Println("🔹 macOS:")
        fmt.Println("- brew install openssh")
    }
}

func EnsureSCP() {
    _, err := exec.LookPath("scp")
    if err == nil {
        return
    }

    fmt.Println("❌ 'scp' was not found on your system.")
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Would you like pushy to try installing it automatically? [y/N]: ")
    answer, _ := reader.ReadString('\n')
    answer = strings.ToLower(strings.TrimSpace(answer))

    if answer != "y" && answer != "yes" {
        fmt.Println("ℹ️ Please install 'scp' manually and try again.")
        showInstallInstructions()
        os.Exit(1)
    }

    fmt.Println("🛠️ Attempting to install 'scp'...")

    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/C", "winget install OpenSSH.Client")
    case "linux":
        cmd = exec.Command("sudo", "apt", "install", "-y", "openssh-client")
    case "darwin":
        cmd = exec.Command("brew", "install", "openssh")
    default:
        fmt.Println("⚠️ Unsupported operating system for automatic installation.")
        os.Exit(1)
    }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        fmt.Println("❌ Failed to install 'scp':", err)
        os.Exit(1)
    }

    fmt.Println("✅ 'scp' installed successfully. Proceeding with deploy.")
}

func loadUserConfig() (*PushyUserConfig, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("failed to get user home directory: %w", err)
    }

    configPath := filepath.Join(homeDir, ".pushy", "config.json")

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("could not open ~/.pushy/config.json: %w", err)
    }
    defer file.Close()

    var config PushyUserConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("failed to parse config.json: %w", err)
    }

    return &config, nil
}

func createTarGz(filename string, exclude []string) error {
    file, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    gz := gzip.NewWriter(file)
    defer gz.Close()

    tarWriter := tar.NewWriter(gz)
    defer tarWriter.Close()

    return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        for _, ex := range exclude {
            if strings.Contains(path, ex) {
                return nil
            }
        }

        header, err := tar.FileInfoHeader(info, path)
        if err != nil {
            return err
        }

        header.Name = path
        if err := tarWriter.WriteHeader(header); err != nil {
            return err
        }

        if !info.Mode().IsRegular() {
            return nil
        }

        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()

        _, err = io.Copy(tarWriter, f)
        return err
    })
}

func RunDeploy(args []string) {
	env := "default"
	if len(args) > 0 {
		env = args[0]
	}

	EnsureSCP()

	// Load SSH key config
	userConfig, _ := loadUserConfig()

	// Load environment config
	project, err := LoadEnvironmentConfig(env)
	if err != nil {
		fmt.Println("❌ Failed to load environment config:", err)
		return
	}

	if len(project.Exclude) == 0 {
		project.Exclude = []string{
			".git", "node_modules", "venv", "__pycache__",
			".idea", ".vscode", ".pushy", "pushy.json",
		}
	}

	archiveName := project.ArchiveName
	if archiveName == "" {
		archiveName = "pushy_deploy.tar.gz"
	}

	// Compress project
	fmt.Println("📦 Compressing project...")
	if err := createTarGz(archiveName, project.Exclude); err != nil {
		fmt.Println("❌ Failed to compress project:", err)
		return
	}

	remotePath := project.RemotePath
	if remotePath == "" {
		remotePath = "."
	}

	// Create remote dir if necessary
	fmt.Println("📁 Ensuring remote directory exists:", remotePath)
	sshArgs := []string{}
	if userConfig != nil && userConfig.SSHKeyPath != "" {
		sshArgs = append(sshArgs, "-i", userConfig.SSHKeyPath)
	}
	sshArgs = append(sshArgs, project.Host, fmt.Sprintf("mkdir -p %s", remotePath))
	ssh := exec.Command("ssh", sshArgs...)
	ssh.Stdout = os.Stdout
	ssh.Stderr = os.Stderr
	ssh.Stdin = os.Stdin
	if err := ssh.Run(); err != nil {
		fmt.Println("❌ Failed to create remote directory:", err)
		return
	}

	// Upload via SCP
	fmt.Println("📤 Uploading to", project.Host)
	scpArgs := []string{}
	if userConfig != nil && userConfig.SSHKeyPath != "" {
		scpArgs = append(scpArgs, "-i", userConfig.SSHKeyPath)
	}
	dest := fmt.Sprintf("%s:%s/", project.Host, strings.TrimSuffix(remotePath, "/"))
	scpArgs = append(scpArgs, archiveName, dest)

	scp := exec.Command("scp", scpArgs...)
	scp.Stdout = os.Stdout
	scp.Stderr = os.Stderr
	scp.Stdin = os.Stdin
	if err := scp.Run(); err != nil {
		fmt.Println("❌ Upload failed:", err)
		return
	}

	// Build post-deploy command with automatic extraction
	postCommands := project.PostDeploy
	archiveCmd := fmt.Sprintf("tar -xzf %s", archiveName)

	// Only inject tar command if user didn’t provide one
	hasTar := false
	for _, cmd := range postCommands {
		if strings.Contains(cmd, "tar -xzf") {
			hasTar = true
			break
		}
	}

	if !hasTar {
		if remotePath != "." && remotePath != "" {
			archiveCmd = fmt.Sprintf("cd %s && %s", remotePath, archiveCmd)
		}
		postCommands = append([]string{archiveCmd}, postCommands...)
	}

	// Execute post-deploy
	if len(postCommands) > 0 {
		fmt.Println("🚀 Running post-deploy commands...")
		sshArgs := []string{}
		if userConfig != nil && userConfig.SSHKeyPath != "" {
			sshArgs = append(sshArgs, "-i", userConfig.SSHKeyPath)
		}
		sshArgs = append(sshArgs, project.Host, strings.Join(postCommands, " && "))

		ssh := exec.Command("ssh", sshArgs...)
		ssh.Stdout = os.Stdout
		ssh.Stderr = os.Stderr
		ssh.Stdin = os.Stdin
		if err := ssh.Run(); err != nil {
			fmt.Println("❌ Post-deploy command failed:", err)
		}
	}

	_ = os.Remove(archiveName)
	fmt.Println("✅ Deploy completed successfully!")
}
