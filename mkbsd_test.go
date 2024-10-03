package main

import "testing"

func TestGetExtension(t *testing.T) {
	tests := []struct {
		url      string
		expected string
	}{
		{"https://example.com/image.jpg", ".jpg"},
		{"https://example.com/image.png", ".png"},
		{"https://example.com/image.jpg", ".jpg"},
		{"://invalid-url`gwerf", ""},
	}

	for _, tt := range tests {
		got, err := getExtension(tt.url)
		if err != nil && tt.expected != "" {
			t.Errorf("expected no error, got %v", err)
		}
		if got != tt.expected {
			t.Errorf("got %v, want %v", got, tt.expected)
		}
	}

}
