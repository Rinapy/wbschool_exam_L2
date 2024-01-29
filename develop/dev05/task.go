package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var (
	A int
	B int
	C int
	c bool
	i bool
	v bool
	F bool
	n bool
)

type Line struct {
	Text string
}
type lineSlice []Line

type File struct {
	name     string
	lines    lineSlice
	strCount int
	findStr  []int
}

type FileSlice []File

func parseNameOrPattern(arg string) []string {
	fs := make([]string, 0)
	isPattern := false
	for _, val := range arg {
		if string(val) == "*" {
			isPattern = true
			break
		}
	}
	if isPattern {
		files, err := filepath.Glob(arg)
		if err != nil {
			fmt.Println("Ошибка поиска файлов:", err)
		}

		for _, file := range files {
			fs = append(fs, file)
		}
	} else {
		fs = append(fs, arg)
	}
	return fs
}

func parseFlags() ([]string, string) {
	flag.IntVar(&A, "A", 0, "\"after\" печатать +N строк после совпадения")
	flag.IntVar(&B, "B", 0, "\"before\" печатать +N строк до совпадения")
	flag.IntVar(&C, "C", 0, "\"context\" (A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&c, "c", false, "\"count\" (количество строк)")
	flag.BoolVar(&i, "i", false, "\"ignore-case\" (игнорировать регистр)")
	flag.BoolVar(&v, "v", false, "\"invert\" (вместо совпадения, исключать)")
	flag.BoolVar(&F, "F", false, "\"fixed\", точное совпадение со строкой, не паттерн")
	flag.BoolVar(&n, "n", false, "\"line num\", печатать номер строки")
	flag.Parse()
	fs := parseNameOrPattern(flag.Arg(1))
	str := flag.Arg(0)
	return fs, str
}

func fillFileSlice(filenames []string) (FileSlice, error) {
	fileSlice := FileSlice{}

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		lines, err := readLines(file)
		if err != nil {
			return nil, err
		}

		fileInfo := File{
			name:     filename,
			lines:    lines,
			strCount: len(lines),
			findStr:  []int{},
		}

		fileSlice = append(fileSlice, fileInfo)
	}

	return fileSlice, nil
}

func readLines(file *os.File) ([]Line, error) {
	reader := bufio.NewReader(file)
	var lines []Line

	for {
		lineText, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}

		cleanStr := strings.TrimRight(lineText, "\r\n")
		line := Line{Text: cleanStr}
		lines = append(lines, line)
	}

	return lines, nil
}

func printer(fileSlice FileSlice) {
	match := "Match -- %v\n"
	before := "Before -- %v\n"
	after := "After  -- %v\n"
	matchN := "Match -- %v [%v]\n"
	beforeN := "Before -- %v [%v]\n"
	afterN := "After  -- %v [%v]\n"

	for _, file := range fileSlice {
		if len(file.findStr) > 1 {
			fmt.Printf("---------------------------%v---------------------------\n", file.name)
		}

		if c {
			fmt.Printf(match, len(file.findStr))
			continue
		}

		for _, idx := range file.findStr {
			if C > 0 {
				A = C
				B = C
			}

			if B > 0 {
				for iB := 1; iB <= B; iB++ {
					beforeIdx := idx - iB
					if beforeIdx >= 0 {
						if n {
							fmt.Printf(beforeN, file.lines[beforeIdx].Text, beforeIdx+1)
							continue
						}
						fmt.Printf(before, file.lines[beforeIdx].Text)
					}
				}
			}

			if n {
				fmt.Printf(matchN, file.lines[idx].Text, idx+1)
			} else {
				fmt.Printf(match, file.lines[idx].Text)
			}

			if A > 0 {
				for iA := 1; iA <= A; iA++ {
					afterIdx := idx + iA
					if n {
						fmt.Printf(afterN, file.lines[afterIdx].Text, afterIdx+1)
						continue
					}
					fmt.Printf(after, file.lines[afterIdx].Text)
				}
			}
		}
	}
}

func Finder(fileSlice FileSlice, text string) {
	found := false
	for idx := range fileSlice {
		file := &fileSlice[idx]
		for idx, line := range file.lines {
			if !F {
				if !v {
					if (i && strings.Contains(strings.ToLower(line.Text), strings.ToLower(text))) ||
						(!i && strings.Contains(line.Text, text)) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				} else {
					if A != 0 || B != 0 || C != 0 {
						A = 0
						B = 0
						C = 0
						fmt.Println("Использование флагов -A -B -C при флаге -v недопустимо, данные флаги отключены.")
					}

					if (i && !strings.Contains(strings.ToLower(line.Text), strings.ToLower(text))) ||
						(!i && !strings.Contains(line.Text, text)) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				}
			} else {
				if !v {
					if (i && strings.ToLower(line.Text) == strings.ToLower(text)) ||
						(!i && line.Text == text) {
						file.findStr = append(file.findStr, idx)
						found = true

					}
				} else {
					if A != 0 || B != 0 || C != 0 {
						A = 0
						B = 0
						C = 0
						fmt.Println("Использование флагов -A -B -C при флаге -v недопустимо, данные флаги отключены.")
					}

					if (i && strings.ToLower(line.Text) != strings.ToLower(text)) ||
						(!i && line.Text != text) {
						file.findStr = append(file.findStr, idx)
						found = true
					}
				}
			}

		}

	}
	if found {
		printer(fileSlice)
	} else {
		fmt.Println("Совпадений не найдено")
	}

}

func main() {
	fsn, str := parseFlags()
	//fmt.Println(fsn)
	//fmt.Println(str)
	fs, err := fillFileSlice(fsn)
	if err != nil {
		fmt.Println(err)
	}
	Finder(fs, str)

}
