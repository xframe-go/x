package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"sort"
	"strings"
)

func VerifyShopifyHMAC(params url.Values, hmacValue string, apiSecret string) bool {
	hmacGen := hmac.New(sha256.New, []byte(apiSecret))

	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "hmac" && k != "signature" {
			keys = append(keys, k)
		}
	}

	sort.Strings(keys)

	for _, k := range keys {
		hmacGen.Write([]byte(k))
		hmacGen.Write([]byte("="))
		hmacGen.Write([]byte(params.Get(k)))
	}

	calculatedHMAC := hex.EncodeToString(hmacGen.Sum(nil))
	return hmacEqual(calculatedHMAC, hmacValue)
}

func VerifyShopifyHMACLegacy(query string, hmacValue string, apiSecret string) bool {
	hmacGen := hmac.New(sha256.New, []byte(apiSecret))

	parts := strings.Split(query, "&")
	for _, part := range parts {
		if !strings.Contains(part, "hmac=") && !strings.Contains(part, "signature=") {
			hmacGen.Write([]byte(part))
		}
	}

	calculatedHMAC := hex.EncodeToString(hmacGen.Sum(nil))
	return hmacEqual(calculatedHMAC, hmacValue)
}

func hmacEqual(a, b string) bool {
	if len(a) != len(b) {
		return false
	}

	result := byte(0)
	for i := 0; i < len(a); i++ {
		result |= a[i] ^ b[i]
	}

	return result == 0
}
