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

    fmt.Print("🔒 Host remoto (ex: root@192.168.0.10): ")
    host, _ := reader.ReadString('\n')

    fmt.Print("📂 Caminho remoto de destino (ex: /var/www/app): ")
    remotePath, _ := reader.ReadString('\n')

    fmt.Print("📦 Nome do arquivo .tar.gz (padrão: pushy_deploy.tar.gz): ")
    archiveName, _ := reader.ReadString('\n')
    if strings.TrimSpace(archiveName) == "" {
        archiveName = "pushy_deploy.tar.gz\n"
    }

    fmt.Print("🚫 Pastas/arquivos para ignorar (separados por vírgula): ")
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

    fmt.Println("🧩 Comandos pós-deploy (digite 1 por linha, ENTER vazio para finalizar):")
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
        fmt.Println("❌ Erro ao criar pushy.json:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(config); err != nil {
        fmt.Println("❌ Erro ao escrever pushy.json:", err)
        return
    }

    fmt.Println("✅ Arquivo pushy.json criado com sucesso!")
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