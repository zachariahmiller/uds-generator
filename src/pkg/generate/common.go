package generate

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/defenseunicorns/uds-generator/src/config"
)

// Write an embedded folder to the localDir
func writeFolder(folder embed.FS) {
	// Walk through the embedded filesystem
	err := fs.WalkDir(folder, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Construct destination file path
		destPath := filepath.Join(config.GenerateOutputDir, path)

		if d.IsDir() {
			// Create directory if it doesn't exist
			return os.MkdirAll(destPath, 0755)
		}

		// Open source file from embedded filesystem
		srcFile, err := folder.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		// Create destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Copy contents from source to destination
		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}
}
