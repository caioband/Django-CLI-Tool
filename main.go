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
    case "version", "--version", "-v":
        commands.RunVersion()
    case "help", "--help", "-h":
        printHelp()
    default:
        fmt.Println("❌ Invalid command:", os.Args[1])
        printHelp()
    }
}

func printHelp() {
	fmt.Print(`
🚀 pushy - Lightweight CLI for automated SSH-based deployments

Usage:
  pushy <command> [arguments]

Available commands:
  init                  Create a new environment configuration (default.json)
  deploy [env]          Deploy the current directory using the specified environment
  config ssh-key <path> Set the SSH private key path used for all deployments
  ssh [env]             Open an SSH session with the environment's host
  show [env]            Display the config of a given environment (default is "default")
  version               Show the current version of pushy
  help                  Show this help message

Examples:
  pushy init
  pushy config ssh-key ~/.ssh/id_rsa
  pushy deploy production
  pushy ssh
  pushy show staging
`)
}
