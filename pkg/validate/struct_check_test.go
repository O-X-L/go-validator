package validate

import (
	"reflect"
	"strings"
	"testing"
)

type testCase struct {
	Name        string
	Input       interface{}
	WantErr     bool
	ErrContains string
}

type RequiredInt struct {
	SomeNum int `required:"true"`
}
type RequiredList struct {
	List []string `required:"true"`
}
type RequiredNested struct {
	Num RequiredInt
	Arr RequiredList
}
type NotRequiredPort struct {
	Port int `validate:"port"`
}

type BadTagFormat struct {
	FieldA string `required_if:"badformat"`
}
type BadFieldName struct {
	FieldA string `required_if:"BadFieldName=^.*$"`
}
type RegexWrongType struct {
	FieldA int `validate:"alphanumeric"`
}
type CustomWrongType struct {
	FieldA string `validate:"port"`
}
type UnknownTag struct {
	FieldA string `validate:"unknown_tag"`
}
type NestedConfig struct {
	Address string `validate_regex:"^[-a-zA-Z0-9\\s]{1,50}$"`
}
type Config struct {
	NameSafe string `validate:"alphanumeric_dash_underscore"`
	Port     int    `validate:"port"`
	Location NestedConfig
	Manager  *NestedConfig
	Domain   string `validate:"domain"`
}
type SliceItem struct {
	Name   string `validate:"alphanumeric_dash_underscore"`
	Port   int    `validate:"port"`
	Domain string `required:"true"`
}
type SliceConfig struct {
	Items []SliceItem
}
type PtrSliceConfig struct {
	Items []*SliceItem
}

func runTests(t *testing.T, tests []testCase, v *StructValidator) {
	if v == nil {
		v = &StructValidator{}
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			errors := v.Validate(tt.Input)
			errOccurred := len(errors) != 0

			// 1. Check if error presence matches expectation
			if errOccurred != tt.WantErr {
				t.Errorf("FAIL: %s\n  Expected error: %v, Got error: %v", tt.Name, tt.WantErr, errOccurred)
				if errOccurred {
					t.Logf("  Actual Errors encountered:")
					for _, e := range errors {
						t.Logf("    - %v", e)
					}
				}
				return
			}

			// 2. If errors were expected, verify the content and print failures
			if tt.WantErr && tt.ErrContains != "" {
				anyMatches := false
				errorsStr := []string{}
				for _, e := range errors {
					errorsStr = append(errorsStr, e.Error())
					if strings.Contains(e.Error(), tt.ErrContains) {
						anyMatches = true
					}
				}
				if !anyMatches {
					t.Errorf("FAIL: %s\n  Expected error to contain: %q\n  Actual errors returned:\n    %s",
						tt.Name, tt.ErrContains, strings.Join(errorsStr, "\n    "))
				}
			}
		})
	}
}

func TestValidateBase(t *testing.T) {
	tests := []testCase{
		{
			Name:    "Base (Input nil)",
			Input:   nil,
			WantErr: false, // Should not panic or error
		},
		{
			Name:    "Base (Input not a struct)",
			Input:   "this is just a string",
			WantErr: false, // Should not panic or error
		},
		{
			Name: "Base (Valid Nested Pointer Nil)",
			Input: &Config{
				NameSafe: "valid_name-123",
				Port:     1,
				Location: NestedConfig{Address: "123 Valid St"},
				Manager:  nil, // Nil pointer should be skipped, not cause panic
			},
			WantErr: false,
		},
		{
			Name: "Base (Valid nested Config)",
			Input: &Config{
				NameSafe: "valid_name-123",
				Port:     1,
				Location: NestedConfig{Address: "123 Valid St"},
				Manager:  &NestedConfig{Address: "456 Manager Ave"},
			},
			WantErr: false,
		},
	}
	runTests(t, tests, nil)
}

func TestValidateRequired(t *testing.T) {
	tests := []testCase{
		// failing
		{
			Name:        "Required (int no)",
			Input:       &RequiredInt{},
			WantErr:     true,
			ErrContains: "'SomeNum' is required but is empty",
		},
		{
			Name:        "Required (int zero)",
			Input:       &RequiredInt{SomeNum: 0}, // 0 is the zero value for int
			WantErr:     true,
			ErrContains: "'SomeNum' is required but is empty",
		},
		{
			Name:        "Required (list no)",
			Input:       &RequiredList{},
			WantErr:     true,
			ErrContains: "'List' is required but is empty",
		},
		{
			Name:        "Required (list zero)",
			Input:       &RequiredList{List: []string{}},
			WantErr:     true,
			ErrContains: "'List' is required but is empty",
		},
		// working
		{
			Name:    "Required (int ok)",
			Input:   &RequiredInt{SomeNum: 2},
			WantErr: false,
		},
		{
			Name:    "Required (list ok)",
			Input:   &RequiredList{List: []string{"test"}},
			WantErr: false,
		},
		{
			Name:    "Not-Required (port validation-skipped)",
			Input:   &NotRequiredPort{Port: 0}, // would not pass validation - but is skipped because not required
			WantErr: false,
		},
	}
	runTests(t, tests, nil)
}

func TestValidateRequiredNested(t *testing.T) {
	tests := []testCase{
		// failing
		{
			Name:        "Nested-Required (empty all)",
			Input:       &RequiredNested{},
			WantErr:     true,
			ErrContains: "Num - field 'SomeNum' is required but is empty",
		},
		{
			Name:        "Nested-Required (null-values)",
			Input:       &RequiredNested{Num: RequiredInt{SomeNum: 0}, Arr: RequiredList{List: []string{}}},
			WantErr:     true,
			ErrContains: "Num - field 'SomeNum' is required but is empty",
		},
		{
			Name:        "Nested-Required (empty list)",
			Input:       &RequiredNested{Num: RequiredInt{SomeNum: 1}, Arr: RequiredList{List: []string{}}},
			WantErr:     true,
			ErrContains: "Arr - field 'List' is required but is empty",
		},
		// working
		{
			Name:    "Nested-Required (ok)",
			Input:   &RequiredNested{Num: RequiredInt{SomeNum: 1}, Arr: RequiredList{List: []string{"a"}}},
			WantErr: false,
		},
	}
	runTests(t, tests, nil)
}

func TestValidateCustom(t *testing.T) {
	tests := []testCase{
		// failing
		{
			Name:        "custom-validate (invalid Port <)",
			Input:       &Config{NameSafe: "ok", Port: -1},
			WantErr:     true,
			ErrContains: "field 'Port' failed validation 'port'",
		},
		{
			Name:        "custom-validate (invalid Port >)",
			Input:       &Config{NameSafe: "ok", Port: 70000},
			WantErr:     true,
			ErrContains: "field 'Port' failed validation 'port'",
		},
		{
			Name:        "custom-validate (wrong type)",
			Input:       &CustomWrongType{FieldA: "not-an-int"}, // port expects an int
			WantErr:     true,
			ErrContains: "field 'FieldA' failed validation 'port'",
		},
		{
			Name:        "custom-validate (tag not found)",
			Input:       &UnknownTag{FieldA: "test"},
			WantErr:     true,
			ErrContains: "no validator registered for tag: 'unknown_tag' on field 'FieldA'",
		},
		{
			Name:        "custom-validate (bad Domain 1)",
			Input:       &Config{NameSafe: "ok", Domain: "bad"},
			WantErr:     true,
			ErrContains: "field 'Domain' failed validation 'domain'",
		},
		{
			Name:        "custom-validate (bad Domain 2)",
			Input:       &Config{NameSafe: "ok", Domain: "ox!.at"},
			WantErr:     true,
			ErrContains: "field 'Domain' failed validation 'domain'",
		},
		// working
		{
			Name:    "custom-validate (good Port)",
			Input:   &Config{NameSafe: "ok", Port: 30000},
			WantErr: false,
		},
		{
			Name:    "custom-validate (good Domain 1)",
			Input:   &Config{NameSafe: "ok", Domain: "oxl.at"},
			WantErr: false,
		},
		{
			Name:    "custom-validate (good Domain 2)",
			Input:   &Config{NameSafe: "ok", Domain: "localhost"},
			WantErr: false,
		},
	}
	runTests(t, tests, nil)
}

func TestValidateRegex(t *testing.T) {
	tests := []testCase{
		// failing
		{
			Name:    "regex-validate (non-string)",
			Input:   &RegexWrongType{FieldA: 123},
			WantErr: false,
		},
		{
			Name:        "regex-validate (invalid NameSafe)",
			Input:       &Config{NameSafe: "invalid!"},
			WantErr:     true,
			ErrContains: "field 'NameSafe' failed validation 'alphanumeric_dash_underscore'",
		},
		{
			Name:        "regex-validate (invalid Address - nested)",
			Input:       &Config{NameSafe: "ok", Port: 1111, Location: NestedConfig{Address: "123 Bad St!"}},
			WantErr:     true,
			ErrContains: "Location - field 'Address' failed validation",
		},
		{
			Name:        "regex-validate (invalid Address - nested ptr)",
			Input:       &Config{NameSafe: "ok", Port: 1111, Manager: &NestedConfig{Address: "456 Bad Ptr!"}},
			WantErr:     true,
			ErrContains: "Manager - field 'Address' failed validation",
		},
	}
	runTests(t, tests, nil)
}

func TestRequiredIf(t *testing.T) {
	type testGeoIP struct {
		// fallback provider: github.com/O-X-L/geoip-asn
		Provider   string `json:"provider" yaml:"provider" default:"oxl" validate_regex:"^(oxl|ipinfo_lite|maxmind_lite)$"`
		APIAccount string `json:"api_account" yaml:"api_account" required_if:"Provider=maxmind.*"`
		APIToken   string `json:"api_token" yaml:"api_token" required_if:"Provider=(ipinfo|maxmind).*"`
	}

	tests := []testCase{
		// failing
		{
			Name:        "required_if (maxmind, account missing)",
			Input:       &testGeoIP{Provider: "maxmind_lite", APIAccount: ""},
			WantErr:     true,
			ErrContains: "field 'APIAccount' is required because 'Provider'",
		},
		{
			Name:        "required_if (maxmind, token missing)",
			Input:       &testGeoIP{Provider: "maxmind_lite", APIAccount: "acc123", APIToken: ""},
			WantErr:     true,
			ErrContains: "field 'APIToken' is required because 'Provider'",
		},
		{
			Name:        "required_if (ipinfo, token missing)",
			Input:       &testGeoIP{Provider: "ipinfo_lite", APIToken: ""},
			WantErr:     true,
			ErrContains: "field 'APIToken' is required because 'Provider'",
		},
		{
			Name:        "required_if (bad tag format)",
			Input:       &BadTagFormat{FieldA: ""},
			WantErr:     true,
			ErrContains: "invalid required_if tag format on field 'FieldA'",
		},
		{
			Name:        "required_if (bad field name)",
			Input:       &BadFieldName{FieldA: ""},
			WantErr:     true,
			ErrContains: "invalid field name 'BadFieldName' in required_if tag on field 'FieldA'",
		},
		// working
		{
			Name: "Valid GeoIP (Default)",
			Input: &testGeoIP{
				Provider: "oxl",
			},
			WantErr: false,
		},
		{
			Name: "Valid GeoIP (Maxmind)",
			Input: &testGeoIP{
				Provider:   "maxmind_lite",
				APIAccount: "acc123",
				APIToken:   "tok123",
			},
			WantErr: false,
		},
		{
			Name: "Valid GeoIP (ipinfo)",
			Input: &testGeoIP{
				Provider: "ipinfo_lite",
				APIToken: "tok123",
			},
			WantErr: false,
		},
	}
	runTests(t, tests, nil)
}

func TestValidateSliceOfStructs(t *testing.T) {
	tests := []testCase{
		// --- Valid Cases ---
		{
			Name: "Slice (Valid)",
			Input: &SliceConfig{
				Items: []SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"},
					{Name: "item-2", Port: 443, Domain: "test.com"},
				},
			},
			WantErr: false,
		},
		{
			Name: "Slice (Valid Pointers)",
			Input: &PtrSliceConfig{
				Items: []*SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"},
					{Name: "item-2", Port: 443, Domain: "test.com"},
				},
			},
			WantErr: false,
		},
		{
			Name:    "Slice (Valid Empty Slice)",
			Input:   &SliceConfig{Items: []SliceItem{}},
			WantErr: false,
		},
		{
			Name:    "Slice (Valid Nil Slice)",
			Input:   &SliceConfig{Items: nil},
			WantErr: false,
		},
		{
			Name: "Slice (Valid with Nil Pointer in Slice)",
			Input: &PtrSliceConfig{
				Items: []*SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"},
					nil, // Nil pointers should be skipped
					{Name: "item-2", Port: 443, Domain: "test.com"},
				},
			},
			WantErr: false,
		},

		// --- Invalid Cases (Slice of Structs) ---
		{
			Name: "Slice (Invalid - required)",
			Input: &SliceConfig{
				Items: []SliceItem{
					{Name: "item-1", Port: 80, Domain: ""}, // Domain is required
				},
			},
			WantErr:     true,
			ErrContains: "Items[0] - field 'Domain' is required but is empty",
		},
		{
			Name: "Slice (Invalid - custom validation)",
			Input: &SliceConfig{
				Items: []SliceItem{
					{Name: "item-1", Port: -1, Domain: "ok.com"}, // Port is invalid
				},
			},
			WantErr:     true,
			ErrContains: "Items[0] - field 'Port' failed validation 'port'",
		},
		{
			Name: "Slice (Invalid - regex validation)",
			Input: &SliceConfig{
				Items: []SliceItem{
					{Name: "bad!", Port: 80, Domain: "ok.com"}, // Name is invalid
				},
			},
			WantErr:     true,
			ErrContains: "Items[0] - field 'Name' failed validation 'alphanumeric_dash_underscore'",
		},
		{
			Name: "Slice (Invalid - second item)",
			Input: &SliceConfig{
				Items: []SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"},      // Valid
					{Name: "item-2", Port: 99999, Domain: "test.com"}, // Invalid Port
				},
			},
			WantErr:     true,
			ErrContains: "Items[1] - field 'Port' failed validation 'port'",
		},

		// --- Invalid Cases (Slice of Pointers to Structs) ---
		{
			Name: "Slice (Invalid Pointers - required)",
			Input: &PtrSliceConfig{
				Items: []*SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"}, // Valid
					{Name: "item-2", Port: 443, Domain: ""},      // Invalid Domain
				},
			},
			WantErr:     true,
			ErrContains: "Items[1] - field 'Domain' is required but is empty",
		},
		{
			Name: "Slice (Invalid Pointers - with nil)",
			Input: &PtrSliceConfig{
				Items: []*SliceItem{
					{Name: "item-1", Port: 80, Domain: "ok.com"}, // Valid
					nil, // Skipped
					{Name: "item-3", Port: -1, Domain: "test.com"}, // Invalid Port
				},
			},
			WantErr:     true,
			ErrContains: "Items[2] - field 'Port' failed validation 'port'",
		},
	}
	runTests(t, tests, nil)
}

func TestValidateSliceInvalidType(t *testing.T) {
	type SlicePort struct {
		Ports []string `validate:"port"` // Port validator expects []uint or uint
	}

	input := &SlicePort{Ports: []string{"80"}}
	v := StructValidator{}
	errs := v.Validate(input)

	if len(errs) == 0 {
		t.Error("Expected error when passing string slice to port validator, but got none")
	}
}

func TestValidateRequiredIfLogic(t *testing.T) {
	// Testing the specific logic in checkRequiredIfField used in global configs
	type TestReqIf struct {
		Provider string
		Token    string `required_if:"Provider=maxmind.*"`
	}

	value := reflect.ValueOf(TestReqIf{Provider: "maxmind_lite", Token: ""})
	field, _ := value.Type().FieldByName("Token")

	v := StructValidator{}

	err := v.checkRequiredIfField(value.FieldByName("Token"), field, value)
	if err == nil {
		t.Error("Expected error for missing Token when Provider is maxmind_lite")
	}

	vValid := reflect.ValueOf(TestReqIf{Provider: "maxmind_lite", Token: "secret"})
	errValid := v.checkRequiredIfField(vValid.FieldByName("Token"), field, vValid)
	if errValid != nil {
		t.Errorf("Expected no error for valid Token, got: %v", errValid)
	}
}
