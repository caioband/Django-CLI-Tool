package commands

import (
    "fmt";
    "os";
	"os/exec"
)

func RunInit() {
    run("python", "manage.py", "makemigrations")
    run("python", "manage.py", "migrate")
    run("python", "manage.py", "createsuperuser")
}

func run(name string, args ...string) {
    cmd := exec.Command(name, args...)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    if err := cmd.Run(); err != nil {
        fmt.Println("Erro ao executar:", err)
    }
}