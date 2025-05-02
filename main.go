package main

import (
    "fmt"
    "os"
    "djctl/commands"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Use: djctl <comando> [argumentos]")
        fmt.Println("Comandos disponíveis: new, app, init")
        return
    }

    switch os.Args[1] {
    case "new":
        //commands.RunNew(os.Args[2:])
    case "app":
        //commands.RunApp(os.Args[2:])
    case "init":
        commands.RunInit()
    default:
        fmt.Println("Comando inválido:", os.Args[1])
        fmt.Println("Comandos disponíveis: new, app, init")
    }
}