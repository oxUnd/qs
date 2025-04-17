# QS - Quick Setup for CMake

A simple command-line tool to quickly set up and manage CMake projects.

## Installation

```
go build -o qs
```

Move the compiled binary to a location in your PATH.

## Usage

### Initialize a CMake project

```
qs init
```

This creates a basic CMakeLists.txt file in the current directory.

### Add an executable target

```
qs add <target_name> [files...]
```

Examples:
- `qs add hello` - Adds target "hello" with hello.cpp or hello.c (if found)
- `qs add hello hello.cpp utils.cpp` - Adds target with explicit source files
- `qs add hello src` - Adds all source files from the "src" directory
- `qs add hello *.cpp` - Adds all .cpp files in the current directory
- `qs add hello src/*.cpp include/*.hpp` - Adds files matching multiple glob patterns

Glob patterns support:
- `*` - Matches any sequence of characters
- `?` - Matches any single character
- `[abc]` - Matches any character in the brackets

### Add standard CMake configuration

```
qs std [cxx_std]
```

This adds common CMake configuration to your project:
- Compiler flags (-Wall -Wextra)
- Organized output directories (bin, lib)
- Standard include directories
- Testing support

You can optionally specify the C++ standard:
- `qs std 11` - Sets C++11 standard
- `qs std 14` - Sets C++14 standard (default)
- `qs std 17` - Sets C++17 standard
- `qs std 20` - Sets C++20 standard

## Examples

Create a basic hello world project:

```
mkdir hello-project
cd hello-project
qs init
echo '#include <iostream>
int main() {
    std::cout << "Hello, world!" << std::endl;
    return 0;
}' > hello.cpp
qs add hello
qs std 17  # Use C++17 standard
mkdir build && cd build
cmake ..
make
```

The executable will be in the `build/bin` directory.

Create a project with multiple source files:

```
mkdir multi-src
cd multi-src
qs init
# Create several source files...
qs add myapp *.cpp *.c
qs std
``` 
