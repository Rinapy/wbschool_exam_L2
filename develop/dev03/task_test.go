package main

import "testing"

func TestSorter(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    string
		wantErr error
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

		})
	}
}
