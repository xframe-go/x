package auth

type Config struct {
	Secret     string
	Expiration int64 // hours
}
