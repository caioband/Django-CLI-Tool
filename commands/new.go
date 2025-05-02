package commands

import (
    "fmt"
    "os"
    "os/exec"
)

func RunNew(args []string) {
    if len(args) < 1 {
        fmt.Println("Uso: djctl new <nome_do_projeto>")
        return
    }

    name := args[0]
    cmd := exec.Command("django-admin", "startproject", name)

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        fmt.Println("Erro ao criar projeto:", err)
    }
}