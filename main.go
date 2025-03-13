package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	folderName := ".rev"
	exeFileName := "file-rewind.exe"
	var changes string
	reader := bufio.NewReader(os.Stdin)
	now := time.Now()

	date := now.Format(time.DateOnly)
	hashTime := now.Format(time.Stamp)
	time := now.Format(time.Kitchen)

	formattedTime := strings.Replace(time, ":", ".", -1)
	hash := md5.Sum([]byte(hashTime)) // MD5 hashs
	savedFolderName := date + " (" + formattedTime + ") " + hex.EncodeToString(hash[:])

	sourceDir := "."
	destDir := folderName + "/" + savedFolderName

	// File to exclude
	excludeFolder := folderName

	// Check if folder exists
	_, err := os.Stat(folderName)
	if err != nil {
		if os.IsNotExist(err) {
			// Folder does not exist
			err = os.Mkdir(folderName, 0755)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)

				return
			}
			err = os.Mkdir(destDir, 0755)
			if err != nil {
				fmt.Printf("Error creating directory: %v\n", err)

				return
			}

			fmt.Printf("Directory '%s' created successfully\n", folderName)
		}

	} else {
		err = os.Mkdir(destDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)

			return
		}
	}

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the destination directory to avoid recursion
		if info.IsDir() && path == destDir {
			fmt.Printf("Skipping destination directory: %s\n", path)
			return filepath.SkipDir
		}

		// Skip the excluded folder and its contents
		if info.IsDir() && filepath.Base(path) == excludeFolder {
			fmt.Printf("Skipping excluded directory: %s\n", path)
			return filepath.SkipDir
		}

		if info.Name() == exeFileName {
			// Skip exe
			return nil
		}

		// Skip the root directory itself
		if path == sourceDir {
			return nil
		}

		// Calculate the relative path from source directory
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Create the destination path
		destPath := filepath.Join(destDir, relPath)

		// If it's a directory, create it in the destination
		if info.IsDir() {
			fmt.Printf("Creating directory: %s\n", destPath)
			return os.MkdirAll(destPath, info.Mode())
		}

		// Otherwise, copy the file
		fmt.Printf("Saving file: %s\n", relPath)
		return copyFile(path, destPath)
	})

	if err != nil {
		fmt.Printf("Error copying files: %v\n", err)
		return
	}

	fmt.Printf("All files and directories copied from current directory to '%s' (excluding '%s' folder)\n", destDir, excludeFolder)
	fmt.Print("Declare file changes for this revision (optional): ")
	changes, _ = reader.ReadString('\n')
	if changes != "" {
		err := os.WriteFile(destDir+"/changes.txt", []byte(changes), 0644) // 0644 = Read & Write permission
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
		fmt.Println("File written successfully!")
	}
	// For debug purposes
	// bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func copyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy the contents
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Sync to ensure write is complete
	err = destFile.Sync()
	if err != nil {
		return err
	}

	// Get and set file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, sourceInfo.Mode())
}
