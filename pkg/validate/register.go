package validate

import (
	"git.oxl.at/go-validator/pkg/validate/custom"
	"git.oxl.at/go-validator/pkg/validate/regex"
)

type ValidatorCustom func(value interface{}) bool

func GetDefaultRegexValidators() map[string]string {
	v := make(map[string]string)

	v["uuid4"] = regex.REGEX_UUID4
	v["path_simple"] = regex.REGEX_PATH_SIMPLE
	v["cert_public"] = regex.REGEX_CERT_PUBLIC
	v["cert_private"] = regex.REGEX_CERT_PRIVATE
	v["url_simple"] = regex.REGEX_URL_SIMPLE
	v["alphanumeric"] = regex.REGEX_ALPHANUMERIC
	v["alphanumeric_dash_underscore"] = regex.REGEX_ALPHANUMERIC_DASH_UNDERSCORE

	return v
}

func GetDefaultCustomValidators() map[string]ValidatorCustom {
	v := make(map[string]ValidatorCustom)

	v["domain"] = custom.ValidateDomain
	v["ip_or_domain"] = custom.ValidateIPDomain
	v["ip"] = custom.ValidateIP
	v["port"] = custom.ValidatePort
	v["email"] = custom.ValidateEmail
	v["no_duplicates"] = custom.ValidateEmail

	return v
}
