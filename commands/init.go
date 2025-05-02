package commands

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strings"
)

type PushyConfig struct {
    Host        string   `json:"host"`
    RemotePath  string   `json:"remote_path"`
    ArchiveName string   `json:"archive_name"`
    Exclude     []string `json:"exclude"`
    PostDeploy  []string `json:"post_deploy"`
}

func RunInit() {
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("🔒 Remote host (e.g., root@192.168.0.10): ")
    host, _ := reader.ReadString('\n')

    fmt.Print("📂 Remote path (e.g., /var/www/app) [default: ~/]: ")
	remotePath, _ := reader.ReadString('\n')
	remotePath = strings.TrimSpace(remotePath)

    fmt.Print("📦 Archive name (.tar.gz) [default: pushy_deploy.tar.gz]: ")
    archiveName, _ := reader.ReadString('\n')
    if strings.TrimSpace(archiveName) == "" {
        archiveName = "pushy_deploy.tar.gz\n"
    }

    fmt.Print("🚫 Files/folders to exclude (comma-separated): ")
    excludeInput, _ := reader.ReadString('\n')
    excludeList := splitAndTrim(excludeInput)

    if len(excludeList) == 0 {
        excludeList = []string{
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

    fmt.Println("🧩 Post-deploy commands (one per line, empty line to finish):")
    postDeploy := []string{}
    for {
        fmt.Print("> ")
        line, _ := reader.ReadString('\n')
        line = strings.TrimSpace(line)
        if line == "" {
            break
        }
        postDeploy = append(postDeploy, line)
    }

    config := PushyConfig{
        Host:        strings.TrimSpace(host),
        RemotePath:  strings.TrimSpace(remotePath),
        ArchiveName: strings.TrimSpace(archiveName),
        Exclude:     excludeList,
        PostDeploy:  postDeploy,
    }

    file, err := os.Create("pushy.json")
    if err != nil {
        fmt.Println("❌ Failed to create pushy.json:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(config); err != nil {
        fmt.Println("❌ Failed to write pushy.json:", err)
        return
    }

    fmt.Println("✅ pushy.json file created successfully!")
}

func splitAndTrim(input string) []string {
    items := strings.Split(input, ",")
    result := []string{}
    for _, item := range items {
        trimmed := strings.TrimSpace(item)
        if trimmed != "" {
            result = append(result, trimmed)
        }
    }
    return result
}
