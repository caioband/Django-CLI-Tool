package commands

import (
    "encoding/json"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)


func RunSSH() {
    // Lê pushy.json
    file, err := os.Open("pushy.json")
    if err != nil {
        fmt.Println("❌ Arquivo pushy.json não encontrado.")
        return
    }
    defer file.Close()

    var project PushyProjectConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&project); err != nil {
        fmt.Println("❌ Erro ao ler pushy.json:", err)
        return
    }

    if project.Host == "" {
        fmt.Println("❌ Host não especificado em pushy.json")
        return
    }

    // Lê config do usuário
    userConfig, err := loadUserConfig()
    if err != nil || userConfig.SSHKeyPath == "" {
        fmt.Println("❌ Caminho da chave SSH não configurado.")
        fmt.Println("Use: pushy config ssh-key <caminho>")
        return
    }

    // Executa ssh
    args := []string{"-i", filepath.Clean(userConfig.SSHKeyPath), project.Host}
    cmd := exec.Command("ssh", args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    fmt.Println("🔐 Conectando com:", project.Host)
    if err := cmd.Run(); err != nil {
        fmt.Println("❌ Erro ao conectar via SSH:", err)
    }
}
