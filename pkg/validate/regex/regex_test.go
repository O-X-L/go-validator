package regex

import (
	"strings"
	"testing"
)

func TestValidateRegexURLSimple(t *testing.T) {
	if ValidateRegex(REGEX_URL_SIMPLE, "1") {
		t.Error("URL validation failed #1")
	}
	if ValidateRegex(REGEX_URL_SIMPLE, "a") {
		t.Error("URL validation failed #2")
	}
	if !ValidateRegex(REGEX_URL_SIMPLE, "https://oxl.at") {
		t.Error("URL validation failed #3")
	}
	if !ValidateRegex(REGEX_URL_SIMPLE, "https://github.com/O-X-L/micro-queue/blob/latest/cmd/main.go") {
		t.Error("URL validation failed #4")
	}
	if !ValidateRegex(REGEX_URL_SIMPLE, "https://demo.ansible-webui.oxl.app/ui#jobs&edit=2") {
		t.Error("URL validation failed #5")
	}
	if !ValidateRegex(REGEX_URL_SIMPLE, "https://python-opnsense.oxl.app/usage/4_multi.html") {
		t.Error("URL validation failed #6")
	}
}

func TestValidatePathSimple(t *testing.T) {
	if ValidateRegex(REGEX_PATH_SIMPLE, "1") {
		t.Error("Path validation failed #1")
	}
	if ValidateRegex(REGEX_PATH_SIMPLE, "a") {
		t.Error("Path validation failed #2")
	}
	if !ValidateRegex(REGEX_PATH_SIMPLE, "/var/tmp") {
		t.Error("Path validation failed #3")
	}
	if !ValidateRegex(REGEX_PATH_SIMPLE, "/home/user/.local/lib/test") {
		t.Error("Path validation failed #4")
	}
	if !ValidateRegex(REGEX_PATH_SIMPLE, ".local/lib/test") {
		t.Error("Path validation failed #5")
	}
}

var CERT_PK_RSA = strings.ReplaceAll("-----BEGIN PRIVATE KEY-----\\n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\\nv82eLLqtEfFwbQmYXGzatw8=\\n-----END PRIVATE KEY-----", `\n`, "\n")
var CERT_PK_EC = strings.ReplaceAll("-----BEGIN EC PARAMETERS-----\\nBggqhkjOPQMBBw==\\n-----END EC PARAMETERS-----\\n-----BEGIN EC PRIVATE KEY-----\\nMHcCAQEEIMYm45kGcui2bQvbkXEGLY4jsHZlT81+UO8DDSQt+MbyoAoGCCqGSM49\\nAwEHoUQDQgAEggV7rnUyxyTCm+UbfB22zeyDJicphgetV6rAjLLMHVNeMVRDtxgj\\nm1wSIy1IDmmeeDZ6tpdEnROx4EpRfbT0ng==\\n-----END EC PRIVATE KEY-----", `\n`, "\n")
var CERT_PUB_EC = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\\nMIIBtjCCAV2gAwIBAgIUCoj8G7BMKZAK6W6QQH8B1hG9o3QwCgYIKoZIzj0EAwIw\\nFzEVMBMGA1UEAwwMSEFQcm94eSBUZXN0MB4XDTI1MTEwOTIxMDMwOVoXDTM1MTEw\\nNzIxMDMwOVowFzEVMBMGA1UEAwwMSEFQcm94eSBUZXN0MFkwEwYHKoZIzj0CAQYI\\nKoZIzj0DAQcDQgAEjWWmjG74753Q+R03UpqRm6e3kxbgCfEge7k4EZia1eqZkQrK\\nssbHYlRsH+5ahLRNcaDHn2XrqkQfjq/abs3jGaOBhjCBgzAdBgNVHQ4EFgQU+DD/\\nBgnioWKyaDldZGqcTgFfmiswHwYDVR0jBBgwFoAU+DD/BgnioWKyaDldZGqcTgFf\\nmiswDwYDVR0TAQH/BAUwAwEB/zAwBgNVHREEKTAnghRhcHAudGVtcGxhdGUub3hs\\nLmFwcIIJbG9jYWxob3N0hwR/AAABMAoGCCqGSM49BAMCA0cAMEQCIBQE6FWocj9+\\nK9Gh543vxIkZPpZB17kQ0dy2DtCxQJWHAiB7vMvDH/EA9g4BbfheFbJHOWOFZeSy\\njFEIpZn9wF9COQ==\\n-----END CERTIFICATE-----", `\n`, "\n")
var CERT_PUB_RSA = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\\nMIIEQzCCAqugAwIBAgIUGJM6pFqsRWwm7zZks8d9Ue65i2swDQYJKoZIhvcNAQEL\\nBQAwFzEVMBMGA1UEAwwMSEFQcm94eSBUZXN0MB4XDTI1MTEwOTIxMDUxN1oXDTM1\\nMTEwNzIxMDUxN1owFzEVMBMGA1UEAwwMSEFQcm94eSBUZXN0MIIBojANBgkqhkiG\\n9w0BAQEFAAOCAY8AMIIBigKCAYEAwg+pFi69l+9bTDsVxBxA+iQwhdfmSH2CEBaY\\ntI0/8N2HVvYM90dLjhl1HeU8mnfvVJrxGYgXGkx9igBcBGtApBYTpZ/Zu8POmScb\\nePBoGCyc/CwVbUATVqj/6saO/ipX70c6QVkMl59NisDF1F8m2BEMUsUvT2bFqoGu\\n6SwA/zjJmgXAWd1EdQNsZLzvMDSgD4xDLWSp98GrLBHOHtj5X9bwtDScDDyG0+Cc\\nQFV2sFykz+XaGqPnnIAf8gvXYyrwoFctIh5eooUNF2g2pZGjnvvmrLM4UKTftJSB\\nV6ewCjMBvissNLZdH9oHlLXzvYrlwVoSNgonVylnyBDHJayz3PRFRBXD+EFRRBQP\\nea3jzAgLSeOFvvlBldYfxg91hV7tXlUDZHOor/8b3IRpybn7uE0wK/1jG2zF/lAv\\nycDKOBzslkuiNNdJmAuBwbUbl8HQk5aoDBC1P3o9J5qaiMujrgr8Z/MK/VCmvxID\\nJlK+RoTre1686hgMMU6X+NhCtNxpAgMBAAGjgYYwgYMwHQYDVR0OBBYEFF6wQNNy\\nwPpFBzme7MrXDAMp29rIMB8GA1UdIwQYMBaAFF6wQNNywPpFBzme7MrXDAMp29rI\\nMA8GA1UdEwEB/wQFMAMBAf8wMAYDVR0RBCkwJ4IUYXBwLnRlbXBsYXRlLm94bC5h\\ncHCCCWxvY2FsaG9zdIcEfwAAATANBgkqhkiG9w0BAQsFAAOCAYEAaooRle4KxTBj\\nkU4MFJK2HQ6MChQbKEa40pYyOasig9X+ei4nsxdu8Lqjmoq58W3jejKS6wnYHP2o\\nVk6LPU70Rk1DeRyRiO10D+IawRi+yhJMuxAsZnBqCPDJP7BFfbBDlFckYoGPgQQJ\\nXTyltK5lQ7WgNklpvtEgzBlAcvn5JdEHG0oee1xIielfQ7QFhZvolr0yrIJA66jk\\n3da0SUf/QReYRFAyI1GlbkgD/uG6kuftjCnmrWwglt9eiasu6V5lgw7JU4PbA5Y0\\nDCTb/ug/a++evttWNeWg5TeVVDQflpFiiNiS8OKlQGHtqROuJOEt9ibI8I0KUL18\\neI2nh+Urbzsm/iYzkQ1LpGotLmvcZMkOqiqbQiT0eZadL2AQYqQEsF+4/WNLnDyC\\n9n4Z04xPTfq/W0JkqHsegT2Rvt6DJ0gPSDki2ujJiqWOxRuoq8dco7Jf9s2ZqhuS\\nuyHZPSqXHYZO6ZEAVO55Azh6hdDb+QZedDyBFq2R71IuibxLSyMz\\n-----END CERTIFICATE-----", `\n`, "\n")

func TestValidateCertPrivate(t *testing.T) {
	if !ValidateRegex(REGEX_CERT_PRIVATE, CERT_PK_RSA) {
		t.Error("Cert-Private validation failed #1")
	}
	if ValidateRegex(REGEX_CERT_PRIVATE, CERT_PUB_EC) {
		t.Error("Cert-Private validation failed #2")
	}
	if !ValidateRegex(REGEX_CERT_PRIVATE, CERT_PK_EC) {
		t.Error("Cert-Private validation failed #3")
	}
	if ValidateRegex(REGEX_CERT_PRIVATE, CERT_PUB_RSA) {
		t.Error("Cert-Private validation failed #4")
	}
}

func TestValidateCertPublic(t *testing.T) {
	if ValidateRegex(REGEX_CERT_PUBLIC, CERT_PK_RSA) {
		t.Error("Cert-Public validation failed #1")
	}
	if !ValidateRegex(REGEX_CERT_PUBLIC, CERT_PUB_EC) {
		t.Error("Cert-Public validation failed #2")
	}
	if ValidateRegex(REGEX_CERT_PUBLIC, CERT_PK_EC) {
		t.Error("Cert-Public validation failed #3")
	}
	if !ValidateRegex(REGEX_CERT_PUBLIC, CERT_PUB_RSA) {
		t.Error("Cert-Public validation failed #4")
	}
}

func TestValidateUUID4(t *testing.T) {
	tests := []struct {
		uuid string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"not-a-uuid", false},
		{"550e8400-e29b-41d4-a716", false}, // too short
	}
	for _, tt := range tests {
		if ValidateRegex(REGEX_UUID4, tt.uuid) != tt.want {
			t.Errorf("UUID validation failed for %s", tt.uuid)
		}
	}
}

func TestValidateHex(t *testing.T) {
	if !ValidateRegex(REGEX_HASH_HEX, "a1b2c3d4e5f6") {
		t.Error("Hash validation failed valid hex")
	}
	if ValidateRegex(REGEX_HASH_HEX, "not-hex-!!") {
		t.Error("Hash validation failed to reject invalid hex")
	}
}

func TestValidateDomainsSimple(t *testing.T) {
	tests := []struct {
		domain string
		want   bool
	}{
		{"example.com", true},
		{"sub.domain.org", true},
		{"my-server.io", true},
		// {"-bad.com", false},
		{"invalid_char!.com", false},
		{"too-long." + strings.Repeat("a", 64), false}, // TLD segment > 63
	}
	for _, tt := range tests {
		if ValidateRegex(REGEX_DOMAINS_SIMPLE, tt.domain) != tt.want {
			t.Errorf("REGEX_DOMAINS_SIMPLE validation failed for %s", tt.domain)
		}
	}
}
