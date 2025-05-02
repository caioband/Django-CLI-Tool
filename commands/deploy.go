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
    EnsureSCP()

    userConfig, _ := loadUserConfig()

    configFile, err := os.Open("pushy.json")
    if err != nil {
        fmt.Println("❌ pushy.json file not found.")
        return
    }
    defer configFile.Close()

    var project PushyProjectConfig
    decoder := json.NewDecoder(configFile)
    if err := decoder.Decode(&project); err != nil {
        fmt.Println("❌ Failed to read pushy.json:", err)
        return
    }

    if len(project.Exclude) == 0 {
        project.Exclude = []string{
            ".git",
            "node_modules",
            "venv",
            "__pycache__",
            ".idea",
            ".vscode",
            ".pushy",
            "pushy.json",
        }
    }

    archiveName := project.ArchiveName
    if archiveName == "" {
        archiveName = "pushy_deploy.tar.gz"
    }

    fmt.Println("📦 Compressing project...")
    if err := createTarGz(archiveName, project.Exclude); err != nil {
        fmt.Println("❌ Failed to compress project:", err)
        return
    }

    var remoteDest string
    if project.RemotePath == "" {
        remoteDest = "~" // Home directory
    } else {
        remoteDest = fmt.Sprintf("~/%s", strings.TrimPrefix(project.RemotePath, "/"))
    }

    // Create remote directory
    if userConfig != nil && userConfig.SSHKeyPath != "" {
        fmt.Println("📁 Ensuring remote directory exists:", remoteDest)
        mkdirArgs := []string{"-i", userConfig.SSHKeyPath, project.Host, "mkdir", "-p", remoteDest}
        mkdir := exec.Command("ssh", mkdirArgs...)
        mkdir.Stdout = os.Stdout
        mkdir.Stderr = os.Stderr
        if err := mkdir.Run(); err != nil {
            fmt.Println("❌ Failed to create remote directory:", err)
            return
        }
    }

    // Build final scp command
    scpArgs := []string{"-i", userConfig.SSHKeyPath, archiveName, fmt.Sprintf("%s:%s", project.Host, remoteDest)}

    fmt.Println("📤 Uploading to", project.Host)
    scp := exec.Command("scp", scpArgs...)
    scp.Stdout = os.Stdout
    scp.Stderr = os.Stderr
    scp.Stdin = os.Stdin
    if err := scp.Run(); err != nil {
        fmt.Println("❌ Failed to upload archive:", err)
        return
    }

    for i, cmd := range project.PostDeploy {
        if strings.Contains(cmd, project.RemotePath) && project.RemotePath != "" {
            project.PostDeploy[i] = strings.ReplaceAll(
                cmd,
                project.RemotePath,
                remoteDest,
            )
        }
    }
    

    if len(project.PostDeploy) > 0 {
        sshArgs := []string{}
        if userConfig != nil && userConfig.SSHKeyPath != "" {
            sshArgs = append(sshArgs, "-i", userConfig.SSHKeyPath)
        }
        sshArgs = append(sshArgs, project.Host, strings.Join(project.PostDeploy, " && "))

        fmt.Println("🚀 Running post-deploy commands...")
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
