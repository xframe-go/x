package db

import "fmt"

type (
	WithConnection interface {
		Connection() string
	}

	Config struct {
		Databases map[string]DriverConf
	}

	DriverConf struct {
		Driver   string
		Host     string
		Port     uint
		Username string
		Password string
		DB       string
		Charset  string
		Debug    bool
	}
)

func (d DriverConf) Dsn() (string, error) {
	if d.Driver == "mysql" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", d.Username, d.Password, d.Host, d.Port, d.DB, d.Charset), nil
	}
	return "", fmt.Errorf("unsupported driver: %s", d.Driver)
}
