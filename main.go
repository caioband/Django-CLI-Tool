package main

import (
    "fmt"
    "os"
    "pushy/commands"
)

func main() {
    if len(os.Args) < 2 {
        printHelp()
        return
    }

    switch os.Args[1] {
    case "init":
        commands.RunInit()
    case "deploy":
        commands.RunDeploy(os.Args[2:])
    case "config":
        commands.RunConfig(os.Args[2:])
    case "ssh":
        commands.RunSSH()
    case "help", "--help", "-h":
        printHelp()
    default:
        fmt.Println("❌ Comando inválido:", os.Args[1])
        printHelp()
    }
}

func printHelp() {
fmt.Print(`
🔧 pushy - Facilitador de Deploys via SSH

Comandos disponíveis:

  pushy init               Inicia o pushy.json com perguntas interativas
  pushy deploy             Envia o projeto para o servidor remoto via scp
  pushy config ssh-key     Define o caminho da chave SSH (salva em ~/.pushy)

Exemplo:

  pushy config ssh-key ~/.ssh/deploy_key
  pushy init
  pushy deploy
`)
}