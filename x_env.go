package x

// 环境变量常量定义
var (
	// EnvDevelopment 开发环境
	EnvDevelopment = "development"
	// EnvTesting 测试环境
	EnvTesting = "testing"
	// EnvProduction 生产环境
	EnvProduction = "production"
)

// Env 返回当前应用程序的运行环境
func Env() string {
	return xApp.Get().AppConfig().Env
}

// IsDevelopment 判断当前是否为开发环境
func IsDevelopment() bool {
	return Env() == EnvDevelopment
}

// IsTesting 判断当前是否为测试环境
func IsTesting() bool {
	return Env() == EnvTesting
}

// IsProduction 判断当前是否为生产环境
func IsProduction() bool {
	return Env() == EnvProduction
}
