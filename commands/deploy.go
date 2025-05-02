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

func exibirInstrucoes() {
    switch runtime.GOOS {
    case "windows":
        fmt.Println("🔹 Windows:")
        fmt.Println("- Ative o OpenSSH Client nas Configurações do Windows")
        fmt.Println("- Ou use: winget install OpenSSH.Client")
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
        return // scp está disponível
    }

    fmt.Println("❌ 'scp' não foi encontrado no sistema.")
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Deseja que o pushy tente instalar automaticamente? [s/N]: ")
    resposta, _ := reader.ReadString('\n')
    resposta = strings.ToLower(strings.TrimSpace(resposta))

    if resposta != "s" && resposta != "sim" {
        fmt.Println("ℹ️  Instale o 'scp' manualmente e tente novamente.")
        exibirInstrucoes()
        os.Exit(1)
    }

    fmt.Println("🛠️ Tentando instalar 'scp'...")

    var cmd *exec.Cmd

    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/C", "winget install OpenSSH.Client")
    case "linux":
        cmd = exec.Command("sudo", "apt", "install", "-y", "openssh-client")
    case "darwin":
        cmd = exec.Command("brew", "install", "openssh")
    default:
        fmt.Println("⚠️ Sistema não suportado para instalação automática.")
        os.Exit(1)
    }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        fmt.Println("❌ Falha ao instalar 'scp':", err)
        os.Exit(1)
    }

    fmt.Println("✅ 'scp' instalado com sucesso. Continue com o deploy.")
}


func loadUserConfig() (*PushyUserConfig, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("erro ao obter diretório do usuário: %w", err)
    }

    configPath := filepath.Join(homeDir, ".pushy", "config.json")

    file, err := os.Open(configPath)
    if err != nil {
        return nil, fmt.Errorf("não foi possível abrir ~/.pushy/config.json: %w", err)
    }
    defer file.Close()

    var config PushyUserConfig
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("erro ao decodificar config.json: %w", err)
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

        // Ignora arquivos da lista
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

    // Carrega config global do usuário
    userConfig, _ := loadUserConfig()

    // Lê pushy.json
    configFile, err := os.Open("pushy.json")
    if err != nil {
        fmt.Println("❌ Arquivo pushy.json não encontrado.")
        return
    }
    defer configFile.Close()

    var project PushyProjectConfig
    decoder := json.NewDecoder(configFile)
    if err := decoder.Decode(&project); err != nil {
        fmt.Println("❌ Erro ao ler pushy.json:", err)
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

    // Cria .tar.gz excluindo arquivos/pastas definidos
    fmt.Println("📦 Compactando projeto...")
    if err := createTarGz(archiveName, project.Exclude); err != nil {
        fmt.Println("❌ Erro ao compactar:", err)
        return
    }

    // Monta comando SCP
    scpArgs := []string{}
    if userConfig != nil && userConfig.SSHKeyPath != "" {
        scpArgs = append(scpArgs, "-i", userConfig.SSHKeyPath)
    }
    scpArgs = append(scpArgs, archiveName, fmt.Sprintf("%s:%s", project.Host, project.RemotePath))

    fmt.Println("📤 Enviando para", project.Host)
    scp := exec.Command("scp", scpArgs...)
    scp.Stdout = os.Stdout
    scp.Stderr = os.Stderr
    scp.Stdin = os.Stdin
    if err := scp.Run(); err != nil {
        fmt.Println("❌ Erro ao enviar:", err)
        return
    }

    // Executa post-deploy via SSH
    if len(project.PostDeploy) > 0 {
        sshArgs := []string{}
        if userConfig != nil && userConfig.SSHKeyPath != "" {
            sshArgs = append(sshArgs, "-i", userConfig.SSHKeyPath)
        }
        sshArgs = append(sshArgs, project.Host, strings.Join(project.PostDeploy, " && "))

        fmt.Println("🚀 Executando comandos pós-deploy...")
        ssh := exec.Command("ssh", sshArgs...)
        ssh.Stdout = os.Stdout
        ssh.Stderr = os.Stderr
        ssh.Stdin = os.Stdin
        if err := ssh.Run(); err != nil {
            fmt.Println("❌ Erro no pós-deploy:", err)
        }
    }

    // Remove arquivo local
    _ = os.Remove(archiveName)
    fmt.Println("✅ Deploy finalizado com sucesso!")
}