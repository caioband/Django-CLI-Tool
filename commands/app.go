package commands

import (
    "fmt"
    "os"
    "os/exec"
)

func RunApp(args []string) {
    if len(args) < 1 {
        fmt.Println("Uso: djctl app <nome_do_app>")
        return
    }

    name := args[0]
    cmd := exec.Command("python", "manage.py", "startapp", name)

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    if err := cmd.Run(); err != nil {
        fmt.Println("Erro ao criar app:", err)
    }
}