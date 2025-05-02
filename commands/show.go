package commands

import (
	"fmt"
)

func RunShow(args []string) {
	env := "default"
	if len(args) > 0 {
		env = args[0]
	}

	config, err := LoadEnvironmentConfig(env)
	if err != nil {
		fmt.Println("❌ Failed to load environment config:", err)
		return
	}

	fmt.Printf("📄 Environment: %s\n", env)
	fmt.Println("─────────────────────────────────────────────")
	fmt.Printf("🔒 Host:         %s\n", config.Host)
	fmt.Printf("📂 Remote Path:  %s\n", config.RemotePath)
	fmt.Printf("📦 Archive Name: %s\n", config.ArchiveName)
	fmt.Printf("🚫 Exclude:      %v\n", config.Exclude)
	fmt.Printf("🧩 Post-Deploy:\n")
	if len(config.PostDeploy) == 0 {
		fmt.Println("  (none)")
	} else {
		for _, cmd := range config.PostDeploy {
			fmt.Println("  -", cmd)
		}
	}
	fmt.Println("─────────────────────────────────────────────")
}