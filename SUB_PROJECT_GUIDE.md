# Working with Sub-Projects in QS

This guide explains how to create and use sub-projects in your CMake projects using the QS tool.

## What is a Sub-Project?

A sub-project is a directory within your main project that contains its own CMakeLists.txt file. Sub-projects are typically used to:

- Create reusable libraries
- Organize code into logical components
- Manage dependencies between components
- Improve build times by isolating changes

## Creating a Sub-Project

To create a sub-project, use the `qs init sub` command:

```
qs init sub <name>
```

For example:

```
qs init sub utils
```

This command:

1. Creates a directory named `utils`
2. Creates the following structure:
   ```
   utils/
   ├── CMakeLists.txt
   ├── include/
   │   ├── utils.h
   │   └── utils/
   │       └── example.h
   └── src/
       └── utils.cpp
   ```
3. Updates the parent CMakeLists.txt to include the subdirectory

## Sub-Project Structure

Each sub-project created by `qs init sub` follows a standard layout:

- **include/**: Contains header files (public API)
  - **include/<subproject_name>/**: Contains organized headers for proper inclusion
- **src/**: Contains implementation files
- **CMakeLists.txt**: Configuration for building the sub-project as a library

The generated CMakeLists.txt is configured to:

- Build a library with the same name as the sub-project
- Set up include directories properly
- Define installation rules
- Configure proper interface include paths

## Using a Sub-Project's Library

### In the Same Project

Once you've created a sub-project, you can use it in your main project by:

1. **Including the headers** in your source code:

   ```cpp
   #include "<subproject_name>/example.h"
   ```

   This is the recommended way to include headers from sub-projects.

   For backward compatibility, you can also use the direct header at the root:
   
   ```cpp
   #include "<subproject_name>.h"
   ```

2. **Linking with the library** in your CMakeLists.txt:

   ```cmake
   add_executable(my_app main.cpp)
   target_link_libraries(my_app PRIVATE <subproject_name>)
   ```

The include paths are automatically set up because of the `target_include_directories` command in the sub-project's CMakeLists.txt.

### IMPORTANT: Including Headers Correctly

When working with a sub-project, especially one named with common terms like "hello", "math", etc., always use the namespaced include format to avoid conflicts:

```cpp
#include "hello/example.h"  // Good - uses the namespaced format
```

Don't rely on the direct include, as it might not work in all cases:

```cpp
#include "hello.h"  // Not recommended - may cause conflicts
```

#### Example with a sub-project named "hello"

```cpp
#include "hello/example.h"

int main() {
    hello::Example example;
    example.doSomething();
    return 0;
}
```

### In External Projects

To use a built and installed library from a sub-project in another project:

1. Install your library:
   ```
   cmake --build . --target install
   ```

2. In the other project, find and link to the package:
   ```cmake
   find_package(<subproject_name> REQUIRED)
   target_link_libraries(your_target PRIVATE <subproject_name>)
   ```

## Example: Complete Project with Sub-Project

Here's a complete example of creating a project with a sub-project and using it:

```bash
# Create the main project
mkdir my-project
cd my-project
qs init

# Create a utilities sub-project
qs init sub utils

# Create a main program that uses the utils library
cat > main.cpp << EOF
#include <iostream>
#include "utils.h"

int main() {
    std::cout << "Main program starting..." << std::endl;
    
    utils::Example utils;
    utils.doSomething();
    
    std::cout << "Main program completed." << std::endl;
    return 0;
}
EOF

# Add the main executable target
qs add main main.cpp

# Build the project
qs build

# List all targets
qs list

# Run the program
qs run main
```

## Adding More Headers to a Sub-Project

To add more headers to your sub-project:

1. Create your header files in the `<subproject_name>/include/` directory
2. Use them in your source files in the `<subproject_name>/src/` directory
3. Include them in your main project:
   ```cpp
   #include "<subproject_name>/<header_name>.h"
   ```

For example:

```bash
# Create a new header in the utils sub-project
cat > utils/include/math_utils.h << EOF
#pragma once

namespace utils {

class MathUtils {
public:
    static int add(int a, int b);
    static int multiply(int a, int b);
};

} // namespace utils
EOF

# Create the implementation
cat > utils/src/math_utils.cpp << EOF
#include "math_utils.h"

namespace utils {

int MathUtils::add(int a, int b) {
    return a + b;
}

int MathUtils::multiply(int a, int b) {
    return a * b;
}

} // namespace utils
EOF

# Update the main program to use it
cat > main.cpp << EOF
#include <iostream>
#include "utils.h"
#include "math_utils.h"

int main() {
    utils::Example utils;
    utils.doSomething();
    
    int sum = utils::MathUtils::add(5, 7);
    int product = utils::MathUtils::multiply(5, 7);
    
    std::cout << "5 + 7 = " << sum << std::endl;
    std::cout << "5 * 7 = " << product << std::endl;
    
    return 0;
}
EOF

# Rebuild and run
qs build
qs run main
```

## Dependencies Between Sub-Projects

You can create dependencies between sub-projects by linking them:

```bash
# Create another sub-project
qs init sub network

# Update its CMakeLists.txt to depend on utils
sed -i '' 's/# target_link_libraries(network PRIVATE dependency1 dependency2)/target_link_libraries(network PRIVATE utils)/' network/CMakeLists.txt

# Create a header that uses utils
cat > network/include/client.h << EOF
#pragma once
#include "utils.h"  // Including from another sub-project

namespace network {

class Client {
public:
    Client();
    void connect();
private:
    utils::Example m_utils;
};

} // namespace network
EOF

# Create the implementation
cat > network/src/client.cpp << EOF
#include "client.h"
#include <iostream>

namespace network {

Client::Client() {}

void Client::connect() {
    std::cout << "Client connecting..." << std::endl;
    m_utils.doSomething();
    std::cout << "Client connected." << std::endl;
}

} // namespace network
EOF

# Update main to use both sub-projects
cat > main.cpp << EOF
#include <iostream>
#include "utils.h"
#include "math_utils.h"
#include "client.h"

int main() {
    utils::Example utils;
    utils.doSomething();
    
    network::Client client;
    client.connect();
    
    return 0;
}
EOF

# Update main's CMakeLists.txt to link to both libraries
# (You need to edit the target_link_libraries line for the main target)
```

## Best Practices

1. **Namespace your code**: Use namespaces to prevent conflicts between sub-projects
2. **Minimize dependencies**: Keep dependencies between sub-projects as minimal as possible
3. **Public vs. Private headers**: Only expose what's necessary in your public API
4. **Use consistent naming**: Keep your naming consistent across sub-projects
5. **Document interfaces**: Document the public interfaces of your libraries 