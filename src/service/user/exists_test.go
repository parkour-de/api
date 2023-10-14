package user

import "testing"

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			"empty",
			"",
			true,
		},
		{
			"too short",
			"ab",
			true,
		},
		{
			"too long",
			"1234567890123456789012345678901",
			true,
		},
		{
			"dollar",
			"abc$",
			true,
		},
		{
			"percent",
			"abc%",
			true,
		},
		{
			"only digits",
			"12345",
			true,
		},
		{
			"valid",
			"abc123",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateCustomKey(tt.username); (err != nil) != tt.wantErr {
				t.Errorf("ValidateCustomKey(%#v) error = %v, wantErr %v", tt.username, err, tt.wantErr)
			}
		})
	}
}
