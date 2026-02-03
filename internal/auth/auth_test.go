package auth

import "testing"

func TestAuthenticator_Verify(t *testing.T) {
	a := NewAuthenticator("admin", "123456")

	tests := []struct {
		name    string
		encoded string
		want    bool
	}{
		{"Valid Credentials", "YWRtaW46MTIzNDU2", true},
		{"Invalid Credentials", "YWRtaW46d3Jvbmc=", false},
		{"Empty String", "", false},
		{"Invalid Base64", "!!!", false},
		{"Wrong Format", "YWRtaW4=", false}, // "admin"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := a.Verify(tt.encoded); got != tt.want {
				t.Errorf("Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthenticator_VerifyHTTPHeader(t *testing.T) {
	a := NewAuthenticator("user", "pass")
	
	if !a.VerifyHTTPHeader("Basic dXNlcjpwYXNz") {
		t.Error("VerifyHTTPHeader failed for valid header")
	}

	if a.VerifyHTTPHeader("Bearer token") {
		t.Error("VerifyHTTPHeader should fail for non-Basic header")
	}
}
