package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestInitFunc(t *testing.T) {
	fmt.Println("Testing init function")
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-A", "2", "-B", "2", "-C", "3", "-i", "-v", "-F", "-n", "Искомое", "input_1.txt"} // Пример аргументов командной строки "-B", "2", "-C", "3", "-i", "-v", "-F", "-n",
	// Проверка результатов
	parseFlags()
	expectedA := 2
	expectedB := 2
	expectedC := 3
	expectedI := true
	expectedV := true
	expectedF := true
	expectedN := true
	expectedStr := "Искомое"

	if A != expectedA {
		t.Errorf("Expected A to be %d, but got %d", expectedA, A)
	}
	if B != expectedB {
		t.Errorf("Expected B to be %v, but got %v", expectedB, B)
	}
	if C != expectedC {
		t.Errorf("Expected C to be %v, but got %v", expectedC, C)
	}
	if i != expectedI {
		t.Errorf("Expected i to be %v, but got %v", expectedI, i)
	}
	if v != expectedV {
		t.Errorf("Expected v to be %v, but got %v", expectedV, v)
	}
	if F != expectedF {
		t.Errorf("Expected F to be %v, but got %v", expectedF, F)
	}
	if n != expectedN {
		t.Errorf("Expected n to be %v, but got %v", expectedN, n)
	}
	if searchStr != expectedStr {
		t.Errorf("Expected searchStr to be %v, but got %v", expectedStr, searchStr)
	}
}

func TestParseNameOrPatternFunc(t *testing.T) {
	fmt.Println("Testing parseNameOrPattern function")
	tests := []struct {
		name    string
		arg     []string
		want    []string
		wantErr error
	}{
		{
			name:    "Одиночный файл test_1.txt",
			arg:     []string{"cmd", "слово", "test_1.txt"},
			want:    []string{"test_1.txt"},
			wantErr: nil,
		},
		{
			name:    "Паттерн файлов test*.txt",
			arg:     []string{"cmd", "слово", "test*.txt"},
			want:    []string{"test_1.txt", "test_2.txt", "test_3.txt"},
			wantErr: nil,
		},
		{
			name:    "Ошибка поиска файла",
			arg:     []string{"cmd", "слово", "test.txt"},
			want:    nil,
			wantErr: ErrFindFiles,
		},
		{
			name:    "Ошибка паттерна",
			arg:     []string{"cmd", "слово", "*/*.txt"},
			want:    nil,
			wantErr: ErrFindFiles,
		},
	}
	for _, tt := range tests {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = tt.arg
		parseFlags()
		fmt.Println(tt.name)
		fsn, err := parseNameOrPattern(flag.Arg(1))
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("parseNameOrPattern error = %v, want.err %v", err, tt.wantErr)
			return
		}
		for i := 0; i < len(tt.want); i++ {
			if fsn[i] != tt.want[i] {
				t.Errorf("Expected file name '%s', but got '%s'", tt.want[i], fsn[i])
			}
		}
	}
}

func TestReadLinesFunc(t *testing.T) {
	fmt.Println("Testing parseNameOrPattern function")
	tests := []struct {
		name    string
		arg     []string
		want    []Line
		wantErr error
	}{
		{
			name:    "Одиночный файл test_1.txt",
			arg:     []string{"cmd", "слово", "test_1.txt"},
			want:    []string{"test_1.txt"},
			wantErr: nil,
		},
		{
			name:    "Паттерн файлов test*.txt",
			arg:     []string{"cmd", "слово", "test*.txt"},
			want:    []string{"test_1.txt", "test_2.txt", "test_3.txt"},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = tt.arg
		parseFlags()
		fsn, err := parseNameOrPattern(flag.Arg(1))
		fmt.Println(tt.name)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Sorter() error = %v, want.err %v", err, tt.wantErr)
			return
		}
		for i := 0; i < len(tt.want); i++ {
			if fsn[i] != tt.want[i] {
				t.Errorf("Expected file name '%s', but got '%s'", tt.want[i], fsn[i])
			}
		}
	}
}
