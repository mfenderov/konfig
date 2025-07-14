# konfig

Welcome to konfig! This library provides a simple and powerful way to manage configurations in your Go applications. If
you're familiar with the Spring Framework, you'll feel right at home with konfig's approach to configuration management.

konfig is a **simplified, high-performance** Go configuration library offering:

* **Spring-like Configuration:** Familiar configuration style for easy organization and access
* **Struct-based Loading:** Type-safe configuration with automatic population using struct tags
* **Profile Support:** Environment-specific configurations (dev, prod, test) with command-line switching
* **Environment Variable Integration:** Seamless overrides and substitution support
* **YAML Configuration:** Clean, readable configuration files with nested structure support

![img.png](docs/schema.png)

## Features

* **YAML-based Configuration:** Load your application configurations from clearly structured and human-readable YAML
  files.
* **Struct-based Configuration (NEW):** Automatically populate Go structs with configuration values using struct tags,
  providing type safety and better IDE support.
* **Profile-specific Configurations:** Define default configurations in `application.yaml` and override or extend them
  with profile-specific files like `application-dev.yaml` or `application-prod.yaml`.
* **Environment Variable Substitution:** Use environment variables directly within your YAML files for dynamic
  configuration values (e.g., `db_host: ${DB_HOST:localhost}`). If the environment variable `DB_HOST` is set, its value
  will be used; otherwise, `localhost` will be the default.
* **Environment Variable Overrides:** Configuration values loaded from files are automatically exposed as environment
  variables. This allows you to override any configuration setting by simply setting an environment variable with the
  corresponding key (e.g., setting `SERVER_PORT=8081` in your environment will override the `server.port` value from
  your YAML files).
* **Profile Selection via Command-Line:** Easily switch between configuration profiles using command-line flags (`-p` or
  `--profile`). For example, running your application with `-p dev` will load the `application-dev.yaml` profile.
* **Automatic Key Flattening:** Nested YAML keys are automatically flattened into a dot-separated format. For instance,
  a YAML structure like `server: {port: 8080}` will be accessible as `server.port`.
* **Nested Struct Support:** Deep nesting of configuration structures with automatic prefix handling.
* **Default Value Support:** Built-in default values using struct tags when configuration keys are not provided.

## Installation

To install the `konfig` library, use the following `go get` command:

```bash
go get github.com/mfenderov/konfig
```

## Usage

This section explains how to integrate and use the `konfig` library in your Go applications.

### Basic Initialization

To load the default configuration, which typically includes `application.yaml` and any active profile-specific YAML
file (e.g., `application-dev.yaml`), simply call:

```go
import "github.com/mfenderov/konfig"

func main() {
err := konfig.Load()
if err != nil {
// Handle error
log.Fatalf("Failed to load configuration: %v", err)
}
// Your application logic here
}
```

By default, `konfig` looks for configuration files in a `resources/` directory relative to the application's working
directory.

### Loading Specific Files

You can also load a configuration from a specific YAML file:

```go
err := konfig.LoadFrom("path/to/your/custom-config.yaml")
if err != nil {
// Handle error
}
```

This will load only the specified file. If you want to include default application and profile configurations as well,
load them first or separately.

### Profile Usage

`konfig` supports profile-specific configurations. Profiles allow you to tailor your application settings for different
environments (e.g., development, testing, production).

**Activating Profiles:**

Profiles are typically activated using command-line flags when running your application:

```bash
go run main.go -p dev  # Activates the 'dev' profile
# or
go run main.go --profile prod # Activates the 'prod' profile
```

This will load `application-<profile>.yaml` (e.g., `application-dev.yaml`) after `application.yaml`, with profile values
overriding the defaults.

**Checking Active Profiles:**

You can programmatically check the active profile:

```go
if konfig.IsDevProfile() {
fmt.Println("Running with development profile")
}

if konfig.IsProdProfile() {
fmt.Println("Running with production profile")
}

if konfig.IsProfile("staging") {
fmt.Println("Running with staging profile")
}

activeProfile := konfig.GetProfile() // Returns the current profile name (e.g., "dev", "prod", or "" if none)
fmt.Printf("Current active profile: %s\n", activeProfile)
```

### Accessing Configuration Values

There are two main ways to access configuration values with konfig:

#### 1. Environment Variable Access (Traditional)

`konfig` loads all configuration key-value pairs into the application's environment. This means you can access any
configuration value using the standard `os.Getenv()` function. Keys are flattened, so a nested key like
`server: { port: 8080 }` in YAML becomes `server.port`.

```go
import (
"fmt"
"os"
"github.com/mfenderov/konfig"
)

func main() {
err := konfig.Load() // Load default configs
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}

// Example: Accessing server.port
serverPort := os.Getenv("server.port")
fmt.Printf("Server port: %s\n", serverPort)

// Example: Accessing database.url
dbURL := os.Getenv("database.url")
fmt.Printf("Database URL: %s\n", dbURL)

// Example: Accessing a value that might use environment variable substitution
apiKey := os.Getenv("security.api_key")
fmt.Printf("API Key: %s\n", apiKey)
}
```

#### 2. Struct-Based Configuration (Recommended)

**New in v0.17.0:** konfig now supports automatic population of Go structs using struct tags. This approach provides
type safety, better IDE support, and cleaner code organization.

```go
import (
"fmt"
"log"
"github.com/mfenderov/konfig"
)

// Define your configuration structure with konfig tags
type DatabaseConfig struct {
Host     string `konfig:"host" default:"localhost"`
Port     string `konfig:"port" default:"5432"`
Name     string `konfig:"name" default:"myapp"`
User     string `konfig:"user" default:"postgres"`
Password string `konfig:"password" default:"secret"`
}

type ServerConfig struct {
Host string `konfig:"host" default:"0.0.0.0"`
Port string `konfig:"port" default:"8080"`
}

type AppConfig struct {
Name     string         `konfig:"application.name" default:"MyApp"`
Version  string         `konfig:"application.version" default:"1.0.0"`
Database DatabaseConfig `konfig:"database"`
Server   ServerConfig   `konfig:"server"`
}

func main() {
// Load configuration into struct
var config AppConfig
err := konfig.LoadInto(&config)
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}

// Access configuration with type safety
fmt.Printf("App: %s v%s\n", config.Name, config.Version)
fmt.Printf("Server: %s:%s\n", config.Server.Host, config.Server.Port)
fmt.Printf("Database: %s@%s:%s/%s\n",
config.Database.User, config.Database.Host,
config.Database.Port, config.Database.Name)
}
```

**Struct Tag Reference:**

- `konfig:"key.path"` - Maps the field to a configuration key (supports dot notation for nested keys)
- `default:"value"` - Provides a default value if the configuration key is not set
- Fields without `konfig` tags are ignored during population

**Benefits of Struct-Based Configuration:**

- **Type Safety:** Compile-time checking prevents runtime errors
- **IDE Support:** Auto-completion and refactoring support
- **Documentation:** Struct fields serve as self-documenting configuration
- **Validation:** Easy to add validation logic to your configuration structs
- **Testing:** Simple to create test configurations by initializing structs

#### Advanced Usage

**Nested structures and profile-aware configuration:**

```go
type Config struct {
App struct {
Name    string `konfig:"name" default:"MyApp"`
Version string `konfig:"version" default:"1.0.0"`
} `konfig:"application"`

Database struct {
Host string `konfig:"host" default:"localhost"`
Port string `konfig:"port" default:"5432"`
} `konfig:"database"`
}

func main() {
var config Config
if err := konfig.LoadInto(&config); err != nil {
log.Fatal(err)
}

// Profile-specific logic
if konfig.IsDevProfile() {
config.Database.Host = "localhost"
}

fmt.Printf("%s v%s connecting to %s:%s\n",
config.App.Name, config.App.Version,
config.Database.Host, config.Database.Port)
}
```

### Example YAML Configuration

Here's an example of an `application.yaml` file:

```yaml
# application.yaml

application:
  name: MyApp
  version: "1.0.0"

server:
  host: localhost
  port: ${APP_PORT:8080} # Uses APP_PORT env var, defaults to 8080 if not set

database:
  url: ${DB_URL} # Must be provided via environment variable
  username: user
  password: ${DB_PASSWORD:secret}

security:
  api_key: ${API_KEY} # Must be provided via environment variable
  jwt:
    secret: ${JWT_SECRET:verysecretkey}
    expiration_hours: 24

logging:
  level: info # Default logging level
```

If a profile, for example `application-dev.yaml`, is active, it can override these values:

```yaml
# application-dev.yaml

server:
  port: 3000 # Overrides default port for dev environment

database:
  url: "postgres://devuser:devpass@localhost:5432/devdb?sslmode=disable"

logging:
  level: debug # More verbose logging for dev
```

When the `dev` profile is active, `server.port` will be `3000` and `database.url` will be the one specified in
`application-dev.yaml`. `application.name` would still be `MyApp`.

## API Reference

### Core Functions

#### `konfig.Load() error`

Loads the default configuration files (`application.yaml` and active profile files) and populates environment variables.

```go
err := konfig.Load()
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}
```

#### `konfig.LoadFrom(filepath string) error`

Loads configuration from a specific YAML file.

```go
err := konfig.LoadFrom("config/custom.yaml")
if err != nil {
log.Fatalf("Failed to load custom config: %v", err)
}
```

#### `konfig.LoadInto(config interface{}) error` (NEW)

Loads configuration into a Go struct using struct tags.

```go
type Config struct {
Port string `konfig:"server.port" default:"8080"`
}

var cfg Config
err := konfig.LoadInto(&cfg)
if err != nil {
log.Fatalf("Failed to load into struct: %v", err)
}
```

**Requirements:**

- `config` must be a pointer to a struct
- Struct fields must have `konfig` tags to be populated
- Supports nested structs with recursive population

### Profile Functions

#### `konfig.GetProfile() string`

Returns the currently active profile name.

```go
profile := konfig.GetProfile() // Returns "dev", "prod", or ""
```

#### `konfig.IsDevProfile() bool`

Returns true if the current profile is "dev".

```go
if konfig.IsDevProfile() {
// Development-specific logic
}
```

#### `konfig.IsProdProfile() bool`

Returns true if the current profile is "prod".

```go
if konfig.IsProdProfile() {
// Production-specific logic
}
```

#### `konfig.IsProfile(name string) bool`

Returns true if the current profile matches the given name.

```go
if konfig.IsProfile("staging") {
// Staging-specific logic
}
```

### Struct Tag Reference

When using `LoadInto()`, struct fields can use the following tags:

- **`konfig:"key.path"`** - Maps the field to a configuration key (required for field to be populated)
- **`default:"value"`** - Provides a default value if the configuration key is not found

**Examples:**

```go
type Config struct {
// Simple key mapping with default
Host string `konfig:"server.host" default:"localhost"`

// Nested key mapping
DBPassword string `konfig:"database.auth.password" default:"secret"`

// No default - will be empty string if not found
APIKey string `konfig:"security.api_key"`

// Ignored field (no konfig tag)
RuntimeValue string
}
```

### Error Handling

```go
if err := konfig.LoadInto(&config); err != nil {
log.Fatalf("Configuration error: %v", err)
}
```

## Configuration File Structure

`konfig` expects configuration files to be located in a specific directory and follow a naming convention to enable
automatic loading and profile management.

### Default Directory

* By default, `konfig` looks for configuration files in a `resources/` directory relative to the application's working
  directory.

### File Naming Conventions

* **Default Configuration File:** The primary configuration file should be named `application.yaml` or
  `application.yml`. This file contains the base configuration for your application.
* **Profile-Specific Configuration Files:** To define configurations for specific profiles (e.g., development,
  production, testing), create files named `application-<profile>.yaml` or `application-<profile>.yml`.
    * Example for a `dev` profile: `resources/application-dev.yaml`
    * Example for a `prod` profile: `resources/application-prod.yml`

### Loading Order and Overrides

When `konfig.Load()` is called (or when a profile is activated):

1. `konfig` first attempts to load the default configuration file (`resources/application.yaml` or
   `resources/application.yml`).
2. If a profile is active (e.g., `dev`), `konfig` then attempts to load the corresponding profile-specific file (e.g.,
   `resources/application-dev.yaml`).
3. Values from the profile-specific configuration file will override any values with the same keys that were loaded from
   the default configuration file. This allows you to define common settings in `application.yaml` and tailor specific
   settings in profile files.

If you use `konfig.LoadFrom("custom/path/config.yaml")`, these default loading mechanisms and directory structures are
bypassed, and only the specified file is loaded.

## Examples

See the `examples/` directory for working code samples:

```bash
cd examples
go run simple_example.go
go run struct_config_example.go -p dev
```

## Project Quality

konfig maintains high standards through systematic simplification:

- **Focused test suite**: 6 essential test files covering all functionality
- **Performance benchmarks**: 3 core benchmarks measuring real-world scenarios
- **Clean documentation**: Concise, example-driven documentation
- **100% functionality preservation**: All features maintained through simplification process

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Submit a Pull Request

## License

This project is licensed under the MIT License.