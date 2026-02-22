package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/bmaca/go-version-manager/internal/dl"
	"github.com/bmaca/go-version-manager/internal/install"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		installCmd(os.Args[2:])
	default:
		usage()
		os.Exit(1)
	}
}

func getOsArch() (string, string) {
	return runtime.GOOS, runtime.GOARCH
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  gotool install --version <version> --dir <install-dir> [--force]")
}

func installCmd(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	version := fs.String("version", "", "Go version to install (e.g. 1.22.0)")
	dir := fs.String("dir", "", "Install root directory (e.g. /opt/go)")
	force := fs.Bool("force", false, "Reinstall if already exists without prompting")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *version == "" || *dir == "" {
		fs.Usage()
		os.Exit(1)
	}
	if err := runInstall(*version, *dir, *force); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runInstall(version, dir string, force bool) error {
	goOS, goArch := getOsArch()

	rel, file, dlUrl, err := dl.Resolve(version, goOS, goArch)
	if err != nil {
		return fmt.Errorf("resolve version: %w", err)
	}

	targetDir := install.TargetDir(dir, rel.Version)
	if install.Exists(targetDir) {
		if !force {
			ok, err := install.PromptYesNo(
				fmt.Sprintf("Go %s already exists at %s. Reinstall?", rel.Version, targetDir),
			)
			if err != nil {
				return err
			}
			if !ok {
				fmt.Println("Aborted.")
				return nil
			}
		}
		if err := install.RemoveDir(targetDir); err != nil {
			return fmt.Errorf("remove existing: %w", err)
		}
	}

	tmpFile, err := install.TempFile()
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile)

	if err := install.Download(dlUrl, tmpFile); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	if err := install.VerifyChecksum(tmpFile, file.SHA256); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}
	fmt.Println("Checksum verified")

	if err := install.ExtractTarGz(tmpFile, targetDir); err != nil {
		return fmt.Errorf("extract: %w", err)
	}

	fmt.Printf("Installed Go %s to %s\n", rel.Version, targetDir)
	return nil
}
