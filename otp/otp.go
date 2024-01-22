package otp

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

func GenerateKey(issuer string, accountName string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountName,
	})
}
func GeneratePasscode(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}

func Validate(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}
