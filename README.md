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

Note: For safety reasons, this command cannot be run in your home directory. Create a specific directory for your project first.

### Add an executable target

```
qs add <target_name> [files...]
```

Examples:
- `qs add hello` - Adds target "hello" with hello.cpp, hello.cc, hello.c, or hello.cxx (if found)
- `qs add hello hello.cpp utils.cpp` - Adds target with explicit source files
- `qs add hello src` - Adds all source files from the "src" directory
- `qs add hello *.cpp` - Adds all .cpp files in the current directory
- `qs add hello src/*.cpp include/*.hpp` - Adds files matching multiple glob patterns

Glob patterns support:
- `*` - Matches any sequence of characters
- `?` - Matches any single character
- `[abc]` - Matches any character in the brackets

If a target with the same name already exists, the new source files will be appended to that target. This allows you to add more source files to an existing target without manually editing the CMakeLists.txt file.

### Build project

```
qs build
```

Creates a build directory, runs cmake and make to build your project.

### Run project

```
qs run [target]
```

Runs the specified executable target (or the default target if not specified).

### List targets

```
qs list
```

Lists all available targets in the project, including:
- Executable targets defined in CMakeLists.txt
- Library targets defined in CMakeLists.txt
- Actual built executables in the build directory

This command helps you see what targets are available for building and running.

### Add standard CMake configuration

```
qs std [cxx_std]
```

This adds common CMake configuration to your project:
- Compiler flags (-Wall -Wextra)
- Organized output directories (bin, lib)
- Standard include directories
- Testing support
- Install targets

You can optionally specify the C++ standard:
- `qs std 11` - Sets C++11 standard
- `qs std 14` - Sets C++14 standard (default)
- `qs std 17` - Sets C++17 standard
- `qs std 20` - Sets C++20 standard

### Other commands

```
qs doc          - Open CMake documentation in the default browser
qs version      - Display version information
qs help         - Show help message
```

## Examples

Create a basic hello world project:

```
mkdir hello-project
cd hello-project
qs init
qs build
qs run
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