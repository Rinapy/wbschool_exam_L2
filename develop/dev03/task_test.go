package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

func TestFillNewFile(t *testing.T) {
	slice := LineSlice{
		Line{[]string{"Field1", "Field2", "Field3"}},
		Line{[]string{"Value1", "Value2"}},
	}
	// Вызов функции fillNewFile
	err := fillNewFile(&slice)
	if err != nil {
		t.Errorf("Func fileNewFile return err: %s", err)
	}

	// Проверка содержимого файла
	expectedLines := []string{
		"Field1 Field2 Field3\n",
		"Value1 Value2\n",
	}
	file, _ := os.Open("output.txt")
	reader := bufio.NewReader(file)
	actualLines := make([]string, 0)
	for {
		lineText, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		}
		actualLines = append(actualLines, lineText)
	}
	file.Close()
	os.Remove(file.Name())

	if len(actualLines) != len(expectedLines) {
		t.Errorf("Expected %d lines, but got %d", len(expectedLines), len(actualLines))
	}

	for i := 0; i < len(actualLines); i++ {
		if actualLines[i] != expectedLines[i] {
			t.Errorf("Expected line '%s', but got '%s'", expectedLines[i], actualLines[i])
		}
	}

}

func TestFillLineSlice(t *testing.T) {
	expectedSlice := LineSlice{
		Line{[]string{"Field1", "Field2", "Field3"}},
		Line{[]string{"Value1", "Value2"}},
	}

	file, _ := os.Create("input.txt")
	writer := bufio.NewWriter(file)

	for _, v := range expectedSlice {
		writer.WriteString(strings.Join(v.Fields, " ") + "\n")
	}
	writer.Flush()
	file.Close()

	testSlice, err := fillLineSlice("input.txt")
	os.Remove(file.Name())
	if err != nil {
		t.Errorf("Func fillLineSlice return err: %s", err)
	}

	if len(expectedSlice) != len(*testSlice) {
		t.Errorf("Expected %d lines, but got %d", len(expectedSlice), len(*testSlice))
	}

	for i := 0; i < len(expectedSlice); i++ {
		for j := 0; j < len(expectedSlice[i].Fields); j++ {
			if (*testSlice)[i].Fields[j] != expectedSlice[i].Fields[j] {
				t.Errorf("Expected line '%s', but got '%s'", expectedSlice[i].Fields[j], (*testSlice)[i].Fields[j])
			}
		}
	}

}

func TestParseFlags(t *testing.T) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Args = []string{"cmd", "-k", "1", "-n", "-r", "-u", "input.txt"} // Пример аргументов командной строки

	// Вызов функции Sorter
	parseFlags()

	// Проверка результатов
	expectedK := 1
	expectedN := true
	expectedR := true
	expectedU := true
	expectedInputFile := "input.txt"

	if k != expectedK {
		t.Errorf("Expected k to be %d, but got %d", expectedK, k)
	}
	if n != expectedN {
		t.Errorf("Expected n to be %v, but got %v", expectedN, n)
	}
	if r != expectedR {
		t.Errorf("Expected r to be %v, but got %v", expectedR, r)
	}
	if u != expectedU {
		t.Errorf("Expected u to be %v, but got %v", expectedU, u)
	}
	if inputFile != expectedInputFile {
		t.Errorf("Expected inputFile to be %s, but got %s", expectedInputFile, inputFile)
	}

}

func TestSorter(t *testing.T) {
	tests := []struct {
		name    string
		arg     []string
		want    LineSlice
		wantErr error
	}{
		{
			name: "Сортировка чисел по k = 0 -n",
			arg:  []string{"cmd", "-k", "0", "-n", "./testfiles/test.txt"},
			want: LineSlice{
				Line{[]string{"2", "8", "4"}},
				Line{[]string{"4", "7", "1"}},
				Line{[]string{"6", "1", "5"}},
				Line{[]string{"7", "3", "9"}},
				Line{[]string{"9", "5", "3"}},
			},
			wantErr: nil,
		},
		{
			name: "Сортировка по k = 1 -n",
			arg:  []string{"cmd", "-k", "1", "-n", "./testfiles/test.txt"},
			want: LineSlice{
				Line{[]string{"6", "1", "5"}},
				Line{[]string{"7", "3", "9"}},
				Line{[]string{"9", "5", "3"}},
				Line{[]string{"4", "7", "1"}},
				Line{[]string{"2", "8", "4"}},
			},
			wantErr: nil,
		},
		{
			name: "Обратная cортировка чисел флаг -k 0 -n -r",
			arg:  []string{"cmd", "-r", "-n", "./testfiles/test.txt"},
			want: LineSlice{
				Line{[]string{"9", "5", "3"}},
				Line{[]string{"7", "3", "9"}},
				Line{[]string{"6", "1", "5"}},
				Line{[]string{"4", "7", "1"}},
				Line{[]string{"2", "8", "4"}},
			},
			wantErr: nil,
		},
		{
			name: "Удаление дубликатов cортировка чисел флаг -k 0 -n -r",
			arg:  []string{"cmd", "-r", "-n", "-u", "./testfiles/testDup.txt"},
			want: LineSlice{
				Line{[]string{"8"}},
				Line{[]string{"7"}},
				Line{[]string{"6"}},
				Line{[]string{"5"}},
				Line{[]string{"4"}},
				Line{[]string{"3"}},
				Line{[]string{"2"}},
				Line{[]string{"1"}},
			},
			wantErr: nil,
		},
		{
			name:    "Ошибка сортировки по выходу за пределы доступных k",
			arg:     []string{"cmd", "-k", "23", "./testfiles/test.txt"},
			want:    nil,
			wantErr: &ErrIndexFile{},
		},
		{
			name:    "Ошибка отсутствуя переданного файла",
			arg:     []string{"cmd", "input.txt"},
			want:    nil,
			wantErr: &ErrOpenFile{},
		},
	}
	for _, tt := range tests {
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		os.Args = tt.arg
		lineSlice, err := Sorter()
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("Sorter() error = %v, want.err %v", err, tt.wantErr)
			return
		}
		for i := 0; i < len(tt.want); i++ {
			if (*lineSlice)[i].Fields[k] != tt.want[i].Fields[k] {
				t.Errorf("Expected line '%s', but got '%s'", tt.want[i].Fields, (*lineSlice)[i].Fields)
			}
		}
		os.Remove("output.txt")
	}
}
