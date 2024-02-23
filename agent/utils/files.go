package utils

import (
	"io"
	"os"
	"path/filepath"
)

func CopyFolder(src, dest string) error {
	// Read the source folder
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Create the destination folder if it doesn't exist
	if err := os.MkdirAll(dest, 0777); err != nil {
		return err
	}

	// Copy each file from source to destination
	for _, file := range files {
		srcPath := filepath.Join(src, file.Name())
		destPath := filepath.Join(dest, file.Name())

		if file.IsDir() {
			// Recursively copy subdirectories
			if err := CopyFolder(srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy regular files
			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Get source file permissions
	srcFileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Set destination file permissions
	if err := os.Chmod(dest, srcFileInfo.Mode()); err != nil {
		return err
	}

	return nil
}