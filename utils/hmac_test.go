package utils

import (
	"net/url"
	"testing"
)

func TestVerifyShopifyHMAC(t *testing.T) {
	apiSecret := "my_secret_key"

	params := url.Values{}
	params.Set("code", "test_code")
	params.Set("shop", "test-shop.myshopify.com")
	params.Set("state", "random_state")

	hmacValue := "calculated_hmac"

	result := VerifyShopifyHMAC(params, hmacValue, apiSecret)

	if !result {
		t.Errorf("HMAC verification failed unexpectedly")
	}
}

func TestHmacEqual(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "equal strings",
			a:    "hello",
			b:    "hello",
			want: true,
		},
		{
			name: "different strings",
			a:    "hello",
			b:    "world",
			want: false,
		},
		{
			name: "different lengths",
			a:    "hello",
			b:    "hello!",
			want: false,
		},
		{
			name: "empty strings",
			a:    "",
			b:    "",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hmacEqual(tt.a, tt.b); got != tt.want {
				t.Errorf("hmacEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHmacEqualTimingSafety(t *testing.T) {
	a := "test_string_with_some_length"
	b := "test_string_with_some_length"
	c := "different_string_with_length"

	result1 := hmacEqual(a, b)
	result2 := hmacEqual(a, c)

	if !result1 {
		t.Error("Expected equal HMACs")
	}

	if result2 {
		t.Error("Expected different HMACs")
	}
}
