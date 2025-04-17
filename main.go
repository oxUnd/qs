package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const version = "0.1.0"

func printHelp() {
	fmt.Println("qs - Quick Setup for CMake projects")
	fmt.Println("Usage:")
	fmt.Println("  qs init                   Initialize a new CMake project")
	fmt.Println("  qs add <target> [files]   Add executable or library target")
	fmt.Println("                            [files] can include glob patterns like *.cpp")
	fmt.Println("  qs std [cxx_std]          Add standard CMake configuration with optional C++ standard (11/14/17/20)")
	fmt.Println("  qs build                  Create build directory, run cmake and make")
	fmt.Println("  qs run [target]           Run the specified executable target (or default target if not specified)")
	fmt.Println("  qs version                Show version information")
	fmt.Println("  qs help                   Show this help message")
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "init":
		initProject()
		newFirstProject()
		addTarget(getTargetName(), []string{"src/main.cc"})

	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Error: 'add' requires a target name")
			return
		}
		targetName := os.Args[2]
		var sourceFiles []string
		if len(os.Args) > 3 {
			sourceFiles = os.Args[3:]
		}
		addTarget(targetName, sourceFiles)
	case "std":
		cxxStd := 0
		if len(os.Args) > 2 {
			val, err := strconv.Atoi(os.Args[2])
			if err == nil && (val == 11 || val == 14 || val == 17 || val == 20) {
				cxxStd = val
			} else {
				fmt.Println("Warning: Invalid C++ standard, using default")
			}
		}
		addStandardConfig(cxxStd)
	case "build":
		buildProject()
	case "run":
		targetName := ""
		if len(os.Args) > 2 {
			targetName = os.Args[2]
		}
		runProject(targetName)
	case "version":
		fmt.Printf("qs version %s\n", version)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}

// buildProject creates a build directory, runs cmake and make
func buildProject() {
	// Check for CMakeLists.txt
	if _, err := os.Stat("CMakeLists.txt"); os.IsNotExist(err) {
		fmt.Println("Error: CMakeLists.txt not found in the current directory.")
		fmt.Println("Run 'qs init' to create a new CMake project.")
		return
	}

	// Create build directory if it doesn't exist
	if _, err := os.Stat("build"); os.IsNotExist(err) {
		fmt.Println("Creating build directory...")
		err := os.Mkdir("build", 0755)
		if err != nil {
			fmt.Printf("Error creating build directory: %s\n", err)
			return
		}
	}

	// Change to build directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %s\n", err)
		return
	}

	buildDir := cwd + "/build"

	// Run cmake
	fmt.Println("Running CMake...")
	cmakeCmd := exec.Command("cmake", "..")
	cmakeCmd.Dir = buildDir
	cmakeCmd.Stdout = os.Stdout
	cmakeCmd.Stderr = os.Stderr

	err = cmakeCmd.Run()
	if err != nil {
		fmt.Printf("Error running cmake: %s\n", err)
		return
	}

	// Run make
	fmt.Println("Running make...")
	makeCmd := exec.Command("make")
	makeCmd.Dir = buildDir
	makeCmd.Stdout = os.Stdout
	makeCmd.Stderr = os.Stderr

	err = makeCmd.Run()
	if err != nil {
		fmt.Printf("Error running make: %s\n", err)
		return
	}

	fmt.Println("Build completed successfully!")
}

// runProject runs a built executable target from the build directory
func runProject(targetName string) {
	// Check for build directory
	if _, err := os.Stat("build"); os.IsNotExist(err) {
		fmt.Println("Error: build directory not found.")
		fmt.Println("Run 'qs build' to build the project first.")
		return
	}

	// Get current directory to construct build path
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %s\n", err)
		return
	}

	// Check if bin directory exists (standard layout)
	binDir := cwd + "/build/bin"
	executablesPath := binDir
	if _, err := os.Stat(binDir); os.IsNotExist(err) {
		// Fall back to just the build directory
		executablesPath = cwd + "/build"
	}

	// If no target specified, try to find one
	if targetName == "" {
		// Try to find an executable in the build directory
		files, err := os.ReadDir(executablesPath)
		if err != nil {
			fmt.Printf("Error reading build directory: %s\n", err)
			return
		}

		// Look for executable files
		var executables []string
		for _, file := range files {
			// Skip directories and files starting with "."
			if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}

			// On Unix systems, check if file is executable
			filePath := filepath.Join(executablesPath, file.Name())
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			// Check if file has execute permission
			if fileInfo.Mode()&0111 != 0 {
				executables = append(executables, file.Name())
			}
		}

		if len(executables) == 0 {
			fmt.Println("Error: No executable targets found in build directory.")
			fmt.Println("Specify a target name or build the project first with 'qs build'.")
			return
		} else if len(executables) == 1 {
			targetName = executables[0]
			fmt.Printf("Running target: %s\n", targetName)
		} else {
			// Multiple executables found, let user choose
			fmt.Println("Multiple targets found:")
			for i, exe := range executables {
				fmt.Printf("  %d. %s\n", i+1, exe)
			}
			fmt.Println("Please specify a target name: qs run <target>")
			return
		}
	}

	// Construct path to the executable
	targetPath := filepath.Join(executablesPath, targetName)
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		fmt.Printf("Error: Target '%s' not found in build directory.\n", targetName)
		return
	}

	// Run the executable
	fmt.Printf("Running %s...\n", targetName)
	cmd := exec.Command(targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error running target: %s\n", err)
		return
	}
}
