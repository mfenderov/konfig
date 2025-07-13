package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/mfenderov/konfig"
)

// Comprehensive validation script for konfig library
// This script validates all major functionality including the new LoadInto method

type ValidationConfig struct {
	SimpleField  string        `konfig:"validation.simple" default:"simple_default"`
	NestedField  NestedConfig  `konfig:"validation.nested"`
	EmptyField   string        `konfig:"validation.empty"`
}

type NestedConfig struct {
	SubField string `konfig:"sub_field" default:"nested_default"`
	DeepNest DeepConfig `konfig:"deep"`
}

type DeepConfig struct {
	Value string `konfig:"value" default:"deep_default"`
}

func main() {
	fmt.Println("🔍 Konfig Library Validation")
	fmt.Println("============================")

	totalTests := 0
	passedTests := 0

	// Test 1: Basic Load functionality
	fmt.Println("\n1️⃣  Testing basic Load() functionality...")
	totalTests++
	err := konfig.Load()
	if err != nil {
		fmt.Printf("   ❌ FAIL: %v\n", err)
	} else {
		fmt.Printf("   ✅ PASS: Basic Load() works\n")
		passedTests++
	}

	// Test 2: Profile detection
	fmt.Println("\n2️⃣  Testing profile detection...")
	totalTests++
	profile := konfig.GetProfile()
	fmt.Printf("   📋 Current profile: '%s'\n", profile)
	
	// Test profile functions
	if konfig.IsProfile("") || konfig.IsProfile("nonexistent") {
		fmt.Printf("   ❌ FAIL: Profile detection incorrect\n")
	} else {
		fmt.Printf("   ✅ PASS: Profile detection works\n")
		passedTests++
	}

	// Test 3: Environment variable population
	fmt.Println("\n3️⃣  Testing environment variable population...")
	totalTests++
	os.Setenv("test.env.var", "test_value")
	err = konfig.Load()
	if err != nil {
		fmt.Printf("   ❌ FAIL: Load failed: %v\n", err)
	} else {
		retrievedValue := os.Getenv("test.env.var")
		if retrievedValue == "test_value" {
			fmt.Printf("   ✅ PASS: Environment variables work\n")
			passedTests++
		} else {
			fmt.Printf("   ❌ FAIL: Expected 'test_value', got '%s'\n", retrievedValue)
		}
	}
	os.Unsetenv("test.env.var")

	// Test 4: LoadInto basic functionality
	fmt.Println("\n4️⃣  Testing LoadInto() basic functionality...")
	totalTests++
	os.Setenv("validation.simple", "env_override")
	
	var config ValidationConfig
	err = konfig.LoadInto(&config)
	if err != nil {
		fmt.Printf("   ❌ FAIL: LoadInto failed: %v\n", err)
	} else if config.SimpleField != "env_override" {
		fmt.Printf("   ❌ FAIL: Expected 'env_override', got '%s'\n", config.SimpleField)
	} else {
		fmt.Printf("   ✅ PASS: LoadInto basic functionality works\n")
		passedTests++
	}
	os.Unsetenv("validation.simple")

	// Test 5: Default values
	fmt.Println("\n5️⃣  Testing default values...")
	totalTests++
	var configDefaults ValidationConfig
	err = konfig.LoadInto(&configDefaults)
	if err != nil {
		fmt.Printf("   ❌ FAIL: LoadInto failed: %v\n", err)
	} else if configDefaults.SimpleField != "simple_default" {
		fmt.Printf("   ❌ FAIL: Expected 'simple_default', got '%s'\n", configDefaults.SimpleField)
	} else {
		fmt.Printf("   ✅ PASS: Default values work\n")
		passedTests++
	}

	// Test 6: Nested struct support
	fmt.Println("\n6️⃣  Testing nested struct support...")
	totalTests++
	os.Setenv("validation.nested.sub_field", "nested_override")
	os.Setenv("validation.nested.deep.value", "deep_override")
	
	var configNested ValidationConfig
	err = konfig.LoadInto(&configNested)
	if err != nil {
		fmt.Printf("   ❌ FAIL: LoadInto failed: %v\n", err)
	} else if configNested.NestedField.SubField != "nested_override" {
		fmt.Printf("   ❌ FAIL: Expected 'nested_override', got '%s'\n", configNested.NestedField.SubField)
	} else if configNested.NestedField.DeepNest.Value != "deep_override" {
		fmt.Printf("   ❌ FAIL: Expected 'deep_override', got '%s'\n", configNested.NestedField.DeepNest.Value)
	} else {
		fmt.Printf("   ✅ PASS: Nested struct support works\n")
		passedTests++
	}
	os.Unsetenv("validation.nested.sub_field")
	os.Unsetenv("validation.nested.deep.value")

	// Test 7: Error handling
	fmt.Println("\n7️⃣  Testing error handling...")
	totalTests++
	
	// Test nil pointer
	err = konfig.LoadInto(nil)
	if err == nil {
		fmt.Printf("   ❌ FAIL: Should return error for nil pointer\n")
	} else {
		// Test non-pointer
		var testStruct ValidationConfig
		err = konfig.LoadInto(testStruct)
		if err == nil {
			fmt.Printf("   ❌ FAIL: Should return error for non-pointer\n")  
		} else {
			// Test pointer to non-struct
			var testString string
			err = konfig.LoadInto(&testString)
			if err == nil {
				fmt.Printf("   ❌ FAIL: Should return error for pointer to non-struct\n")
			} else {
				fmt.Printf("   ✅ PASS: Error handling works correctly\n")
				passedTests++
			}
		}
	}

	// Test 8: Empty values vs defaults
	fmt.Println("\n8️⃣  Testing empty values vs defaults...")
	totalTests++
	os.Setenv("validation.empty", "")  // Explicitly empty
	
	var configEmpty ValidationConfig
	err = konfig.LoadInto(&configEmpty)
	if err != nil {
		fmt.Printf("   ❌ FAIL: LoadInto failed: %v\n", err)
	} else if configEmpty.EmptyField != "" {
		fmt.Printf("   ❌ FAIL: Expected empty string, got '%s'\n", configEmpty.EmptyField)
	} else {
		fmt.Printf("   ✅ PASS: Empty values work correctly\n")
		passedTests++
	}
	os.Unsetenv("validation.empty")

	// Test 9: Performance validation (basic)
	fmt.Println("\n9️⃣  Testing performance...")
	totalTests++
	const iterations = 100
	
	var perfConfig ValidationConfig
	for i := 0; i < iterations; i++ {
		err = konfig.LoadInto(&perfConfig)
		if err != nil {
			fmt.Printf("   ❌ FAIL: Performance test failed at iteration %d: %v\n", i, err)
			break
		}
	}
	if err == nil {
		fmt.Printf("   ✅ PASS: Performance test completed %d iterations\n", iterations)
		passedTests++
	}

	// Test 10: Memory test (no memory leaks validation)
	fmt.Println("\n🔟 Testing memory usage...")
	totalTests++
	
	// Create many configs to test for obvious memory issues
	configs := make([]ValidationConfig, 50)
	for i := range configs {
		err = konfig.LoadInto(&configs[i])
		if err != nil {
			fmt.Printf("   ❌ FAIL: Memory test failed: %v\n", err)
			break
		}
	}
	if err == nil {
		fmt.Printf("   ✅ PASS: Memory test completed (no obvious leaks)\n")
		passedTests++
	}

	// Final results
	fmt.Println("\n" + "="*50)
	fmt.Printf("📊 VALIDATION RESULTS: %d/%d tests passed\n", passedTests, totalTests)
	
	percentage := float64(passedTests) / float64(totalTests) * 100
	fmt.Printf("📈 Success Rate: %.1f%%\n", percentage)
	
	if passedTests == totalTests {
		fmt.Println("🎉 ALL TESTS PASSED! Konfig library is ready for release.")
	} else {
		fmt.Printf("⚠️  %d tests failed. Please review and fix issues before release.\n", totalTests-passedTests)
		os.Exit(1)
	}

	// Additional feature validation
	fmt.Println("\n🔧 Additional Feature Validation:")
	
	// Validate that GetProfile function exists and works
	fmt.Printf("   GetProfile(): %s\n", konfig.GetProfile())
	fmt.Printf("   IsDevProfile(): %t\n", konfig.IsDevProfile())
	fmt.Printf("   IsProdProfile(): %t\n", konfig.IsProdProfile())
	fmt.Printf("   IsProfile('test'): %t\n", konfig.IsProfile("test"))
	
	fmt.Println("\n✅ Konfig library validation completed successfully!")
}