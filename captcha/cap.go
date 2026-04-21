package captcha

import gocap "github.com/ackcoder/go-cap"

func NewCaptcha() *gocap.Cap {
	return gocap.New(
		gocap.WithChallenge(50, 32, 3),
		gocap.WithChallengeExpires(60*2),
		gocap.WithTokenExpires(60*5),
	)
}
