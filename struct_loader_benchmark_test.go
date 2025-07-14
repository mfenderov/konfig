package konfig

import (
	"os"
	"testing"
)

// Essential benchmark tests for LoadInto functionality

func BenchmarkLoadInto_SimpleStruct(b *testing.B) {
	type Config struct {
		Host string `konfig:"bench.host" default:"localhost"`
		Port string `konfig:"bench.port" default:"8080"`
	}

	os.Setenv("bench.host", "example.com")
	os.Setenv("bench.port", "9090")
	defer func() {
		os.Unsetenv("bench.host")
		os.Unsetenv("bench.port")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(&cfg)
		if err != nil {
			b.Fatalf("LoadInto failed: %v", err)
		}
	}
}

func BenchmarkLoadInto_NestedStruct(b *testing.B) {
	type DatabaseConfig struct {
		Host string `konfig:"host" default:"localhost"`
		Port string `konfig:"port" default:"5432"`
		Name string `konfig:"name" default:"myapp"`
	}

	type ServerConfig struct {
		Host string `konfig:"host" default:"0.0.0.0"`
		Port string `konfig:"port" default:"8080"`
	}

	type Config struct {
		Database DatabaseConfig `konfig:"benchnested.database"`
		Server   ServerConfig   `konfig:"benchnested.server"`
	}

	os.Setenv("benchnested.database.host", "db.example.com")
	os.Setenv("benchnested.server.port", "9090")
	defer func() {
		os.Unsetenv("benchnested.database.host")
		os.Unsetenv("benchnested.server.port")
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(&cfg)
		if err != nil {
			b.Fatalf("LoadInto failed: %v", err)
		}
	}
}

// Memory allocation benchmark
func BenchmarkLoadInto_MemoryAllocation(b *testing.B) {
	type Config struct {
		Host string `konfig:"benchmem.host" default:"localhost"`
		Port string `konfig:"benchmem.port" default:"8080"`
		Name string `konfig:"benchmem.name" default:"myapp"`
	}

	os.Setenv("benchmem.host", "example.com")
	defer os.Unsetenv("benchmem.host")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(&cfg)
		if err != nil {
			b.Fatalf("LoadInto failed: %v", err)
		}
	}
}
