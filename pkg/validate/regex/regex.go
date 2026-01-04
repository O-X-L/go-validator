package regex

import (
	"fmt"
	"regexp"
)

const (
	REGEX_UUID4                        = "^([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$"
	REGEX_DOMAINS_SIMPLE               = "^[a-zA-Z0-9\\-\\.]{1,253}\\.[a-zA-Z0-9\\-\\.]{1,63}$"
	REGEX_PATH_SIMPLE                  = "^(\\/|\\.)[\\/\\-_a-zA-Z0-9\\.]*$"
	REGEX_HASH_HEX                     = "^[a-f0-9]{6,50}$"
	REGEX_CERT_PUBLIC                  = "^-----BEGIN CERTIFICATE-----[-A-Za-z0-9+\\/\\n]*={0,3}\\n-----END CERTIFICATE-----$"
	REGEX_CERT_PRIVATE                 = "^(-----BEGIN EC PARAMETERS-----[-A-Za-z0-9+\\/\\n]*={0,3}\\n-----END EC PARAMETERS-----\\n-----BEGIN EC PRIVATE KEY-----[-A-Za-z0-9+\\/\n]*={0,3}\\n-----END EC PRIVATE KEY-----|-----BEGIN PRIVATE KEY-----[-A-Za-z0-9+\\/\\n]*={0,3}\\n-----END PRIVATE KEY-----)$"
	REGEX_URL_SIMPLE                   = "^https:\\/\\/[a-z0-9A-Z\\-\\._~:\\/?#\\[\\]@!$&'\\(\\)*+,;%=]*$"
	REGEX_ALPHANUMERIC                 = "^[a-zA-Z0-9]{1,}$"
	REGEX_ALPHANUMERIC_DASH_UNDERSCORE = "^[-a-zA-Z0-9_]{1,}$"
)

func ValidateRegex(regex string, value interface{}) bool {
	m, _ := regexp.MatchString(regex, fmt.Sprintf("%v", value))
	return m
}
