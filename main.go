package main

import (
	"bufio"
	"fmt"

	//"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	currentDirectory, _ := os.Getwd()
	storageDirectory := filepath.Join(currentDirectory, "storageData")

	// Create the storage directory if it doesn't exist
	if _, err := os.Stat(storageDirectory); os.IsNotExist(err) {
		err := os.Mkdir(storageDirectory, 0755)
		if err != nil {
			fmt.Printf("Failed to create the storage directory: %s\n", err.Error())
			return
		}
	}

	// Initialize virtual file system
	versioning, err := NewVersioning()
	baseDir := "./storageData" // Update the base directory path as per your requirement
	fs := NewFileSystem(baseDir, versioning)

	// Initialize database
	db := NewDatabase()

	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize cache
	cache := NewCache()

	// Main loop for user interaction
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Welcome to the Virtual File System!")
	fmt.Println("Enter 'help' to see available commands.")

	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			//fs.Close()
			db.Close()
			//cache.Close()
			versioning.Close()
			break
		}

		// Process user input
		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {
		case "help":
			printHelp()
		case "debug":
			fmt.Printf("** DEBUG **")
		case "pwd":
			handlePWDCommand(fs)
		case "cd":
			handleChangeDirectory(parts, fs)
		case "mkdir":
			handleCreateDirCommand(parts, storageDirectory)
		case "list":
			handleListCommand(parts, storageDirectory)
		case "rmdir":
			handleDeleteCommand(parts, storageDirectory)
		case "create":
			if len(parts) != 3 {
				fmt.Println("Invalid command. Usage: create <filename>")
				continue
			}
			filename := parts[1]
			err := fs.CreateFile(filename, []byte(parts[2]))
			if err != nil {
				fmt.Printf("Error creating file: %s\n", err.Error())
				continue
			}
			fmt.Println("File created successfully.")
		case "read":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: read <filename>")
				continue
			}
			filename := parts[1]
			content, err := fs.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err.Error())
				continue
			}
			fmt.Printf("File content: %s\n", content)
		case "update":
			if len(parts) != 3 {
				fmt.Println("Invalid command. Usage: update <filename> <content>")
				continue
			}
			filename := parts[1]
			content := []byte(parts[2])
			err := fs.UpdateFile(filename, content)
			if err != nil {
				fmt.Printf("Error updating file: %s\n", err.Error())
				continue
			}
			fmt.Println("File updated successfully.")
		case "delete":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: delete <filename>")
				continue
			}
			filename := parts[1]
			err := fs.DeleteFile(filename)
			if err != nil {
				fmt.Printf("Error deleting file: %s\n", err.Error())
				continue
			}
			fmt.Println("File deleted successfully.")
		case "compress":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: compress <filename>")
				continue
			}
			filename := parts[1]
			content, err := fs.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err.Error())
				continue
			}
			Compress(content, filename+".gz")
			if err != nil {
				fmt.Printf("Error compressing file: %s\n", err.Error())
				continue
			}
			fmt.Printf("Compressed content: " + filename + ".gz\n")
		case "decompress":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: decompress <filename>")
				continue
			}
			filename := parts[1]

			decompressedContent, err := Decompress(filename)
			if err != nil {
				fmt.Printf("Error decompressing file: %s\n", err.Error())
				continue
			}
			fmt.Printf("Decompressed content: %s\n", decompressedContent)
		case "encrypt":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: encrypt <filename>")
				continue
			}
			filename := parts[1]
			content, err := fs.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err.Error())
				continue
			}
			encryptedContent, err := Encrypt(content)
			if err != nil {
				fmt.Printf("Error encrypting file: %s\n", err.Error())
				continue
			}
			fmt.Printf("Encrypted content: %v\n", encryptedContent)
		case "decrypt":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: decrypt <filename>")
				continue
			}
			filename := parts[1]
			content, err := fs.ReadFile(filename)
			if err != nil {
				fmt.Printf("Error reading file: %s\n", err.Error())
				continue
			}
			decryptedContent, err := Decrypt(content)
			if err != nil {
				fmt.Printf("Error decrypting file: %s\n", err.Error())
				continue
			}
			fmt.Printf("Decrypted content: %s\n", decryptedContent)
		case "cache":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: cache <filename>")
				continue
			}
			filename := parts[1]
			content, err := cache.Get(filename)
			if err != false {
				//fmt.Printf("Error retrieving file from cache: %s\n", err.Error())
				continue
			}
			if content != nil {
				fmt.Printf("File content retrieved from cache: %s\n", content)
			} else {
				fileContent, err := fs.ReadFile(filename)
				if err != nil {
					fmt.Printf("Error reading file: %s\n", err.Error())
					continue
				}
				cache.Set(filename, fileContent)
				fmt.Printf("File content retrieved from file system: %s\n", fileContent)
			}
		case "version":
			if len(parts) != 2 {
				fmt.Println("Invalid command. Usage: version <filename>")
				continue
			}
			filename := parts[1]
			latestVersion, err := versioning.GetLatestVersion(filename)
			if err != nil {
				fmt.Printf("Error getting latest version: %s\n", err.Error())
				continue
			}
			fmt.Printf("Latest version of file '%s': %d\n", filename, latestVersion)

			// Retrieve all previous versions of the file
			previousVersions, err := versioning.GetAllVersions(filename)
			if err != nil {
				fmt.Printf("Error getting previous versions: %s\n", err.Error())
				continue
			}

			// Print the content of each previous version
			for _, version := range previousVersions {
				fmt.Printf("Version %d content: %s\n", version.Version, version.Content)
			}
		default:
			fmt.Println("Unknown command. Enter 'help' to see available commands.")
		}
	}
}

func parseInput(input string) []string {
	parts := []string{}
	for _, part := range filepath.SplitList(input) {
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func handlePWDCommand(fs *FileSystem) {
	fmt.Println("Current working directory:", fs.BaseDir)
}

func handleChangeDirectory(parts []string, fs *FileSystem) {
	if len(parts) != 2 {
		fmt.Println("Invalid command. Usage: cd <directory>")
		return
	}

	dirPath := parts[1]
	if dirPath == ".." {
		if fs.BaseDir != "./storageData" {
			// Navigate up one level within the base path
			fs.BaseDir = filepath.Dir(fs.BaseDir)
			fmt.Println("Changed to directory:", fs.BaseDir)
		} else {
			fmt.Println("Cannot navigate up beyond the base path.")
		}
		return
	}

	// Construct the new directory path
	newDirPath := filepath.Join(fs.BaseDir, dirPath)

	// Check if the directory exists
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		fmt.Println("Directory does not exist.")
		return
	}

	// Update the base directory to the new directory path
	fs.BaseDir = newDirPath
	fmt.Println("Changed to directory:", fs.BaseDir)
}

func handleCreateDirCommand(parts []string, storageDirectory string) {
	if len(parts) != 2 {
		fmt.Println("Invalid command. Usage: create <dirname>")
		return
	}

	dirname := parts[1]
	dirPath := filepath.Join(storageDirectory, dirname)
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err.Error())
		return
	}

	fmt.Println("Directory created successfully.")
}

func handleListCommand(parts []string, storageDirectory string) {
	if len(parts) != 1 {
		fmt.Println("Invalid command. Usage: list")
		return
	}

	filepath.Walk(storageDirectory, func(path string, info os.FileInfo, err error) error {
		if err == nil {
			relativePath, err := filepath.Rel(storageDirectory, path)
			if err != nil {
				fmt.Printf("Error getting relative path: %s\n", err.Error())
				return nil
			}

			// Check if it's a directory
			if info.IsDir() {
				fmt.Printf("[%s]\n", relativePath)
			} else {
				fmt.Println(relativePath)
			}
		}
		return nil
	})
}

func handleDeleteCommand(parts []string, storageDirectory string) {
	if len(parts) != 2 {
		fmt.Println("Invalid command. Usage: delete <filename>")
		return
	}

	filename := parts[1]
	filePath := filepath.Join(storageDirectory, filename)

	err := os.Remove(filePath)
	if err != nil {
		fmt.Printf("Error deleting file: %s\n", err.Error())
		return
	}

	fmt.Println("File deleted successfully.")
}

func printHelp() {
	fmt.Println("Available commands:")
	fmt.Println("help - Print this help message")
	fmt.Println("mkdir <dirname> - Create a new directory")
	fmt.Println("list - Lists all files and directories")
	fmt.Println("rmdir <dirname> - Delete a directory")
	fmt.Println("create <filename> <content> - Create a new file")
	fmt.Println("read <filename> - Read the content of a file")
	fmt.Println("update <filename> <content> - Update the content of a file")
	fmt.Println("delete <filename> - Delete a file")
	fmt.Println("compress <filename> - Compress the content of a file")
	fmt.Println("decompress <filename> - Decompress the content of a file")
	fmt.Println("encrypt <filename> - Encrypt the content of a file")
	fmt.Println("decrypt <filename> - Decrypt the content of a file")
	fmt.Println("cache <filename> - Get the content of a file from cache")
	fmt.Println("version <filename> - Get the latest version of a file")
	fmt.Println("exit - Exit the program")
}
