package env

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestEnv(t *testing.T) {
	// 清理环境
	envVars = make(map[string]string)
	loaded = false

	// 创建临时.env文件
	content := `TEST_VAR=from_file
TEST_INT=42
TEST_BOOL=true
TEST_FLOAT=3.14
# 这是注释
EMPTY_VAR=
QUOTED_VAR="quoted value"
`

	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test .env file: %v", err)
	}

	// 加载测试文件
	if err := Load(envFile); err != nil {
		t.Fatalf("failed to load test .env file: %v", err)
	}
}

func TestString(t *testing.T) {
	setupTestEnv(t)

	tests := []struct {
		key          string
		defaultValue string
		want         string
	}{
		{"TEST_VAR", "", "from_file"},
		{"NON_EXISTENT", "default", "default"},
		{"EMPTY_VAR", "default", ""},
		{"QUOTED_VAR", "", "quoted value"},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := String(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("String(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	setupTestEnv(t)

	tests := []struct {
		key          string
		defaultValue int
		want         int
	}{
		{"TEST_INT", 0, 42},
		{"NON_EXISTENT", 100, 100},
		{"TEST_VAR", 0, 0}, // "from_file" 无法转换为int
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := Int(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("Int(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	setupTestEnv(t)

	tests := []struct {
		key          string
		defaultValue bool
		want         bool
	}{
		{"TEST_BOOL", false, true},
		{"NON_EXISTENT", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := Bool(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("Bool(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	setupTestEnv(t)

	got := Float64("TEST_FLOAT", 0.0)
	if got != 3.14 {
		t.Errorf("Float64(TEST_FLOAT) = %v, want 3.14", got)
	}
}

func TestEnvironmentVariableOverride(t *testing.T) {
	setupTestEnv(t)

	// 设置环境变量（应该覆盖.env文件）
	os.Setenv("TEST_VAR", "from_env")
	defer os.Unsetenv("TEST_VAR")

	got := String("TEST_VAR", "")
	if got != "from_env" {
		t.Errorf("String(TEST_VAR) = %v, want from_env", got)
	}
}

func TestHas(t *testing.T) {
	setupTestEnv(t)

	if !Has("TEST_VAR") {
		t.Error("Has(TEST_VAR) should return true")
	}

	if Has("NON_EXISTENT_VAR") {
		t.Error("Has(NON_EXISTENT_VAR) should return false")
	}
}

func TestMustString(t *testing.T) {
	setupTestEnv(t)

	value, err := MustString("TEST_VAR")
	if err != nil {
		t.Errorf("MustString(TEST_VAR) error = %v", err)
	}
	if value != "from_file" {
		t.Errorf("MustString(TEST_VAR) = %v, want from_file", value)
	}

	_, err = MustString("NON_EXISTENT")
	if err == nil {
		t.Error("MustString(NON_EXISTENT) should return error")
	}
}
