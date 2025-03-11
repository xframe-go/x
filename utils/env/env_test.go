package env

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestString(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue string
		want         string
	}{
		{
			name:         "existing value",
			key:          "TEST_STRING",
			value:        "test value",
			defaultValue: "default",
			want:         "test value",
		},
		{
			name:         "empty value",
			key:          "TEST_STRING_EMPTY",
			value:        "",
			defaultValue: "default",
			want:         "default",
		},
		{
			name:         "non-existent key",
			key:          "TEST_STRING_NONEXISTENT",
			defaultValue: "default",
			want:         "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := String(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue int
		want         int
	}{
		{
			name:         "valid integer",
			key:          "TEST_INT",
			value:        "123",
			defaultValue: 0,
			want:         123,
		},
		{
			name:         "invalid integer",
			key:          "TEST_INT_INVALID",
			value:        "not a number",
			defaultValue: 42,
			want:         42,
		},
		{
			name:         "empty value",
			key:          "TEST_INT_EMPTY",
			value:        "",
			defaultValue: 42,
			want:         42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := Int(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("Int() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBool(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue bool
		want         bool
	}{
		{"true value", "TEST_BOOL", "true", false, true},
		{"TRUE value", "TEST_BOOL", "TRUE", false, true},
		{"1 value", "TEST_BOOL", "1", false, true},
		{"false value", "TEST_BOOL", "false", true, false},
		{"FALSE value", "TEST_BOOL", "FALSE", true, false},
		{"0 value", "TEST_BOOL", "0", true, false},
		{"invalid value", "TEST_BOOL", "invalid", true, true},
		{"empty value", "TEST_BOOL", "", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := Bool(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue time.Duration
		want         time.Duration
	}{
		{
			name:         "valid duration",
			key:          "TEST_DURATION",
			value:        "1h30m",
			defaultValue: time.Hour,
			want:         90 * time.Minute,
		},
		{
			name:         "invalid duration",
			key:          "TEST_DURATION_INVALID",
			value:        "invalid",
			defaultValue: time.Hour,
			want:         time.Hour,
		},
		{
			name:         "empty value",
			key:          "TEST_DURATION_EMPTY",
			value:        "",
			defaultValue: time.Hour,
			want:         time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := Duration(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		sep          string
		defaultValue []string
		want         []string
	}{
		{
			name:         "comma separated",
			key:          "TEST_SLICE",
			value:        "a,b,c",
			sep:          ",",
			defaultValue: []string{"default"},
			want:         []string{"a", "b", "c"},
		},
		{
			name:         "space separated",
			key:          "TEST_SLICE",
			value:        "a b c",
			sep:          " ",
			defaultValue: []string{"default"},
			want:         []string{"a", "b", "c"},
		},
		{
			name:         "empty parts filtered",
			key:          "TEST_SLICE",
			value:        "a,,b, ,c",
			sep:          ",",
			defaultValue: []string{"default"},
			want:         []string{"a", "b", "c"},
		},
		{
			name:         "empty value",
			key:          "TEST_SLICE_EMPTY",
			value:        "",
			sep:          ",",
			defaultValue: []string{"default"},
			want:         []string{"default"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			got := StringSlice(tt.key, tt.defaultValue, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue float64
		want         float64
	}{
		{
			name:         "valid float",
			key:          "TEST_FLOAT",
			value:        "123.456",
			defaultValue: 0.0,
			want:         123.456,
		},
		{
			name:         "integer as float",
			key:          "TEST_FLOAT",
			value:        "123",
			defaultValue: 0.0,
			want:         123.0,
		},
		{
			name:         "invalid float",
			key:          "TEST_FLOAT_INVALID",
			value:        "not a number",
			defaultValue: 42.0,
			want:         42.0,
		},
		{
			name:         "empty value",
			key:          "TEST_FLOAT_EMPTY",
			value:        "",
			defaultValue: 42.0,
			want:         42.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := Float64(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		defaultValue int64
		want         int64
	}{
		{
			name:         "valid int64",
			key:          "TEST_INT64",
			value:        "9223372036854775807", // max int64
			defaultValue: 0,
			want:         9223372036854775807,
		},
		{
			name:         "invalid int64",
			key:          "TEST_INT64_INVALID",
			value:        "not a number",
			defaultValue: 42,
			want:         42,
		},
		{
			name:         "empty value",
			key:          "TEST_INT64_EMPTY",
			value:        "",
			defaultValue: 42,
			want:         42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
				defer os.Unsetenv(tt.key)
			}
			if got := Int64(tt.key, tt.defaultValue); got != tt.want {
				t.Errorf("Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}
