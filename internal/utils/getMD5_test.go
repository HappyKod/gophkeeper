package utils

import "testing"

func TestGetMD5Hash(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte("hello"), "5d41402abc4b2a76b9719d911017c592"},
		{[]byte("world"), "7d793037a0760186574b0282f2f435e7"},
		{[]byte(""), "d41d8cd98f00b204e9800998ecf8427e"},
	}

	for _, tt := range tests {
		got := GetMD5Hash(tt.input)
		if got != tt.want {
			t.Errorf("GetMD5Hash(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
