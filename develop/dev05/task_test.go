package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"slices"
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
			name:    "Test-Case: Одиночный файл test_1.txt",
			arg:     []string{"cmd", "слово", "testfiles\\test_1.txt"},
			want:    []string{"testfiles\\test_1.txt"},
			wantErr: nil,
		},
		{
			name:    "Test-Case: Паттерн файлов test*.txt",
			arg:     []string{"cmd", "слово", "./testfiles/test*.txt"},
			want:    []string{"testfiles\\test_1.txt", "testfiles\\test_2.txt", "testfiles\\test_3.txt"},
			wantErr: nil,
		},
		{
			name:    "Test-Case: Ошибка поиска файла",
			arg:     []string{"cmd", "слово", "./testfiles/test.txt"},
			want:    nil,
			wantErr: ErrFindFiles,
		},
		{
			name:    "Test-Case: Ошибка паттерна",
			arg:     []string{"cmd", "слово", "./testfiles/*/*.txt"},
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
	fmt.Println("Testing readLines function")
	files := []string{"./testfiles/input_1.txt", "./testfiles/input_2.txt"}
	expected := []lineSlice{
		{
			Line{Text: "Туту"},
			Line{Text: "Там"},
			Line{Text: "Сюда"},
			Line{Text: "туда"},
			Line{Text: "1"},
			Line{Text: "3"},
		},
		{
			Line{Text: "Там"},
			Line{Text: "сям"},
			Line{Text: "Почему"},
			Line{Text: "Потому"},
			Line{Text: "1"},
			Line{Text: "3"},
		},
	}

	for idx, fn := range files {
		file, err := os.Open(fn)
		if err != nil {
			t.Errorf("Ошибка открытия тестового файла %v error: %v\n", fn, err)
		}
		defer file.Close()
		lines, err := readLines(file)
		if err != nil {
			t.Errorf("Ошибка чтения тестового файла %v error: %v\n", fn, err)
		}
		if !slices.Equal(lines, expected[idx]) {
			t.Errorf("Ошибка чтения строк ожидалось %v, получено %v.\n", expected[idx], lines)
		}
	}
}

func TestFillFileSliceFunc(t *testing.T) {
	fmt.Println("Testing fillFileSlice function")
	files := []string{"./testfiles/input_1.txt", "./testfiles/input_2.txt"}
	expected := FileSlice{
		File{
			name: "./testfiles/input_1.txt",
			lines: lineSlice{
				Line{Text: "Туту"},
				Line{Text: "Там"},
				Line{Text: "Сюда"},
				Line{Text: "туда"},
				Line{Text: "1"},
				Line{Text: "3"},
			},
			findStr: []int{},
		},
		File{
			name: "./testfiles/input_2.txt",
			lines: lineSlice{
				Line{Text: "Там"},
				Line{Text: "сям"},
				Line{Text: "Почему"},
				Line{Text: "Потому"},
				Line{Text: "1"},
				Line{Text: "3"},
			},
			findStr: []int{},
		},
	}

	for idx, fn := range files {
		fsn, err := fillFileSlice(files)
		if err != nil {
			t.Errorf("Ошибка чтения тестового файла %v error: %v\n", fn, err)
		}
		if !reflect.DeepEqual(expected[idx], fsn[idx]) {
			t.Errorf("Ошибка чтения строк ожидалось %v, получено %v.\n", expected[idx], fsn[idx])
		}
	}
}

func TestFinderFunc(t *testing.T) {
	tests := []struct {
		name    string
		arg     []string
		data    FileSlice
		wantRes []int
		wantErr error
	}{
		{
			name: "Test-Case: Есть совпадение в тексте",
			arg:  []string{"cmd", "apple"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
					},
				},
			},
			wantRes: []int{1},
			wantErr: nil,
		},
		{
			name: "Test-Case: Отсутствие совпадений",
			arg:  []string{"cmd", "12313"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
					},
				},
			},
			wantRes: nil,
			wantErr: ErrNoLinesFound,
		},
		{
			name: "Test-Case: Исключение совпадений",
			arg:  []string{"cmd", "-v", "apple"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
					},
				},
			},
			wantRes: []int{0},
			wantErr: nil,
		},
		{
			name: "Test-Case: Игнорировать регистр",
			arg:  []string{"cmd", "-i", "Apple"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
					},
				},
			},
			wantRes: []int{1},
			wantErr: nil,
		},
		{
			name: "Test-Case: Точное совпадение",
			arg:  []string{"cmd", "-F", "Apple"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
						Line{"Apple"},
					},
				},
			},
			wantRes: []int{2},
			wantErr: nil,
		},
		{
			name: "Test-Case: Точное совпадение c без учёта регистра",
			arg:  []string{"cmd", "-F", "-i", "apple"},
			data: FileSlice{
				File{
					lines: lineSlice{
						Line{"Test line 1"},
						Line{"Test line 2 contains the word apple"},
						Line{"Apple"},
					},
				},
			},
			wantRes: []int{2},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = tt.arg
		parseFlags()
		fmt.Println(tt.name)

		fs, err := finder(tt.data, searchStr)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("finder error = %v, want.err %v", err, tt.wantErr)
			return
		}
		if errors.Is(err, tt.wantErr) {
			continue
		}
		if !reflect.DeepEqual(fs[0].findStr, tt.wantRes) {
			t.Errorf("Expected findStr '%v', but got '%v'", tt.wantRes, fs[0].findStr)
		}
	}
}
