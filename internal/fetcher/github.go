package fetcher

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Fetch downloads or locates the Go package source code.
// For standard library packages, it uses GOROOT.
// For remote packages, it clones the git repository.
// knownRepos maps vanity import paths to actual git URLs.
func Fetch(packagePath string, knownRepos map[string]string) (string, error) {
	// Check if it's a standard library package
	if isStdLib(packagePath) {
		return fetchStdLib(packagePath)
	}

	// Check if it's a known vanity import
	if repoURL, ok := findKnownRepo(packagePath, knownRepos); ok {
		return fetchKnownRepo(packagePath, repoURL)
	}

	// For remote packages, clone the repository
	return fetchRemote(packagePath)
}

// findKnownRepo checks if the package matches a known vanity import path
func findKnownRepo(packagePath string, knownRepos map[string]string) (string, bool) {
	// Exact match first
	if url, ok := knownRepos[packagePath]; ok {
		return url, true
	}
	// Check if packagePath is a sub-package of a known repo
	for prefix, url := range knownRepos {
		if strings.HasPrefix(packagePath, prefix+"/") {
			return url, true
		}
	}
	return "", false
}

// fetchKnownRepo clones a known vanity import repository
func fetchKnownRepo(packagePath string, repoURL string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "godoceasy-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	cloneDir := filepath.Join(tmpDir, "repo")

	fmt.Printf("   Cloning %s (resolved from %s)...\n", repoURL, packagePath)

	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to clone %s: %w", repoURL, err)
	}

	// Check and display the latest version
	showLatestVersion(cloneDir, repoURL)

	return cloneDir, nil
}

// isStdLib checks if a package is part of the Go standard library
func isStdLib(pkg string) bool {
	// Standard library packages don't contain a dot in the first path element
	parts := strings.SplitN(pkg, "/", 2)
	return !strings.Contains(parts[0], ".")
}

// fetchStdLib locates a standard library package source
func fetchStdLib(pkg string) (string, error) {
	goroot := os.Getenv("GOROOT")
	if goroot == "" {
		// Try to detect GOROOT
		out, err := exec.Command("go", "env", "GOROOT").Output()
		if err != nil {
			return "", fmt.Errorf("cannot determine GOROOT: %w", err)
		}
		goroot = strings.TrimSpace(string(out))
	}

	srcDir := filepath.Join(goroot, "src", pkg)
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return "", fmt.Errorf("standard library package not found: %s", pkg)
	}

	return srcDir, nil
}

// fetchRemote clones a remote Git repository
func fetchRemote(packagePath string) (string, error) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "godoceasy-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}

	// Determine git URL from package path
	repoURL, subDir := resolveRepoURL(packagePath)

	cloneDir := filepath.Join(tmpDir, "repo")

	fmt.Printf("   Cloning %s...\n", repoURL)

	// Use git clone (shallow clone for speed)
	cmd := exec.Command("git", "clone", "--depth", "1", repoURL, cloneDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// Fallback: try with go mod download
		return fetchWithGoModDownload(packagePath, tmpDir)
	}

	// Check and display the latest version
	showLatestVersion(cloneDir, repoURL)

	// If there's a subdirectory, use that
	if subDir != "" {
		targetDir := filepath.Join(cloneDir, subDir)
		if _, err := os.Stat(targetDir); err == nil {
			return targetDir, nil
		}
	}

	return cloneDir, nil
}

// resolveRepoURL converts a Go import path to a git clone URL
func resolveRepoURL(packagePath string) (string, string) {
	parts := strings.Split(packagePath, "/")

	if len(parts) < 3 {
		return "https://" + packagePath + ".git", ""
	}

	// For GitHub/GitLab/Bitbucket, repo is the first 3 parts
	host := parts[0]
	switch {
	case strings.Contains(host, "github.com"),
		strings.Contains(host, "gitlab.com"),
		strings.Contains(host, "bitbucket.org"):
		repoPath := strings.Join(parts[:3], "/")
		subDir := ""
		if len(parts) > 3 {
			subDir = strings.Join(parts[3:], "/")
		}
		return "https://" + repoPath + ".git", subDir
	default:
		return "https://" + packagePath + ".git", ""
	}
}

// fetchWithGoModDownload uses 'go mod download' as a fallback
func fetchWithGoModDownload(packagePath string, tmpDir string) (string, error) {
	fmt.Println("   Trying go mod download as fallback...")

	// Create a temporary Go module
	modDir := filepath.Join(tmpDir, "gomod")
	if err := os.MkdirAll(modDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create mod dir: %w", err)
	}

	// Create go.mod
	goMod := fmt.Sprintf(`module godoceasy-fetch

go 1.22.0

require %s latest
`, packagePath)
	if err := os.WriteFile(filepath.Join(modDir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", fmt.Errorf("failed to write go.mod: %w", err)
	}

	// Create a dummy main.go
	mainGo := fmt.Sprintf(`package main

import _ "%s"

func main() {}
`, packagePath)
	if err := os.WriteFile(filepath.Join(modDir, "main.go"), []byte(mainGo), 0644); err != nil {
		return "", fmt.Errorf("failed to write main.go: %w", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = modDir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go mod tidy failed: %s: %w", string(out), err)
	}

	// Find the module cache path
	cmd = exec.Command("go", "list", "-m", "-json", packagePath)
	cmd.Dir = modDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go list failed: %s: %w", string(out), err)
	}

	// Parse the directory from output
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, `"Dir"`) {
			dir := strings.TrimPrefix(line, `"Dir": "`)
			dir = strings.TrimSuffix(dir, `",`)
			dir = strings.TrimSuffix(dir, `"`)
			return dir, nil
		}
	}

	return modDir, nil
}

// showLatestVersion checks and displays the latest version tag of the cloned repo
func showLatestVersion(cloneDir string, repoURL string) {
	// Try to get latest tag from remote
	cmd := exec.Command("git", "ls-remote", "--tags", "--sort=-v:refname", repoURL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Fallback: try local tags
		cmd = exec.Command("git", "describe", "--tags", "--abbrev=0")
		cmd.Dir = cloneDir
		out, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Println("   ℹ️  Version: latest (main branch)")
			return
		}
		tag := strings.TrimSpace(string(out))
		fmt.Printf("   📌 Version: %s (latest tag)\n", tag)
		return
	}

	// Parse the latest tag from ls-remote output
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	latestTag := ""
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		if len(parts) < 2 {
			continue
		}
		ref := parts[1]
		// Skip ^{} dereferences
		if strings.HasSuffix(ref, "^{}") {
			continue
		}
		tag := strings.TrimPrefix(ref, "refs/tags/")
		// Prefer semver tags (v1.x.x)
		if strings.HasPrefix(tag, "v") {
			latestTag = tag
			break
		}
		if latestTag == "" {
			latestTag = tag
		}
	}

	if latestTag != "" {
		fmt.Printf("   📌 Latest version: %s\n", latestTag)
		fmt.Printf("   ✅ You are viewing the latest code (main branch)\n")
	} else {
		fmt.Println("   ℹ️  Version: latest (main branch)")
	}
}
