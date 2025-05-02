package commands

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
)

func RunSSH(args []string) {
	env := "default"
	if len(args) > 0 {
		env = args[0]
	}

	// Load environment config from ~/.pushy/environments/<env>.json
	project, err := LoadEnvironmentConfig(env)
	if err != nil {
		fmt.Println("❌ Failed to load environment config:", err)
		return
	}

	if project.Host == "" {
		fmt.Println("❌ Host not specified in environment config.")
		return
	}

	// Load global SSH key path
	userConfig, err := loadUserConfig()
	if err != nil || userConfig.SSHKeyPath == "" {
		fmt.Println("❌ SSH key path not configured.")
		fmt.Println("Use: pushy config ssh-key <path>")
		return
	}

	// Build and run SSH command
	argsList := []string{"-i", filepath.Clean(userConfig.SSHKeyPath), project.Host}
	cmd := exec.Command("ssh", argsList...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	fmt.Println("🔐 Connecting to:", project.Host)
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ SSH connection failed:", err)
	}
}
