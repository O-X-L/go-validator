package custom

import (
	"net"
	"reflect"
	"strings"

	"git.oxl.at/go-validator/pkg/validate/regex"
)

func ValidateDomain(value interface{}) bool {
	s, ok := value.(string)
	if !ok {
		return false
	}
	return value == "localhost" || regex.ValidateRegex(regex.REGEX_DOMAINS_SIMPLE, s)
}

func ValidatePort(value interface{}) bool {
	val := reflect.ValueOf(value)
	var port int
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		port = int(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		port = int(val.Uint())
	default:
		return false
	}
	return port > 0 && port < 65536
}

func ValidateIP(value interface{}) bool {
	str, ok := value.(string)
	if !ok {
		return false
	}

	return net.ParseIP(str) != nil
}

func ValidateIPDomain(value interface{}) bool {
	if ValidateDomain(value) {
		return true
	}
	return ValidateIP(value)
}

func ValidateEmail(value interface{}) bool {
	email, ok := value.(string)
	if !ok {
		return false
	}
	if !strings.Contains(email, "@") {
		return false
	}
	if len(email) > 254 {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	return ValidateDomain(parts[1])
}

func ValidateNoDuplicates(value interface{}) bool {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return false
	}
	seen := make(map[interface{}]struct{})

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		if _, exists := seen[item]; exists {
			return true
		}
		seen[item] = struct{}{}
	}

	return false
}
