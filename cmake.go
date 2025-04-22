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

// initSubProject creates a subdirectory with a CMakeLists.txt file for a sub-project
func initSubProject(subDirName string) {
	// Validate subdirectory name
	if strings.TrimSpace(subDirName) == "" {
		fmt.Println("Error: Please specify a valid subdirectory name.")
		return
	}

	// Check if parent CMakeLists.txt exists
	if !fileExists("CMakeLists.txt") {
		fmt.Println("Error: Main CMakeLists.txt not found in the current directory.")
		fmt.Println("Run 'qs init' in the parent directory before creating a sub-project.")
		return
	}

	// Create the subdirectory if it doesn't exist
	if _, err := os.Stat(subDirName); os.IsNotExist(err) {
		err := os.MkdirAll(subDirName, 0755)
		if err != nil {
			fmt.Printf("Error creating subdirectory '%s': %v\n", subDirName, err)
			return
		}
		fmt.Printf("Created subdirectory '%s'\n", subDirName)
	}

	// Create the CMakeLists.txt for the subdirectory
	subCMakeContent := fmt.Sprintf(`# %s sub-project

# Add include directories
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/include)

# Add source files
file(GLOB SOURCES "*.cpp" "*.cc" "*.c" "src/*.cpp" "src/*.cc" "src/*.c")
file(GLOB HEADERS "*.h" "*.hpp" "include/*.h" "include/*.hpp")

# Add library
add_library(%s ${SOURCES} ${HEADERS})

# Link any dependencies if needed
# target_link_libraries(%s PRIVATE dependency1 dependency2)

# Set include directories for targets that link this library
target_include_directories(%s PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<INSTALL_INTERFACE:include>
)

# Install rules
install(TARGETS %s
    ARCHIVE DESTINATION lib
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
)
install(DIRECTORY include/ DESTINATION include)
`, subDirName, subDirName, subDirName, subDirName, subDirName)

	subCMakePath := filepath.Join(subDirName, "CMakeLists.txt")
	err := os.WriteFile(subCMakePath, []byte(subCMakeContent), 0644)
	if err != nil {
		fmt.Printf("Error creating CMakeLists.txt in '%s': %v\n", subDirName, err)
		return
	}

	// Create include directory in the subdirectory
	includeDir := filepath.Join(subDirName, "include")
	if _, err := os.Stat(includeDir); os.IsNotExist(err) {
		err := os.MkdirAll(includeDir, 0755)
		if err != nil {
			fmt.Printf("Error creating include directory: %v\n", err)
		}
	}

	// Create src directory in the subdirectory
	srcDir := filepath.Join(subDirName, "src")
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		err := os.MkdirAll(srcDir, 0755)
		if err != nil {
			fmt.Printf("Error creating src directory: %v\n", err)
		}
	}

	// Update parent CMakeLists.txt to include the subdirectory
	parentCMake, err := os.ReadFile("CMakeLists.txt")
	if err != nil {
		fmt.Printf("Error reading parent CMakeLists.txt: %v\n", err)
		return
	}

	parentCMakeContent := string(parentCMake)

	// Check if the subdirectory is already included
	addSubdirLine := fmt.Sprintf("add_subdirectory(%s)", subDirName)
	if strings.Contains(parentCMakeContent, addSubdirLine) {
		fmt.Printf("Subdirectory '%s' is already included in the parent CMakeLists.txt\n", subDirName)
	} else {
		// Add the subdirectory to the parent CMakeLists.txt
		parentCMakeContent += fmt.Sprintf("\n# Include the %s sub-project\nadd_subdirectory(%s)\n", subDirName, subDirName)

		err = os.WriteFile("CMakeLists.txt", []byte(parentCMakeContent), 0644)
		if err != nil {
			fmt.Printf("Error updating parent CMakeLists.txt: %v\n", err)
			return
		}
		fmt.Printf("Added subdirectory '%s' to parent CMakeLists.txt\n", subDirName)
	}

	linkSubLibLine := fmt.Sprintf("target_link_libraries(%s PRIVATE %s)", getProjectName(), subDirName)
	if strings.Contains(parentCMakeContent, linkSubLibLine) {
		fmt.Printf("Subdirectory '%s' is already linked in the parent CMakeLists.txt\n", subDirName)
	} else {
		parentCMakeContent += fmt.Sprintf("\n# Link the %s sub-project\ntarget_link_libraries(%s PRIVATE %s)\n", subDirName, getProjectName(), subDirName)

		err = os.WriteFile("CMakeLists.txt", []byte(parentCMakeContent), 0644)
		if err != nil {
			fmt.Printf("Error updating parent CMakeLists.txt: %v\n", err)
			return
		}
		fmt.Printf("Added subdirectory '%s' to parent CMakeLists.txt\n", subDirName)
	}

	// Create a sample header file
	headerContent := fmt.Sprintf(`#pragma once

namespace %s {

class Example {
public:
    Example();
    void doSomething();
};

} // namespace %s
`, subDirName, subDirName)

	headerPath := filepath.Join(includeDir, subDirName+".h")
	err = os.WriteFile(headerPath, []byte(headerContent), 0644)
	if err != nil {
		fmt.Printf("Error creating sample header: %v\n", err)
	}

	// Create a sample source file
	sourceContent := fmt.Sprintf(`#include "%s.h"
#include <iostream>

namespace %s {

Example::Example() {
    // Constructor implementation
}

void Example::doSomething() {
    std::cout << "Hello from %s library!" << std::endl;
}

} // namespace %s
`, subDirName, subDirName, subDirName, subDirName)

	sourcePath := filepath.Join(srcDir, subDirName+".cc")
	err = os.WriteFile(sourcePath, []byte(sourceContent), 0644)
	if err != nil {
		fmt.Printf("Error creating sample source: %v\n", err)
	}

	fmt.Printf("Successfully initialized sub-project '%s'.\n", subDirName)
	fmt.Printf("To link this library to an executable, use: target_link_libraries(your_executable PRIVATE %s)\n", subDirName)
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
