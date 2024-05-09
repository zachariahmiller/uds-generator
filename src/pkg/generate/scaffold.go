package generate

import (
	"embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/uds-generator/src/pkg/common"
	git "github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
)

//go:embed all:repo/*
var repoTemplate embed.FS

func Scaffold(data common.RepoInfo) (string, error) {
	var folder = "uds-package-" + data.Name
	var fullPath = config.GenerateOutputDir + "/" + folder
	log.Println("Starting Scaffold function")

	// Create generated directory if it doesn't exist
	if err := os.MkdirAll(config.GenerateOutputDir, 0755); err != nil {
		log.Println("Failed to create output directory:", err)
		return config.GenerateOutputDir, err
	}

	// Copy template repo
	writeEmbeddedFolder(repoTemplate, "repo/", "uds-package-"+data.Name)

	// List of Markdown file paths
	filePaths := []string{"CONTRIBUTING.md", "README.md", "SECURITY.md", "CODEOWNERS", "chart/Chart.yaml", "renovate.json"}

	// Map of strings to find and their replacements
	replacements := map[string]string{
		"<name>":         strings.ToLower(data.Name),
		"<name-upper>":   SimpleTitleCase(data.Name),
		"<organization>": data.Organization,
		"<author>":       data.Author,
		"<year>":         strconv.Itoa(time.Now().Year()),
	}

	// Loop over each file path
	for _, filePath := range filePaths {
		fmt.Println("Updating file:", filePath)
		err := ReplaceInFile(fullPath+"/"+filePath, replacements)
		if err != nil {
			fmt.Printf("An error occurred while updating file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("File %s updated successfully\n", filePath)
		}
	}
	// Initialize the repository
	repo, err := InitGitRepo(fullPath)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}

	// Adding a remote to the repository
	remoteName := "origin"
	remoteURL := "https://" + data.Platform + "/" + data.Organization + "/uds-package-" + data.Name + ".git"
	err = AddRemote(repo, remoteName, remoteURL)
	if err != nil {
		fmt.Println("Error:", err)
	}

	log.Println("Scaffold function completed successfully")
	return remoteURL, nil
}

func ReplaceInFile(filePath string, replacements map[string]string) error {

	input, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	content := string(input)

	for find, replace := range replacements {
		content = strings.Replace(content, find, replace, -1) // Replace all occurrences
	}

	err = os.WriteFile(filePath, []byte(content), 0666)
	if err != nil {
		return err
	}
	return nil
}

func SimpleTitleCase(input string) string {
	words := strings.Fields(input)
	for i, word := range words {
		words[i] = strings.ToUpper(string(word[0])) + word[1:]
	}
	return strings.Join(words, " ")
}

// InitGitRepo initializes a new git repository at the specified path.
func InitGitRepo(repoPath string) (*git.Repository, error) {
	// Initialize non-bare repository
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to init git repository: %v", err)
	}
	fmt.Println("Repository initialized successfully at", repoPath)
	return repo, nil
}

// AddRemote adds a new remote to the given repository.
func AddRemote(repo *git.Repository, remoteName, remoteURL string) error {
	_, err := repo.CreateRemote(&gitConfig.RemoteConfig{
		Name: remoteName,
		URLs: []string{remoteURL},
	})
	if err != nil {
		return fmt.Errorf("failed to add remote: %v", err)
	}
	fmt.Println("Remote added successfully:", remoteName, remoteURL)
	return nil
}
