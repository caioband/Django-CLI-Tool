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
        fmt.Println("Uso: pushy config ssh-key <caminho/da/chave>")
        return
    }

    subcommand := args[0]

    switch subcommand {
    case "ssh-key":
        caminho := args[1]
        setSSHKey(caminho)
    default:
        fmt.Println("Configuração não reconhecida:", subcommand)
    }
}

func setSSHKey(path string) {
    config := PushyUserConfig{
        SSHKeyPath: path,
    }

    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Erro ao obter diretório do usuário:", err)
        return
    }

    pushyDir := filepath.Join(homeDir, ".pushy")
    configPath := filepath.Join(pushyDir, "config.json")

    // Cria o diretório ~/.pushy se não existir
    if _, err := os.Stat(pushyDir); os.IsNotExist(err) {
        if err := os.Mkdir(pushyDir, 0755); err != nil {
            fmt.Println("Erro ao criar ~/.pushy:", err)
            return
        }
    }

    file, err := os.Create(configPath)
    if err != nil {
        fmt.Println("Erro ao criar config.json:", err)
        return
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    if err := encoder.Encode(config); err != nil {
        fmt.Println("Erro ao escrever config.json:", err)
        return
    }

    fmt.Println("✅ Caminho da chave SSH salvo em ~/.pushy/config.json")
}