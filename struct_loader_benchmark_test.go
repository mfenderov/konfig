package konfig

import (
	"fmt"
	"os"
	"testing"
)

// Benchmark tests for LoadInto functionality

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

func BenchmarkLoadInto_LargeStruct(b *testing.B) {
	type DatabaseConfig struct {
		Host            string `konfig:"host" default:"localhost"`
		Port            string `konfig:"port" default:"5432"`
		Name            string `konfig:"name" default:"myapp"`
		User            string `konfig:"user" default:"postgres"`
		Password        string `konfig:"password" default:"secret"`
		SSLMode         string `konfig:"ssl_mode" default:"disable"`
		MaxConnections  string `konfig:"max_connections" default:"100"`
		ConnMaxLifetime string `konfig:"conn_max_lifetime" default:"1h"`
	}

	type ServerConfig struct {
		Host           string `konfig:"host" default:"0.0.0.0"`
		Port           string `konfig:"port" default:"8080"`
		ReadTimeout    string `konfig:"read_timeout" default:"30s"`
		WriteTimeout   string `konfig:"write_timeout" default:"30s"`
		MaxHeaderBytes string `konfig:"max_header_bytes" default:"1048576"`
		TLSCertFile    string `konfig:"tls_cert_file" default:""`
		TLSKeyFile     string `konfig:"tls_key_file" default:""`
	}

	type LoggingConfig struct {
		Level      string `konfig:"level" default:"info"`
		Format     string `konfig:"format" default:"json"`
		Output     string `konfig:"output" default:"stdout"`
		MaxSize    string `konfig:"max_size" default:"100"`
		MaxBackups string `konfig:"max_backups" default:"3"`
		MaxAge     string `konfig:"max_age" default:"28"`
	}

	type RedisConfig struct {
		Host     string `konfig:"host" default:"localhost"`
		Port     string `konfig:"port" default:"6379"`
		Password string `konfig:"password" default:""`
		DB       string `konfig:"db" default:"0"`
	}

	type Config struct {
		App      string         `konfig:"benchlarge.app.name" default:"myapp"`
		Version  string         `konfig:"benchlarge.app.version" default:"1.0.0"`
		Debug    string         `konfig:"benchlarge.app.debug" default:"false"`
		Database DatabaseConfig `konfig:"benchlarge.database"`
		Server   ServerConfig   `konfig:"benchlarge.server"`
		Logging  LoggingConfig  `konfig:"benchlarge.logging"`
		Redis    RedisConfig    `konfig:"benchlarge.redis"`
	}

	// Set up some environment variables
	envVars := map[string]string{
		"benchlarge.app.name":      "bench-app",
		"benchlarge.database.host": "db.example.com",
		"benchlarge.database.port": "3306",
		"benchlarge.server.port":   "9090",
		"benchlarge.logging.level": "debug",
		"benchlarge.redis.host":    "redis.example.com",
	}

	for key, value := range envVars {
		os.Setenv(key, value)
	}

	defer func() {
		for key := range envVars {
			os.Unsetenv(key)
		}
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

func BenchmarkLoadInto_DeepNesting(b *testing.B) {
	type Level4Config struct {
		Value string `konfig:"value" default:"level4"`
	}

	type Level3Config struct {
		Level4 Level4Config `konfig:"level4"`
		Value  string       `konfig:"value" default:"level3"`
	}

	type Level2Config struct {
		Level3 Level3Config `konfig:"level3"`
		Value  string       `konfig:"value" default:"level2"`
	}

	type Level1Config struct {
		Level2 Level2Config `konfig:"level2"`
		Value  string       `konfig:"value" default:"level1"`
	}

	type Config struct {
		Deep Level1Config `konfig:"benchdeep.level1"`
	}

	os.Setenv("benchdeep.level1.level2.level3.level4.value", "deep_bench")
	defer os.Unsetenv("benchdeep.level1.level2.level3.level4.value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var cfg Config
		err := LoadInto(&cfg)
		if err != nil {
			b.Fatalf("LoadInto failed: %v", err)
		}
	}
}

func BenchmarkLoadInto_ManyFields(b *testing.B) {
	type Config struct {
		Field01 string `konfig:"benchmany.field01" default:"default01"`
		Field02 string `konfig:"benchmany.field02" default:"default02"`
		Field03 string `konfig:"benchmany.field03" default:"default03"`
		Field04 string `konfig:"benchmany.field04" default:"default04"`
		Field05 string `konfig:"benchmany.field05" default:"default05"`
		Field06 string `konfig:"benchmany.field06" default:"default06"`
		Field07 string `konfig:"benchmany.field07" default:"default07"`
		Field08 string `konfig:"benchmany.field08" default:"default08"`
		Field09 string `konfig:"benchmany.field09" default:"default09"`
		Field10 string `konfig:"benchmany.field10" default:"default10"`
		Field11 string `konfig:"benchmany.field11" default:"default11"`
		Field12 string `konfig:"benchmany.field12" default:"default12"`
		Field13 string `konfig:"benchmany.field13" default:"default13"`
		Field14 string `konfig:"benchmany.field14" default:"default14"`
		Field15 string `konfig:"benchmany.field15" default:"default15"`
		Field16 string `konfig:"benchmany.field16" default:"default16"`
		Field17 string `konfig:"benchmany.field17" default:"default17"`
		Field18 string `konfig:"benchmany.field18" default:"default18"`
		Field19 string `konfig:"benchmany.field19" default:"default19"`
		Field20 string `konfig:"benchmany.field20" default:"default20"`
	}

	// Set every other field to test mixed env/default performance
	for i := 1; i <= 20; i += 2 {
		key := "benchmany.field" + fmt.Sprintf("%02d", i)
		value := "env_value" + fmt.Sprintf("%02d", i)
		os.Setenv(key, value)
	}

	defer func() {
		for i := 1; i <= 20; i += 2 {
			key := "benchmany.field" + fmt.Sprintf("%02d", i)
			os.Unsetenv(key)
		}
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

// Memory allocation benchmarks
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
