package main

import (
	"fmt"
	"os"
	"strconv"
)

const version = "0.1.0"

func printHelp() {
	fmt.Println("qs - Quick Setup for CMake projects")
	fmt.Println("Usage:")
	fmt.Println("  qs init                   Initialize a new CMake project")
	fmt.Println("  qs add <target> [files]   Add executable or library target")
	fmt.Println("                            [files] can include glob patterns like *.cpp")
	fmt.Println("  qs std [cxx_std]          Add standard CMake configuration with optional C++ standard (11/14/17/20)")
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
	case "version":
		fmt.Printf("qs version %s\n", version)
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
	}
}
