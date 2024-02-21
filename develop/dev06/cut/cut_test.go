package cut

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseF(t *testing.T) {

	tests := []struct {
		name    string
		args    string
		wantRes []int
		wantErr error
	}{
		{
			name:    "Test -f 1",
			args:    "1",
			wantRes: []int{0},
			wantErr: nil,
		},
		{
			name:    "Test -f 1,3",
			args:    "1,3",
			wantRes: []int{0, 2},
			wantErr: nil,
		},
		{
			name:    "Test -f 1-3",
			args:    "1-3",
			wantRes: []int{0, 1, 2},
			wantErr: nil,
		},
		{
			name:    "Test -f 5-3",
			args:    "5-3",
			wantRes: nil,
			wantErr: &IndexValueError{},
		},
		{
			name:    "Test -f 0",
			args:    "0",
			wantRes: nil,
			wantErr: &IndexValueError{},
		},
		{
			name:    "Test -f str",
			args:    "str",
			wantRes: nil,
			wantErr: &IndexValueError{},
		},
		{
			name:    "Test -f str,str",
			args:    "str,str",
			wantRes: nil,
			wantErr: &IndexValueError{},
		},
		{
			name:    "Test -f str-str",
			args:    "str-str",
			wantRes: nil,
			wantErr: &IndexValueError{},
		},
	}
	for _, tt := range tests {
		cfg := NewCfg()
		fmt.Println(tt.name)
		err := cfg.parseF(tt.args)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("parseF() error = %v, want.err %v", err, tt.wantErr)
			return
		}
		for i, v := range tt.wantRes {
			if len(tt.wantRes) != len(cfg.f) {
				t.Errorf("Error, expected length of cfg.f %v but got %v", len(tt.wantRes), len(cfg.f))
				return
			}
			if v != cfg.f[i] {
				t.Errorf("Error result, want %v, got %v", v, cfg.f[i])
				return
			}
		}
	}
}

func TestParseData(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		wantRes lineSlice
		wantErr error
	}{
		{
			name: "Test 123:321",
			args: "123:321",
			wantRes: lineSlice{
				line{text: "123:321"},
			},
			wantErr: nil,
		},
		{
			name: "Test 123:321 456:654",
			args: "123:321 456:654",
			wantRes: lineSlice{
				line{text: "123:321"},
				line{text: "456:654"},
			},
			wantErr: nil,
		},
		{
			name:    "Test nil data",
			args:    "",
			wantRes: nil,
			wantErr: &DataNotFound{},
		},
	}
	for _, tt := range tests {
		cfg := NewCfg()
		cfg.ld = " "
		fmt.Println(tt.name)
		data, err := cfg.parseData(tt.args)
		fmt.Println(data)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("parseData() error = %v, want.err %v", err, tt.wantErr)
			return
		}
		for i, v := range tt.wantRes {
			if data[i].text != v.text {
				t.Errorf("Error result, want %v, got %v", v.text, data[i].text)
			}
		}
	}
}
