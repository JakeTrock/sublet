package main

import (
	"fmt"
	"os"
	"os/exec"
)

const NIXDIR = "/etc/nixos"

/**
 * init configuration.nix
 */
func initConfiguration() error {
	// initialize a git repo in nixos directory
	cmd := exec.Command("git", "init")
	cmd.Dir = NIXDIR
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Failed to initialize git repo: %v", err)
	}

	return nil
}

/**
 * list nix files
 */
func listNixFiles() ([]string, error) {
	files, err := os.ReadDir(NIXDIR)
	if err != nil {
		return nil, fmt.Errorf("Failed to read directory: %v", err)
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

/**
 * get contents of all nix files and print as gpt-compatible markdown
 */
func getNixFilesContents() (string, error) {
	files, err := listNixFiles()
	if err != nil {
		return "", err
	}

	var contents string
	for _, file := range files {
		content, err := os.ReadFile(fmt.Sprintf("%s/%s", NIXDIR, file))
		if err != nil {
			return "", fmt.Errorf("Failed to read file: %v", err)
		}
		contents += fmt.Sprintf("### %s\n\n```nix\n%s\n```\n\n", file, string(content))
	}
	return contents, nil
}

/**
 * set nix file contents, key value pairs of filename and content
 */
func setNixFilesContents(files map[string]string) error {
	for file, content := range files {
		err := os.WriteFile(fmt.Sprintf("%s/%s", NIXDIR, file), []byte(content), 0644)
		if err != nil {
			return fmt.Errorf("Failed to write file: %v", err)
		}
	}
	return nil
}

/**
 * Run a dry build
 * This command builds the new configuration but does not switch to it.
 * It allows you to see what would be built and changed without actually applying the changes.
 */
func runDryBuild() (string, error) {
	cmd := exec.Command("nixos-rebuild", "dry-build")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running dry-build: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

/**
 * Run a test
 * This command builds the new configuration and switches to it,
 * but the changes are not persistent across reboots.
 * It is useful for testing new configurations temporarily.
 * If you reboot the system, it will revert to the previous configuration
 */
func runTest() (string, error) {
	cmd := exec.Command("nixos-rebuild", "test")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running test: %v, output: %s", err, string(output))
	}
	return string(output), nil
}

/**
 * Run a switch
 * This command builds the new configuration and then switches the system to use it.
 * It applies the changes immediately, restarting services as necessary.
 */
func runSwitch() (string, error) {
	/*
		// add files to git now that build succeeded
		cmd := exec.Command("git", "add", ".")
		cmd.Dir = NIXDIR
		err := cmd.Run()
		if err != nil {
			return "", fmt.Errorf("Failed to add files to git: %v", err)
		}
		// commit with timestamp
		cmd = exec.Command("git", "commit", "-m", "Configuration updated at"+time.Now().String())
		cmd.Dir = NIXDIR
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("Failed to commit changes: %v", err)
		}
		// push to remote
		cmd = exec.Command("git", "push")
		cmd.Dir = NIXDIR
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("Failed to push changes: %v", err)
		}
	*/
	cmd := exec.Command("nixos-rebuild", "switch")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error running switch: %v, output: %s", err, string(output))
	}
	return string(output), nil
}
