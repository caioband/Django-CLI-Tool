package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
)

type PushyUserConfig struct {
    SSHKeyPath string `json:"ssh_key_path"`
}

func RunConfig(args []string) {
    if len(args) < 2 {
        fmt.Println("Usage: pushy config ssh-key <path/to/key>")
        return
    }

    subcommand := args[0]

    switch subcommand {
    case "ssh-key":
        path := args[1]
        setSSHKey(path)
    default:
        fmt.Println("Unrecognized config command:", subcommand)
    }
}

func setSSHKey(path string) {
    config := PushyUserConfig{
        SSHKeyPath: path,
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Failed to get user home directory:", err)
        return
    }

    pushyDir := filepath.Join(homeDir, ".pushy")
    configPath := filepath.Join(pushyDir, "config.json")

    // Create ~/.pushy if it doesn't exist
    if _, err := os.Stat(pushyDir); os.IsNotExist(err) {
        if err := os.Mkdir(pushyDir, 0755); err != nil {
            fmt.Println("Failed to create ~/.pushy directory:", err)
            return
        }
    }

    file, err := os.Create(configPath)
    if err != nil {
        fmt.Println("Failed to create config.json:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(config); err != nil {
        fmt.Println("Failed to write to config.json:", err)
        return
    }

    fmt.Println("✅ SSH key path saved to ~/.pushy/config.json")
}
