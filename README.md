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

### Create a sub-project

```
qs init sub <name>
```

Creates a subdirectory with its own CMakeLists.txt file configured as a library. 
This command:
- Creates the specified subdirectory with src/ and include/ folders
- Adds a CMakeLists.txt in the subdirectory configured to build a library
- Updates the parent CMakeLists.txt to include the subdirectory
- Creates sample header and source files

This is useful for organizing larger projects with multiple components.

#### Using a sub-project's include files

After creating a sub-project, you can use its header files in other targets:

1. **Include the header file in your source**:
   ```cpp
   #include "<subproject_name>/<subproject_name>.h"
   // Or for the sample header directly:
   #include "<subproject_name>.h"
   ```

2. **Link the library to your executable** in your main CMakeLists.txt:
   ```cmake
   # After your add_executable command
   target_link_libraries(your_executable PRIVATE <subproject_name>)
   ```

The sub-project is already set up with proper include paths, so once you link against it, your executable will have access to all of its public headers.

Example usage:
```cpp
// In main.cpp
#include "<subproject_name>.h"

int main() {
    <subproject_name>::Example example;
    example.doSomething();
    return 0;
}
```

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

Create a project with a sub-project:

```
# Create main project
mkdir my-project
cd my-project
qs init

# Create a sub-project library
qs init sub utils

# Create a main executable that uses the utils library
echo '#include <iostream>
#include "utils.h"

int main() {
    utils::Example utils;
    utils.doSomething();
    std::cout << "Main program says hello!" << std::endl;
    return 0;
}' > main.cpp

# Add the main executable
qs add main main.cpp

# Build and run
qs build
qs run
```

This example creates a project structure like:
```
my-project/
├── CMakeLists.txt
├── main.cpp
├── utils/
│   ├── CMakeLists.txt
│   ├── include/
│   │   └── utils.h
│   └── src/
│       └── utils.cpp
└── build/
```