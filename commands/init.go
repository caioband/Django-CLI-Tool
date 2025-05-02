package commands

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "strings"
    "path/filepath"
)

type PushyConfig struct {
    Host        string   `json:"host"`
    RemotePath  string   `json:"remote_path"`
    ArchiveName string   `json:"archive_name"`
    Exclude     []string `json:"exclude"`
    PostDeploy  []string `json:"post_deploy"`
}

func LoadEnvironmentConfig(name string) (*PushyConfig, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	path := filepath.Join(home, ".pushy", "environments", name+".json")

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file '%s': %w", path, err)
	}
	defer file.Close()

	var cfg PushyConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &cfg, nil
}

func SaveEnvironmentConfig(name string, cfg PushyConfig) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	dir := filepath.Join(home, ".pushy", "environments")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	configPath := filepath.Join(dir, name+".json")
	file, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}

	fmt.Printf("✅ Saved environment config to %s\n", configPath)
	return nil
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

    if err := SaveEnvironmentConfig("default", config); err != nil {
        fmt.Println("❌", err)
        return
    }
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
