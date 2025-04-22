package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// initProject creates a basic CMakeLists.txt file for a new project
func initProject() {
	// Check if current directory is user's home directory
	homeDir, err := os.UserHomeDir()
	if err == nil && getCurrentDir() == homeDir {
		fmt.Println("Error: Cannot initialize a CMake project in your home directory.")
		fmt.Println("Please create a new directory for your project and run 'qs init' there.")
		return
	}

	// if CMakeLists.txt exists, return
	if fileExists("CMakeLists.txt") {
		fmt.Println("Error: CMakeLists.txt already exists. Run 'qs add' to add targets.")
		return
	}

	projectName := getProjectName()
	content := fmt.Sprintf(`cmake_minimum_required(VERSION 3.10)
project(%s)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Compiler options
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wall -Wextra")

# Output directories
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin)
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

# Include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include)

# Enable testing
enable_testing()
`, projectName)

	err = os.WriteFile("CMakeLists.txt", []byte(content), 0644)
	if err != nil {
		fmt.Printf("Error creating CMakeLists.txt: %v\n", err)
		return
	}

	fmt.Printf("Initialized CMake project '%s'\n", projectName)
	newFirstProject()
	addTarget(getProjectName(), []string{"src/main.cc"})
}

// getProjectName returns the project name for the project
func getProjectName() string {
	projectName := filepath.Base(getCurrentDir())
	return projectName
}

// addTarget adds an executable or library target to CMakeLists.txt
func addTarget(targetName string, sourceFiles []string) {
	// Determine if target exists in current directory
	cmakelists, err := os.ReadFile("CMakeLists.txt")
	if err != nil {
		fmt.Println("Error: CMakeLists.txt not found. Run 'qs init' first.")
		return
	}

	cmakeContent := string(cmakelists)

	// Process source files with glob support
	var expandedSourceFiles []string
	if len(sourceFiles) == 0 {
		// Check if target is a directory with main.cpp/main.c
		if isDir(targetName) {
			// Auto-add source files from directory
			expandedSourceFiles = findSourceFiles(targetName)
			if len(expandedSourceFiles) == 0 {
				fmt.Printf("No source files found in directory '%s'\n", targetName)
				return
			}
		} else {
			// Assume it's a single source file with same name as target
			// Try all common C/C++ file extensions
			extensions := []string{".cpp", ".cc", ".c", ".cxx"}
			mainFile := ""

			for _, ext := range extensions {
				candidateFile := targetName + ext
				if fileExists(candidateFile) {
					mainFile = candidateFile
					break
				}
			}

			if mainFile == "" {
				fmt.Printf("No source file found for target '%s' (tried .cpp, .cc, .c, .cxx extensions)\n", targetName)
				return
			}

			expandedSourceFiles = []string{mainFile}
		}
	} else {
		// Process each source file or glob pattern
		for _, pattern := range sourceFiles {
			// Check if pattern contains glob characters
			if containsGlobChar(pattern) {
				// Expand glob pattern
				matches, err := filepath.Glob(pattern)
				if err != nil {
					fmt.Printf("Invalid glob pattern '%s': %v\n", pattern, err)
					continue
				}
				if len(matches) == 0 {
					fmt.Printf("Warning: No files match pattern '%s'\n", pattern)
					continue
				}
				// Filter to only include source files
				for _, match := range matches {
					if isSourceFile(match) {
						expandedSourceFiles = append(expandedSourceFiles, match)
					}
				}
			} else if fileExists(pattern) {
				// Add the file directly
				expandedSourceFiles = append(expandedSourceFiles, pattern)
			} else if isDir(pattern) {
				// If it's a directory, find all source files in it
				dirFiles := findSourceFiles(pattern)
				expandedSourceFiles = append(expandedSourceFiles, dirFiles...)
			} else {
				fmt.Printf("Warning: File '%s' not found\n", pattern)
			}
		}
	}

	// Check if we have any source files after expansion
	if len(expandedSourceFiles) == 0 {
		fmt.Println("Error: No source files found for target")
		return
	}

	// Remove duplicates from expanded source files
	expandedSourceFiles = removeDuplicates(expandedSourceFiles)

	// Format source files list
	for i, file := range expandedSourceFiles {
		expandedSourceFiles[i] = strings.Replace(file, "\\", "/", -1)
	}

	// Check if target already exists
	targetPattern := fmt.Sprintf(`add_executable\(%s[\s\n]`, regexp.QuoteMeta(targetName))
	targetExists, _ := regexp.MatchString(targetPattern, cmakeContent)

	if targetExists {
		// Find the existing target block and append new source files to it
		re := regexp.MustCompile(fmt.Sprintf(`(add_executable\(%s\n)([^)]+)(\))`, regexp.QuoteMeta(targetName)))
		matches := re.FindStringSubmatch(cmakeContent)

		if len(matches) >= 3 {
			// Extract existing source files
			existingSourcesBlock := matches[2]
			existingSources := strings.Split(existingSourcesBlock, "\n")

			// Clean up the existing sources (remove whitespace)
			cleanedSources := []string{}
			for _, src := range existingSources {
				trimmed := strings.TrimSpace(src)
				if trimmed != "" {
					cleanedSources = append(cleanedSources, trimmed)
				}
			}

			// Append new sources, avoiding duplicates
			for _, newFile := range expandedSourceFiles {
				found := false
				for _, existingFile := range cleanedSources {
					if existingFile == newFile {
						found = true
						break
					}
				}
				if !found {
					cleanedSources = append(cleanedSources, newFile)
				}
			}

			// Create the updated source block
			newSourcesBlock := "    " + strings.Join(cleanedSources, "\n    ")

			// Replace the old target block with the new one
			replacement := fmt.Sprintf("${1}%s${3}", newSourcesBlock)
			cmakeContent = re.ReplaceAllString(cmakeContent, replacement)

			err = os.WriteFile("CMakeLists.txt", []byte(cmakeContent), 0644)
			if err != nil {
				fmt.Printf("Error updating CMakeLists.txt: %v\n", err)
				return
			}

			fmt.Printf("Updated existing target '%s' with %d additional source files\n",
				targetName, len(expandedSourceFiles))
			return
		}
	}

	// Create the target addition text (for new targets)
	targetAddition := fmt.Sprintf("\nadd_executable(%s\n    %s\n)\n",
		targetName,
		strings.Join(expandedSourceFiles, "\n    "))

	// Append the target to CMakeLists.txt
	cmakeContent += targetAddition

	err = os.WriteFile("CMakeLists.txt", []byte(cmakeContent), 0644)
	if err != nil {
		fmt.Printf("Error updating CMakeLists.txt: %v\n", err)
		return
	}

	fmt.Printf("Added executable target '%s' with %d source files\n", targetName, len(expandedSourceFiles))
}

// addStandardConfig adds standard CMake configuration
func addStandardConfig(cxxStd int) {
	// Check if CMakeLists.txt exists
	if !fileExists("CMakeLists.txt") {
		fmt.Println("Error: CMakeLists.txt not found. Run 'qs init' first.")
		return
	}

	cmakelists, err := os.ReadFile("CMakeLists.txt")
	if err != nil {
		fmt.Printf("Error reading CMakeLists.txt: %v\n", err)
		return
	}

	cmakeContent := string(cmakelists)

	// Check if standard configuration is already added
	stdConfigAdded := strings.Contains(cmakeContent, "CMAKE_RUNTIME_OUTPUT_DIRECTORY") ||
		strings.Contains(cmakeContent, "CMAKE_ARCHIVE_OUTPUT_DIRECTORY")

	// If user provided a C++ standard, update it in the file
	if cxxStd > 0 {
		// Use regex to replace existing CMAKE_CXX_STANDARD line
		re := regexp.MustCompile(`set\(CMAKE_CXX_STANDARD \d+\)`)
		if re.MatchString(cmakeContent) {
			// Replace existing CMAKE_CXX_STANDARD value
			cmakeContent = re.ReplaceAllString(cmakeContent, fmt.Sprintf("set(CMAKE_CXX_STANDARD %d)", cxxStd))
			fmt.Printf("Updated C++ standard to C++%d\n", cxxStd)
		} else {
			// Add CMAKE_CXX_STANDARD to config
			stdConfig := fmt.Sprintf("\n# C++ Standard\nset(CMAKE_CXX_STANDARD %d)\nset(CMAKE_CXX_STANDARD_REQUIRED ON)\n", cxxStd)
			cmakeContent += stdConfig
			fmt.Printf("Set C++ standard to C++%d\n", cxxStd)
		}
	}

	// Add standard configurations if not already present
	if !stdConfigAdded {
		// First, try to find target names for install command
		re := regexp.MustCompile(`add_executable\(([^):\n\s]+)`)
		matches := re.FindAllStringSubmatch(cmakeContent, -1)

		targetNames := []string{}
		for _, match := range matches {
			if len(match) > 1 {
				targetNames = append(targetNames, match[1])
			}
		}

		installCmd := "# Add install target\n"
		if len(targetNames) > 0 {
			installCmd += fmt.Sprintf("install(TARGETS %s DESTINATION bin)\n", strings.Join(targetNames, " "))
		} else {
			installCmd += "# No targets found to install\n"
		}

		stdConfig := fmt.Sprintf(`

%s`, installCmd)

		cmakeContent += stdConfig
		fmt.Println("Added standard CMake configuration")
	} else {
		// Standard config already exists
		fmt.Println("Standard CMake configuration already present")
	}

	err = os.WriteFile("CMakeLists.txt", []byte(cmakeContent), 0644)
	if err != nil {
		fmt.Printf("Error updating CMakeLists.txt: %v\n", err)
		return
	}
}

// Helper functions
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return ""
	}
	return dir
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

func isSourceFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".cpp" || ext == ".c" || ext == ".cc" || ext == ".cxx" || ext == ".h" || ext == ".hpp" || ext == ".hxx"
}

func containsGlobChar(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

func findSourceFiles(dirPath string) []string {
	var sourceFiles []string

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return sourceFiles
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if isSourceFile(name) {
			sourceFiles = append(sourceFiles, filepath.Join(dirPath, name))
		}
	}

	return sourceFiles
}

func removeDuplicates(files []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, file := range files {
		if !seen[file] {
			seen[file] = true
			result = append(result, file)
		}
	}

	return result
}
