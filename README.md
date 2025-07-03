# konfig

Welcome to konfig! This library provides a simple and powerful way to manage configurations in your Go applications. If you're familiar with the Spring Framework, you'll feel right at home with konfig's approach to configuration management.

konfig aims to simplify your development workflow by offering:

* **Spring-like Configuration:** Enjoy a familiar configuration style, making it easy to organize and access your application settings.
* **Profile Support:** Effortlessly manage different configurations for various environments (e.g., development, production, testing) using profiles.
* **Environment Variable Overrides:** Seamlessly override configuration settings using environment variables, providing flexibility and control.
* **YAML Configuration:** Load your configurations from clearly structured and human-readable YAML files.

![img.png](docs/schema.png)

## Features

*   **YAML-based Configuration:** Load your application configurations from clearly structured and human-readable YAML files.
*   **Profile-specific Configurations:** Define default configurations in `application.yaml` and override or extend them with profile-specific files like `application-dev.yaml` or `application-prod.yaml`.
*   **Environment Variable Substitution:** Use environment variables directly within your YAML files for dynamic configuration values (e.g., `db_host: ${DB_HOST:localhost}`). If the environment variable `DB_HOST` is set, its value will be used; otherwise, `localhost` will be the default.
*   **Environment Variable Overrides:** Configuration values loaded from files are automatically exposed as environment variables. This allows you to override any configuration setting by simply setting an environment variable with the corresponding key (e.g., setting `SERVER_PORT=8081` in your environment will override the `server.port` value from your YAML files).
*   **Profile Selection via Command-Line:** Easily switch between configuration profiles using command-line flags (`-p` or `--profile`). For example, running your application with `-p dev` will load the `application-dev.yaml` profile.
*   **Automatic Key Flattening:** Nested YAML keys are automatically flattened into a dot-separated format. For instance, a YAML structure like `server: {port: 8080}` will be accessible as `server.port`.

## Installation

To install the `konfig` library, use the following `go get` command:

```bash
go get github.com/mfenderov/konfig
```

## Usage

This section explains how to integrate and use the `konfig` library in your Go applications.

### Basic Initialization

To load the default configuration, which typically includes `application.yaml` and any active profile-specific YAML file (e.g., `application-dev.yaml`), simply call:

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

By default, `konfig` looks for configuration files in a `resources/` directory relative to the application's working directory.

### Loading Specific Files

You can also load a configuration from a specific YAML file:

```go
err := konfig.LoadFrom("path/to/your/custom-config.yaml")
if err != nil {
    // Handle error
}
```
This will load only the specified file. If you want to include default application and profile configurations as well, load them first or separately.

### Profile Usage

`konfig` supports profile-specific configurations. Profiles allow you to tailor your application settings for different environments (e.g., development, testing, production).

**Activating Profiles:**

Profiles are typically activated using command-line flags when running your application:

```bash
go run main.go -p dev  # Activates the 'dev' profile
# or
go run main.go --profile prod # Activates the 'prod' profile
```

This will load `application-<profile>.yaml` (e.g., `application-dev.yaml`) after `application.yaml`, with profile values overriding the defaults.

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

`konfig` loads all configuration key-value pairs into the application's environment. This means you can access any configuration value using the standard `os.Getenv()` function. Keys are flattened, so a nested key like `server: { port: 8080 }` in YAML becomes `server.port`.

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
When the `dev` profile is active, `server.port` will be `3000` and `database.url` will be the one specified in `application-dev.yaml`. `application.name` would still be `MyApp`.

## Configuration File Structure

`konfig` expects configuration files to be located in a specific directory and follow a naming convention to enable automatic loading and profile management.

### Default Directory

*   By default, `konfig` looks for configuration files in a `resources/` directory relative to the application's working directory.

### File Naming Conventions

*   **Default Configuration File:** The primary configuration file should be named `application.yaml` or `application.yml`. This file contains the base configuration for your application.
*   **Profile-Specific Configuration Files:** To define configurations for specific profiles (e.g., development, production, testing), create files named `application-<profile>.yaml` or `application-<profile>.yml`.
    *   Example for a `dev` profile: `resources/application-dev.yaml`
    *   Example for a `prod` profile: `resources/application-prod.yml`

### Loading Order and Overrides

When `konfig.Load()` is called (or when a profile is activated):
1.  `konfig` first attempts to load the default configuration file (`resources/application.yaml` or `resources/application.yml`).
2.  If a profile is active (e.g., `dev`), `konfig` then attempts to load the corresponding profile-specific file (e.g., `resources/application-dev.yaml`).
3.  Values from the profile-specific configuration file will override any values with the same keys that were loaded from the default configuration file. This allows you to define common settings in `application.yaml` and tailor specific settings in profile files.

If you use `konfig.LoadFrom("custom/path/config.yaml")`, these default loading mechanisms and directory structures are bypassed, and only the specified file is loaded.

## Contributing

We welcome and encourage contributions to `konfig`! If you'd like to help improve the library, please follow these steps:

1.  **Fork the repository** on GitHub.
2.  **Create a new branch** for your feature or bug fix.
    *   For features: `git checkout -b feature/your-feature-name`
    *   For bug fixes: `git checkout -b fix/issue-number-or-description`
3.  **Make your changes** and commit them with clear, descriptive messages.
4.  **Push your changes** to your forked repository.
5.  **Create a Pull Request** from your branch to the `main` branch of the `mfenderov/konfig` repository.

Please ensure your code adheres to the project's standards. For more detailed information on the project structure, coding standards, testing procedures, and other development guidelines, please refer to the [Developer Guidelines](.junie/guidelines.md).

## License

This project is licensed under the MIT License.