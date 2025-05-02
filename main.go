package main

import (
    "fmt"
    "os"
    "github.com/caioband/pushy/commands"
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
        commands.RunSSH(os.Args[2:])
    case "show":
        commands.RunShow(os.Args[2:])
    case "help", "--help", "-h":
        printHelp()
    default:
        fmt.Println("❌ Invalid command:", os.Args[1])
        printHelp()
    }
}

func printHelp() {
    fmt.Print(`
🔧 pushy - Lightweight SSH deployment tool

Available commands:

  pushy init               Generate pushy.json via interactive prompts
  pushy deploy             Compress, send, and deploy project to remote server
  pushy config ssh-key     Set SSH private key path (stored in ~/.pushy)
  pushy ssh                Connect directly via SSH using saved config

Example:

  pushy config ssh-key ~/.ssh/deploy_key
  pushy init
  pushy deploy
`)
}
