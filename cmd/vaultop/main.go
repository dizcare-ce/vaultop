package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/user/vaultop/internal/config"
)

func main() {
	cfgPath := flag.String("config", config.DefaultConfigFile, "path to vaultop config file")
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	cmd := flag.Arg(0)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	cfg.ApplyDefaults()

	switch cmd {
	case "validate":
		fmt.Printf("Config %q is valid (version=%s, providers=%d)\n",
			*cfgPath, cfg.Version, len(cfg.Providers))
	case "info":
		printInfo(cfg)
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", cmd)
		printUsage()
		os.Exit(1)
	}
}

func printInfo(cfg *config.Config) {
	fmt.Printf("vaultop config\n")
	fmt.Printf("  version:           %s\n", cfg.Version)
	fmt.Printf("  rotation_interval: %s\n", cfg.RotationDuration())
	fmt.Printf("  dry_run:           %v\n", cfg.Defaults.DryRun)
	fmt.Printf("  providers (%d):\n", len(cfg.Providers))
	for name, p := range cfg.Providers {
		fmt.Printf("    - %s (%s)\n", name, p.Type)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: vaultop [flags] <command>")
	fmt.Fprintln(os.Stderr, "Commands:")
	fmt.Fprintln(os.Stderr, "  validate   Validate the configuration file")
	fmt.Fprintln(os.Stderr, "  info       Print configuration summary")
	flag.PrintDefaults()
}
