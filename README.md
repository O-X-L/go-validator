# Dynamic Go Struct Value-Validation

[![Lint](https://github.com/O-X-L/go-validator/actions/workflows/lint.yml/badge.svg?branch=latest)](https://github.com/O-X-L/go-validator/actions/workflows/lint.yml)
[![Unit Tests](https://github.com/O-X-L/go-validator/actions/workflows/unit_test.yml/badge.svg?branch=latest)](https://github.com/O-X-L/go-validator/actions/workflows/unit_test.yml)

This library can validate complex Go struct's by utilizing struct-tags.

## Why

We, like most, have to validate struct-values after loading data - like a config-file or an API-call.

To keep the validation-logic clean we want to define those validators inside the struct-tags.

Its usage is similar to [creasty/defaults](github.com/creasty/defaults) - which can be used to set default values.

----

## Features

* **How we can validate**:
  * Custom validation-functions
  * Pre-defined regex
  * Tag-specific regex
  * Require a value to be set
  * Conditionally require a value to be set (*only very simple conditions for now*)

* **What we can validate**:
  * Single values
  * Every entry inside an array/slice
  * Nested structs
  * Arrays/slices of structs

----

## Examples

See [register.go](https://github.com/O-X-L/go-validator/blob/latest/pkg/validate/register.go) for all built-in validators.

### Basics

You are able to completely override the built-in validators.

See also: [examples](https://github.com/O-X-L/go-validator/blob/latest/examples)

```go
import (
	"log"

	"git.oxl.at/go-validator/pkg/validate"
)

type DummyApp struct {
    // tag-specific regex
	Description string   `validate_regex:"^[-a-z_.:;\\s]*$"`

    // length > 0; every entry is validated
	Domains     []string `validate:"domain" required:"true"`

	// Tip: *bool allows for a null-value
	Enabled     *bool    `required:"true"`

	// conditionally required
	ListenIP    string   `validate:"ip" required_if:"Enabled=true"`

	Options     []string `required:"true" validate:"my_custom_validator"`
}

type MyData struct {
	Name string     `validate:"alphanumeric"`
	Apps []DummyApp `required:"true"` // length > 0;
}

func TestExample() {
	v := &validate.StructValidator{}
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

```

### Loading a config-file

Example struct config:

<details>

```go
type AppCert struct {
	ID             uint              `yaml:"id" required:"true"`
	Provider       string            `yaml:"provider" required:"true"` // URL if HTTP-01 else one of the listed DNS-providers
	ProviderConfig map[string]string `yaml:"provider_config"`          // env-vars for DNS-01
	ChallengeType  string            `yaml:"challenge_type" validate_regex:"^(http-01|dns-01)$" required:"true"`
	Domains        []string          `yaml:"domains" validate:"domain"`
}

type App struct {
	Name  string    `yaml:"name" required:"true"`
	ID    uint      `yaml:"id" required:"true"`
	Certs []AppCert `yaml:"certs"`
}

type ConfigFile struct {
	Email        string      `yaml:"email" validate:"email" required:"true"`
	Apps         []App       `yaml:"apps"`
	Retries      uint        `yaml:"retries" default:"0"`
	CooldownSec  uint        `yaml:"cooldown_sec" default:"1"` // do not overwhelm the ACME service with requests - speed is not that important for requesting certs
	PathWeb      string      `yaml:"path_web" validate:"path"` // only required if certs use http-01
	PathCerts    string      `yaml:"path_certs" required:"true" validate:"path"`
	Concat       *bool       `yaml:"concat" default:"false"`
	FileModeCert os.FileMode `yaml:"file_mode_cert" default:"0644"`
	FileModeKey  os.FileMode `yaml:"file_mode_key" default:"0600"`
	FileGroup    string      `yaml:"file_group"`
	HookCmd      string      `yaml:"hook_cmd"`
}
```

</details>

Loading a YAML config:

```go
func LoadConfig(path string) (*ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cnf ConfigFile
	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	if err := defaults.Set(&cnf); err != nil {
	    return &cnf, fmt.Errorf("Failed to set config-defaults: %v", err)
	}

	v := &validate.StructValidator{}
	validationErrors := v.Validate(cnf)
	if len(validationErrors) > 0 {
	    return &cnf, fmt.Errorf("Got invalid config: %v", validationErrors)
	}

	return &cnf, nil
}
```

----

## Contribute

* Report bugs, problems & ideas on how to improve it
* Extend the [Unit-Tests for the validation-logic](https://github.com/O-X-L/go-validator/blob/latest/pkg/validate/struct_check_test.go) 🙏
* Implement more commonly used default-validators

## License

MIT - free to use
