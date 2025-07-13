# Konfig Examples

This directory contains comprehensive examples demonstrating how to use the konfig library, especially the new struct-based configuration feature.

## Examples Included

### 1. Simple Example (`simple_example.go`)

A basic example showing how to use `LoadInto()` with a simple configuration structure.

**Features demonstrated:**
- Basic struct definition with konfig tags
- Default values
- Environment variable overrides
- Profile detection

**Run it:**
```bash
go run simple_example.go
go run simple_example.go -p dev
go run simple_example.go -p prod
```

### 2. Comprehensive Example (`struct_config_example.go`)

A complete real-world example showing advanced configuration patterns.

**Features demonstrated:**
- Complex nested structure (4+ levels deep)
- Multiple configuration domains (server, database, security, etc.)
- Profile-aware configuration logic
- Configuration validation
- Environment variable substitution
- Mixed defaults and overrides

**Run it:**
```bash
go run struct_config_example.go
go run struct_config_example.go -p dev
go run struct_config_example.go -p prod
```

## Configuration Files

The `resources/` directory contains example YAML files:

- `application.yaml` - Base configuration with defaults
- `application-dev.yaml` - Development profile overrides
- `application-prod.yaml` - Production profile overrides

These demonstrate:
- Environment variable substitution (`${VAR:default}`)
- Profile-specific overrides
- Nested YAML structure mapping to Go structs
- Security-conscious production settings

## Key Patterns Demonstrated

### 1. Struct Tag Usage

```go
type Config struct {
    // Simple field mapping
    Host string `konfig:"server.host" default:"localhost"`
    
    // Nested structure mapping
    Database DatabaseConfig `konfig:"database"`
    
    // Environment variable with default
    Port string `konfig:"server.port" default:"8080"`
    
    // Required field (no default)
    APIKey string `konfig:"security.api_key"`
}
```

### 2. Environment Variable Substitution in YAML

```yaml
server:
  port: ${SERVER_PORT:8080}  # Uses SERVER_PORT env var, defaults to 8080
  host: ${SERVER_HOST}       # Must be provided via environment

database:
  url: ${DB_URL:postgres://localhost/myapp}
```

### 3. Profile-Specific Configuration

```go
func main() {
    var config AppConfig
    konfig.LoadInto(&config)
    
    if konfig.IsDevProfile() {
        // Development-specific logic
        config.Debug = true
    } else if konfig.IsProdProfile() {
        // Production validation
        validateProductionConfig(&config)
    }
}
```

### 4. Configuration Validation

```go
func validateConfig(config *AppConfig) {
    if konfig.IsProdProfile() && config.JWT.Secret == "default-secret" {
        log.Fatal("Production requires secure JWT secret")
    }
    
    if config.Database.Host == "" {
        log.Fatal("Database host is required")
    }
}
```

## Running the Examples

1. **Prerequisites:**
   ```bash
   go mod init example
   go get github.com/mfenderov/konfig
   ```

2. **Basic usage:**
   ```bash
   go run simple_example.go
   ```

3. **With profiles:**
   ```bash
   go run struct_config_example.go -p dev
   go run struct_config_example.go -p prod
   ```

4. **With environment variables:**
   ```bash
   SERVER_PORT=9090 DB_HOST=mydb.com go run simple_example.go
   ```

## Expected Output

The examples will show:
- ✅ Successful configuration loading
- 📱 Application metadata
- 🌐 Server configuration
- 🗄️ Database settings
- 🔐 Security configuration
- 🏷️ Active profile information
- ✅ Validation results

## Best Practices Demonstrated

1. **Struct Organization:** Group related configuration into nested structs
2. **Default Values:** Provide sensible defaults for all optional fields
3. **Environment Integration:** Use environment variables for deployment-specific values
4. **Profile Usage:** Leverage profiles for environment-specific configurations
5. **Validation:** Add validation logic for critical configuration values
6. **Security:** Mask sensitive values in logs and outputs

## Troubleshooting

### Common Issues

1. **"must be a pointer to struct"**
   ```go
   // Wrong
   err := konfig.LoadInto(config)
   
   // Correct
   err := konfig.LoadInto(&config)
   ```

2. **Fields not populating**
   ```go
   // Missing konfig tag
   Host string `default:"localhost"`  // Won't work
   
   // Correct
   Host string `konfig:"server.host" default:"localhost"`
   ```

3. **Nested structures not working**
   ```go
   // Missing konfig tag on nested struct
   Server ServerConfig  // Won't populate nested fields
   
   // Correct
   Server ServerConfig `konfig:"server"`
   ```

For more information, see the main [konfig README](../README.md).