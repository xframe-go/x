// Package env 提供了一个用于获取环境变量的工具包
// 支持多种数据类型，并提供默认值支持
package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// String 获取字符串类型的环境变量
// 如果环境变量不存在，则返回默认值
func String(key string, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// Int 获取整数类型的环境变量
// 如果环境变量不存在或解析失败，则返回默认值
func Int(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// Int64 获取 int64 类型的环境变量
// 如果环境变量不存在或解析失败，则返回默认值
func Int64(key string, defaultValue int64) int64 {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i
		}
	}
	return defaultValue
}

// Float64 获取浮点数类型的环境变量
// 如果环境变量不存在或解析失败，则返回默认值
func Float64(key string, defaultValue float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return defaultValue
}

// Bool 获取布尔类型的环境变量
// 支持的真值：1, t, T, true, TRUE, True
// 支持的假值：0, f, F, false, FALSE, False
// 如果环境变量不存在或解析失败，则返回默认值
func Bool(key string, defaultValue bool) bool {
	if v := os.Getenv(key); v != "" {
		v = strings.ToLower(v)
		switch v {
		case "1", "t", "true":
			return true
		case "0", "f", "false":
			return false
		}
	}
	return defaultValue
}

// Duration 获取时间间隔类型的环境变量
// 支持的格式：time.ParseDuration 支持的所有格式
// 如果环境变量不存在或解析失败，则返回默认值
func Duration(key string, defaultValue time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultValue
}

// StringSlice 获取字符串切片类型的环境变量
// 使用指定的分隔符分割字符串
// 如果环境变量不存在，则返回默认值
func StringSlice(key string, defaultValue []string, sep string) []string {
	if v := os.Getenv(key); v != "" {
		parts := strings.Split(v, sep)
		// 去除空字符串
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return defaultValue
}
