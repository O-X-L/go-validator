package custom

import "testing"

func TestValidateDomain(t *testing.T) {
	if !ValidateDomain("oxl.at") {
		t.Error("Domain validation failed #1")
	}
	if !ValidateDomain("abc.test.oxl.at") {
		t.Error("Domain validation failed #2")
	}
	if !ValidateDomain("localhost") {
		t.Error("Domain validation failed #3")
	}
	if ValidateDomain("2") {
		t.Error("Domain validation failed #4")
	}
	if !ValidateDomain("163.com") {
		t.Error("Domain validation failed #5")
	}
	if !ValidateDomain("host-svc.com") {
		t.Error("Domain validation failed #6")
	}
}

func TestValidatePort(t *testing.T) {
	if !ValidatePort(200) {
		t.Error("Port validation failed #1")
	}
	if ValidatePort(0) {
		t.Error("Port validation failed #2")
	}
	if ValidatePort(-1) {
		t.Error("Port validation failed #3")
	}
	if ValidatePort(70000) {
		t.Error("Port validation failed #4")
	}
	if !ValidatePort(1234) {
		t.Error("Port validation failed #5")
	}
	if !ValidatePort(30000) {
		t.Error("Port validation failed #6")
	}
}

func TestValidateIP(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{"127.0.0.1", true},
		{"8.8.8.8", true},
		{"2001:db8::1", true},
		{"not-an-ip", false},
		{"256.256.256.256", false},
		{123, false}, // Wrong type
	}
	for _, tt := range tests {
		if ValidateIP(tt.input) != tt.want {
			t.Errorf("ValidateIP(%v) = %v; want %v", tt.input, !tt.want, tt.want)
		}
	}
}

func TestValidateIPDomain(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{"oxl.at", true},
		{"1.1.1.1", true},
		{"localhost", true},
		{"invalid_ip", false},
	}
	for _, tt := range tests {
		if ValidateIPDomain(tt.input) != tt.want {
			t.Errorf("ValidateIPDomain(%v) = %v; want %v", tt.input, !tt.want, tt.want)
		}
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{"oxl.at", false},
		{"test@oxl.at", true},
		{"test@localhost", true},
		{"invalid", false},
		{"dddd@dslfdsk@dt", false},
		{"ddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd@oxl.at", false},
	}
	for _, tt := range tests {
		if ValidateEmail(tt.input) != tt.want {
			t.Errorf("ValidateEmail(%v) = %v; want %v", tt.input, !tt.want, tt.want)
		}
	}
}

func TestValidateNoDuplicates(t *testing.T) {
	tests := []struct {
		input interface{}
		want  bool
	}{
		{"oxl.at", false},
		{42, false},
		{[]int{2}, false},
		{[]int{2, 4}, false},
		{[]int{2, 4, 2}, true},
		{[]string{"a"}, false},
		{[]string{"a", "b"}, false},
		{[]string{"a", "b", "a"}, true},
	}
	for _, tt := range tests {
		if ValidateNoDuplicates(tt.input) != tt.want {
			t.Errorf("ValidateNoDuplicates(%v) = %v; want %v", tt.input, !tt.want, tt.want)
		}
	}
}
