package main

import (
	"fmt"
	"os"
)

func newFirstProject() {
	source := `
#include <iostream>

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}
`
	if _, err := os.Stat("src"); os.IsNotExist(err) {
		os.Mkdir("src", 0755)
	}

	err := os.WriteFile("src/main.cc", []byte(source), 0644)
	if err != nil {
		fmt.Printf("Error creating main.cc: %v\n", err)
		return
	}

	fmt.Println("Initialized project with main.cc")
}
