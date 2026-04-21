package env

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

var (
	envVars = make(map[string]string)
	loaded  = false
)

func Load(filenames ...string) error {
	if loaded {
		return nil
	}

	if len(filenames) == 0 {
		filenames = []string{".env"}
	}

	for _, filename := range filenames {
		if err := loadFile(filename); err != nil {
			// 文件不存在不是错误
			if !os.IsNotExist(err) {
				return err
			}
		}
	}

	loaded = true
	return nil
}

func loadFile(filename string) error {
	// 获取项目根目录
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	// 从 pkg/env 回到项目根目录
	rootPath := filepath.Join(basePath, "..", "..")
	filepath := filepath.Join(rootPath, filename)

	file, err := os.Open(filepath)
	if err != nil {
		// 尝试当前目录
		file, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 去除引号
		value = strings.Trim(value, `"'`)

		envVars[key] = value
	}

	return scanner.Err()
}

func get(key string) (string, bool) {
	if !loaded {
		_ = Load()
	}

	// 优先从系统环境变量获取
	if value, exists := os.LookupEnv(key); exists {
		return value, true
	}

	// 从.env文件获取
	value, exists := envVars[key]
	return value, exists
}

func String(key string, defaultValue ...string) string {
	value, exists := get(key)
	if !exists && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return value
}

func Int(key string, defaultValue ...int) int {
	value, exists := get(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return i
}

func Int64(key string, defaultValue ...int64) int64 {
	value, exists := get(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return i
}

func Bool(key string, defaultValue ...bool) bool {
	value, exists := get(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return false
	}
	return b
}

func Float64(key string, defaultValue ...float64) float64 {
	value, exists := get(key)
	if !exists {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}
	return f
}

func MustString(key string) (string, error) {
	value, exists := get(key)
	if !exists {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return value, nil
}

func MustInt(key string) (int, error) {
	value, exists := get(key)
	if !exists {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}
	return strconv.Atoi(value)
}

func Has(key string) bool {
	_, exists := get(key)
	return exists
}
