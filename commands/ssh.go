package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func RunSSH() {
    // Load pushy.json
    file, err := os.Open("pushy.json")
    if err != nil {
        fmt.Println("❌ pushy.json file not found.")
        return
    }
    defer file.Close()

    var project PushyProjectConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&project); err != nil {
        fmt.Println("❌ Failed to read pushy.json:", err)
        return
    }

    if project.Host == "" {
        fmt.Println("❌ Host not specified in pushy.json.")
        return
    }

    // Load user's SSH key config
    userConfig, err := loadUserConfig()
    if err != nil || userConfig.SSHKeyPath == "" {
        fmt.Println("❌ SSH key path not configured.")
        fmt.Println("Use: pushy config ssh-key <path>")
        return
    }

    // Build and execute ssh command
    args := []string{"-i", filepath.Clean(userConfig.SSHKeyPath), project.Host}
    cmd := exec.Command("ssh", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    fmt.Println("🔐 Connecting to:", project.Host)
    if err := cmd.Run(); err != nil {
        fmt.Println("❌ Failed to connect via SSH:", err)
    }
}
