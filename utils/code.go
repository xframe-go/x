package utils

import (
	"fmt"
	"strings"

	"github.com/dromara/carbon/v2"
	"github.com/xframe-go/x/snowflake"
)

var codeGenerator = snowflake.New()

// GenerateCode 生成唯一编码
func GenerateCode() string {
	return codeGenerator.Generate()
}

// GenerateCodeWithPrefix 生成带前缀的唯一编码
func GenerateCodeWithPrefix(prefix string) string {
	return prefix + "_" + codeGenerator.Generate()
}

// SetCodeIfEmpty 如果 code 为空，则自动生成
func SetCodeIfEmpty(code *string) {
	if strings.TrimSpace(*code) == "" {
		*code = GenerateCode()
	}
}

// SetCodeWithPrefix 如果 code 为空，则使用指定前缀生成
func SetCodeWithPrefix(code *string, prefix string) {
	if strings.TrimSpace(*code) == "" {
		*code = GenerateCodeWithPrefix(prefix)
	}
}

func GenCodeWithPrefix(prefix string) string {
	now := carbon.Now()
	return fmt.Sprintf("%s-%02d%d-%04d", prefix, now.Month(), now.Day(), now.Microsecond()%10000)
}
