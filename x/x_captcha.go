package x

import (
	gocap "github.com/ackcoder/go-cap"
)

func RegisterCaptcha(fn func() *gocap.Cap) {
	rocket.captcha = fn()
}

func Captcha() *gocap.Cap {
	return rocket.captcha
}
