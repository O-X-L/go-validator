package test

import (
	"log"
	"slices"
	"strings"
	"testing"

	"git.oxl.at/go-validator/pkg/validate"
)

type DummyApp struct {
	Description string   `validate_regex:"^[-a-z_.:;\\s]*$"`        // tag-specific regex
	Domains     []string `validate:"domain" required:"true"`        // length > 0; every entry is validated
	Enabled     *bool    `required:"true"`                          // Tip: *bool allows for a null-value
	ListenIP    string   `validate:"ip" required_if:"Enabled=true"` // conditionally required
	Options     []string `required:"true" validate:"my_custom_validator"`
}

type MyData struct {
	Name string     `validate:"alphanumeric"`
	Apps []DummyApp `required:"true"` // length > 0;
}

func TestExample(t *testing.T) {
	v := &validate.StructValidator{}
	v.ValidatorsCustom = validate.GetDefaultCustomValidators()
	v.ValidatorsRegex = validate.GetDefaultRegexValidators()
	/*
		or omit/replace the defaults:
		v.ValidatorsCustom = map[string]validate.ValidatorCustom{}
		v.ValidatorsRegex = map[string]string{}
	*/

	v.ValidatorsCustom["my_custom_validator"] = func(value interface{}) bool {
		// should be key=value pairs
		str, ok := value.(string)
		if !ok {
			return false
		}
		if !strings.Contains(str, "=") {
			return false
		}
		parts := strings.Split(str, "=")
		if len(parts) > 2 {
			return false
		}
		if !slices.Contains([]string{"opt1", "opt2"}, parts[0]) {
			// invalid option
			return false
		}
		return true
	}

	data := MyData{
		Apps: []DummyApp{
			{
				Domains: []string{"test.oxl.at"},
				Options: []string{"opt1=test"},
			},
		},
	}
	errors := v.Validate(data)
	if len(errors) > 0 {
		log.Fatalf("Data validation failed: %v\n", errors)
	}
}
