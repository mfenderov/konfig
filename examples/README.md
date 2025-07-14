# Konfig Examples

Working examples demonstrating konfig library usage.

## Quick Start

```bash
# Basic example
go run simple_example.go

# Basic example with profiles
go run simple_example.go -p dev
go run simple_example.go -p prod

# Comprehensive example with advanced features
go run simple_example.go struct_config_example.go

# With environment variables
SERVER_PORT=9090 go run simple_example.go
```

## Files

- `simple_example.go` - Runnable basic struct-based configuration example
- `struct_config_example.go` - Advanced nested configuration library (call `RunStructExample()`)
- `resources/` - Example YAML configuration files

## Common Issues

**"must be a pointer to struct"** - Use `&config`, not `config`
**Fields not populating** - Add `konfig:"key.path"` tags to struct fields
**Nested structs not working** - Add `konfig:"prefix"` tag to nested struct field